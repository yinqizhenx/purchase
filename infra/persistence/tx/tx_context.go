package tx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"purchase/infra/persistence/dal/db/ent"
)

// AfterCommitHook 事务提交后执行的回调函数
type AfterCommitHook func(ctx context.Context) error

type TransactionContext struct {
	ctx              context.Context
	tx               *ent.Tx
	parent           *TransactionContext
	savePoint        string
	afterCommitHooks []AfterCommitHook // 自管理的 commit 后回调列表
}

func NewTransactionContext(ctx context.Context, cli *ent.Client, parent *TransactionContext) (*TransactionContext, error) {
	tx, err := cli.Tx(ctx)
	if err != nil {
		return nil, err
	}
	return &TransactionContext{
		ctx:    ctx,
		tx:     tx,
		parent: parent,
	}, nil
}

func (c *TransactionContext) Tx() *ent.Tx {
	return c.tx
}

func (c *TransactionContext) Context() context.Context {
	return c.ctx
}

func (c *TransactionContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *TransactionContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *TransactionContext) Err() error {
	return c.ctx.Err()
}

func (c *TransactionContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *TransactionContext) IsRoot() bool {
	return c.parent == nil
}

// root 沿 parent 链找到最顶层的 TransactionContext
func (c *TransactionContext) root() *TransactionContext {
	root := c
	for root.parent != nil {
		root = root.parent
	}
	return root
}

// Derive 同一个事务衍生新的context，但是使用同一个事务，新的context有自己独立的savepoint
func (c *TransactionContext) Derive() *TransactionContext {
	return &TransactionContext{
		ctx:    c.ctx,
		tx:     c.tx,
		parent: c,
	}
}

func (c *TransactionContext) existSavepoint() bool {
	return c.savePoint != ""
}

func (c *TransactionContext) SetSavepoint() error {
	s, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	// 将 UUID 中的连字符替换为下划线，并加上 sp_ 前缀，避免 SQL 标识符解析问题
	name := fmt.Sprintf("sp_%s", strings.ReplaceAll(s.String(), "-", "_"))
	_, err = c.tx.ExecContext(c.ctx, fmt.Sprintf("SAVEPOINT %s", name))
	if err != nil {
		return err
	}
	c.savePoint = name
	return nil
}

// AddAfterCommitHook 注册事务提交后的回调。hook 会被提升到 root TransactionContext 上，
// 只在最外层事务真正 commit 成功后才执行。如果当前子事务回滚，会通过 clearAfterCommitHooks 清理。
func (c *TransactionContext) AddAfterCommitHook(fn AfterCommitHook) {
	c.afterCommitHooks = append(c.afterCommitHooks, fn)
}

// promoteHooksToParent 将当前 context 的 hooks 提升到父 context。
// 在子事务成功时调用，确保 hooks 沿着事务链向上传递，最终由 root 事务提交时执行。
func (c *TransactionContext) promoteHooksToParent() {
	if c.parent != nil && len(c.afterCommitHooks) > 0 {
		c.parent.afterCommitHooks = append(c.parent.afterCommitHooks, c.afterCommitHooks...)
		c.afterCommitHooks = nil
	}
}

// clearAfterCommitHooks 清理当前 context 上注册的所有 hooks。
// 在子事务回滚时调用，确保已回滚的子事务的 hooks 不会被父事务触发。
func (c *TransactionContext) clearAfterCommitHooks() {
	c.afterCommitHooks = nil
}

// executeAfterCommitHooks 执行所有注册的 commit 后回调
func (c *TransactionContext) executeAfterCommitHooks() error {
	for _, hook := range c.afterCommitHooks {
		if err := hook(c.ctx); err != nil {
			return err
		}
	}
	c.afterCommitHooks = nil
	return nil
}

func (c *TransactionContext) Commit() error {
	if c.IsRoot() { // 只有顶层的 TransactionContext 才会真正提交事务
		err := c.tx.Commit()
		if err != nil {
			return err
		}
		// 提交成功后执行所有 hooks
		return c.executeAfterCommitHooks()
	}
	// 非 root：将 hooks 提升到父级，等待 root 提交时统一执行
	c.promoteHooksToParent()
	return nil
}

// Rollback 回滚事务。使用独立的 context 避免原 context 超时导致回滚失败
func (c *TransactionContext) Rollback() error {
	// 回滚时清理当前 context 的 hooks，防止已回滚子事务的 hooks 被父事务触发
	c.clearAfterCommitHooks()

	// 使用 WithoutCancel 确保即使原 context 已超时/取消，回滚操作仍能执行
	rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(c.ctx), 5*time.Second)
	defer cancel()

	if !c.existSavepoint() {
		return c.tx.Rollback()
	}
	_, err := c.tx.ExecContext(rollbackCtx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", c.savePoint))
	return err
}
