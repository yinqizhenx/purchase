package idempotent

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"

	"purchase/infra/persistence/dal/db/ent"
	"purchase/infra/persistence/dal/db/ent/idempotent"
	"purchase/infra/persistence/tx"
)

type MysqlIdempotentImpl struct {
	db *ent.Client
}

func NewMysqlIdempotentImpl(db *ent.Client) Idempotent {
	return &MysqlIdempotentImpl{
		db: db,
	}
}

func (imp *MysqlIdempotentImpl) getIdempotentClient(ctx context.Context) *ent.IdempotentClient {
	txCtx, ok := ctx.(*tx.TransactionContext)
	if ok {
		return txCtx.Tx().Idempotent
	}
	return imp.db.Idempotent
}

func (imp *MysqlIdempotentImpl) SetKeyPendingWithDDL(ctx context.Context, key string, _ time.Duration) (bool, error) {
	_, err := imp.getIdempotentClient(ctx).Create().
		SetKey(key).
		SetState(Pending).
		Save(ctx)
	if err != nil {
		if IsMysqlDuplicateErr(err) {
			return false, nil
		}
	}
	return true, nil
}

func (imp *MysqlIdempotentImpl) GetKeyState(ctx context.Context, key string) (string, error) {
	r, err := imp.getIdempotentClient(ctx).Query().Where(idempotent.Key(key)).Only(ctx)
	if err != nil {
		return "", err
	}
	return r.State, nil
}

func (imp *MysqlIdempotentImpl) UpdateKeyDone(ctx context.Context, key string) error {
	return imp.getIdempotentClient(ctx).Update().
		SetState(Done).
		Where(idempotent.Key(key)).
		Exec(ctx)
}

func (imp *MysqlIdempotentImpl) RemoveFailKey(ctx context.Context, key string) error {
	_, err := imp.getIdempotentClient(ctx).Delete().
		Where(idempotent.Key(key)).
		Exec(ctx)
	return err
}

func IsMysqlDuplicateErr(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		// MySQL错误码1062: ER_DUP_ENTRY
		if mysqlErr.Number == 1062 {
			return true
		}
	}
	return false
}
