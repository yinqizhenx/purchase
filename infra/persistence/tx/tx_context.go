package tx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"purchase/infra/persistence/dal/db/ent"
)

var ErrNotInTransaction = errors.New("not in a transaction")

type TransactionContext struct {
	ctx       context.Context
	tx        *ent.Tx
	parent    *TransactionContext
	savePoint string
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

// func (c *TransactionContext) InTransaction() bool {
// 	return c.tx != nil
// }

// Derive 同一个事务衍生新的context，但是使用同一个事务，新的context有自己独立的savepoint
func (c *TransactionContext) Derive() *TransactionContext {
	return &TransactionContext{
		ctx:    c.ctx,
		tx:     c.tx,
		parent: c,
	}
}

func (c *TransactionContext) SetParent(p *TransactionContext) {
	c.parent = p
}

func (c *TransactionContext) existSavepoint() bool {
	return c.savePoint != ""
}

func (c *TransactionContext) SetSavepoint() error {
	s, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	_, err = c.tx.ExecContext(c.ctx, fmt.Sprintf("SAVEPOINT %s", s.String()))
	if err != nil {
		return err
	}
	c.savePoint = s.String()
	return nil
}

// func (c *TransactionContext) releaseSavepoint() error {
// 	if !c.InTransaction() {
// 		return ErrNotInTransaction
// 	}
// 	if !c.existSavepoint() {
// 		return nil
// 	}
// 	_, err := c.tx.ExecContext(c.ctx, fmt.Sprintf("RELEASE SAVEPOINT %s", c.savePoint))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (c *TransactionContext) Commit() error {
	if c.IsRoot() { // 只有顶层的 TransactionContext 才会真正提交事务
		err := c.tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *TransactionContext) Rollback() error {
	if !c.existSavepoint() {
		err := c.tx.Rollback()
		if err != nil {
			return err
		}
	} else {
		_, err := c.tx.ExecContext(c.ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", c.savePoint))
		if err != nil {
			return err
		}
		// c.savePoint = ""
	}
	return nil
}
