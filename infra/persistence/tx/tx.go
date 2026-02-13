package tx

import (
	"context"
	"fmt"
	"runtime"

	"github.com/google/wire"

	"purchase/infra/logx"
	"purchase/infra/persistence/dal/db/ent"
)

var ProviderSet = wire.NewSet(NewTransactionManager)

type TransactionPropagation int

const (
	PropagationRequired    TransactionPropagation = iota // 支持当前事务，如果当前没有事务，就新建一个事务
	PropagationRequiresNew                               // 新建事务，如果当前存在事务，把当前事务挂起，两个事务互不影响
	PropagationNested                                    // 支持当前事务，如果当前事务存在，则执行一个嵌套事务，如果当前没有事务，就新建一个事务, 在A事务里面嵌套B事务, B回滚不影响A, A回滚会让B回滚，A提交会让B提交
	PropagationNever                                     // 以非事务方式执行，如果当前存在事务，不在事务中执行
)

func defaultPropagation() TransactionPropagation {
	return PropagationRequired
}

type TransactionManager struct {
	db *ent.Client
}

func NewTransactionManager(db *ent.Client) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (m *TransactionManager) Transaction(ctx context.Context, fn func(ctx context.Context) error, propagations ...TransactionPropagation) error {
	propagation := defaultPropagation()
	if len(propagations) > 0 {
		propagation = propagations[0]
	}
	switch propagation {
	case PropagationNever:
		return m.withNeverPropagation(ctx, fn)
	case PropagationNested:
		return m.withNestedPropagation(ctx, fn)
	case PropagationRequired:
		return m.withRequiredPropagation(ctx, fn)
	case PropagationRequiresNew:
		return m.withRequiresNewPropagation(ctx, fn)
	default:
		return fmt.Errorf("not supported tx propagation: %d", propagation)
	}
}

// withNeverPropagation 不在事务中执行
func (m *TransactionManager) withNeverPropagation(ctx context.Context, fn func(ctx context.Context) error) error {
	txCtx, ok := ctx.(*TransactionContext)
	if ok {
		// 当前在事务中，提取原始 context，不在事务中执行
		return fn(txCtx.Context())
	}
	// 当前不在事务中，直接执行
	return fn(ctx)
}

// withNestedPropagation 支持当前事务，如果当前事务存在，则执行一个嵌套事务，如果当前没有事务，就新建一个事务, 在A事务里面嵌套B事务, B回滚不影响A, A回滚会让B回滚，A提交会让B提交
// A提交失败会使得A回滚，从而导致B回滚 (本质上还是使用的同一个事务，使用savepoint来控制)
func (m *TransactionManager) withNestedPropagation(ctx context.Context, fn func(ctx context.Context) error) error {
	txCtx, ok := ctx.(*TransactionContext)
	if !ok {
		// 新建一个新的txCtx
		var err error
		txCtx, err = NewTransactionContext(ctx, m.db, nil)
		if err != nil {
			return err
		}
	} else {
		txCtx = txCtx.Derive()
		// 设置savepoint
		err := txCtx.SetSavepoint()
		if err != nil {
			return err
		}
	}
	return m.runWithTransaction(txCtx, fn)
}

// withRequiredPropagation 使用当前事务，如果当前没有事务，就新建一个事务
func (m *TransactionManager) withRequiredPropagation(ctx context.Context, fn func(ctx context.Context) error) error {
	txCtx, ok := ctx.(*TransactionContext)
	if !ok {
		// 新建一个新的txCtx
		var err error
		txCtx, err = NewTransactionContext(ctx, m.db, nil)
		if err != nil {
			return err
		}
	} else {
		txCtx = txCtx.Derive()
	}
	return m.runWithTransaction(txCtx, fn)
}

// withRequiresNewPropagation 新建事务，如果当前存在事务，把当前事务挂起，两个事务互不影响
func (m *TransactionManager) withRequiresNewPropagation(ctx context.Context, fn func(ctx context.Context) error) error {
	// 如果当前在事务中，提取原始 context，挂起当前事务
	if txCtx, ok := ctx.(*TransactionContext); ok {
		ctx = txCtx.Context()
	}
	newTxCtx, err := NewTransactionContext(ctx, m.db, nil)
	if err != nil {
		return err
	}
	return m.runWithTransaction(newTxCtx, fn)
}

// runWithTransaction 事务中执行
func (m *TransactionManager) runWithTransaction(txCtx *TransactionContext, fn func(ctx context.Context) error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("transaction exec panic: %v", p)
			buf := make([]byte, 64<<10) //nolint:gomnd
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			logx.Errorf(txCtx, "事务执行panic :%v:\n%s\n", p, buf)
			rErr := txCtx.Rollback()
			if rErr != nil {
				logx.Errorf(txCtx, "回滚事务失败: %v", rErr)
			}
		}
	}()

	if err = fn(txCtx); err != nil {
		if rErr := txCtx.Rollback(); rErr != nil {
			logx.Errorf(txCtx, "回滚事务失败: %v", rErr)
		}
		return err
	}

	if err = txCtx.Commit(); err != nil {
		// commit 失败后尝试回滚，释放锁等资源
		if rErr := txCtx.Rollback(); rErr != nil {
			logx.Errorf(txCtx, "commit失败后回滚也失败: %v", rErr)
		}
		return fmt.Errorf("commit transaction fail: %w", err)
	}

	return
}

// RunAfterTxCommit 在事务提交后执行，用于事务中的异步操作，保证如果事务失败，不执行异步操作。
// hook 注册在当前 TransactionContext 上，子事务成功时 hooks 会逐级提升到父级，
// 最终由 root 事务 commit 成功后统一执行。如果子事务回滚，其注册的 hooks 会被清理，不会被父事务触发。
func RunAfterTxCommit(ctx context.Context, fn func(ctx context.Context) error) error {
	txCtx, ok := ctx.(*TransactionContext)
	// 不在事务中，直接执行
	if !ok {
		return fn(ctx)
	}

	txCtx.AddAfterCommitHook(func(ctx context.Context) error {
		return fn(ctx)
	})
	return nil
}
