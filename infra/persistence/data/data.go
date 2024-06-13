package data

import (
	"context"
	"fmt"

	"github.com/google/wire"

	"purchase/infra/persistence/dal/db/ent"
	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewRedis)

// // Data .
// type Data struct {
// 	db  *ent.Client
// 	rdb *redis.Client
// 	mgo *mongo.Database
// 	// cfg *conf.DataConfig
// 	log log.Logger
// }

// NewData .
// func NewData(c config.Config) (*Data, func(), error) {
// 	dataCfg := c.Value("data")
// 	conf := &cfg.DataConfig{}
// 	err := dataCfg.Scan(conf)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	d := &Data{conf: conf}
// 	err = d.initDB()
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	err = d.initRedis()
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	return d, func() {
// 		log.Info("closing the data resources")
// 		if err := d.db.Close(); err != nil {
// 			log.Errorf("d.db.Close error :%s", err)
// 		}
// 		if err := d.rdb.Close(); err != nil {
// 			log.Errorf("d.rdb.Close error :%s", err)
// 		}
// 	}, nil
// }
//
// func (d *Data) initDB() error {
// 	drv, err := sql.Open(
// 		d.conf.Database.Driver,
// 		d.conf.Database.Source,
// 	)
// 	if err != nil {
// 		log.Errorf("failed opening connection to sqlite: %v", err)
// 		return err
// 	}
// 	sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
// 		log.Info(i...)
// 		tracer := otel.Tracer("ent.")
// 		kind := trace.SpanKindServer
// 		_, span := tracer.Start(ctx,
// 			"Query",
// 			trace.WithAttributes(
// 				attribute.String("sql", fmt.Sprint(i...)),
// 			),
// 			trace.WithSpanKind(kind),
// 		)
// 		span.End()
// 	})
// 	client := ent.NewClient(ent.Driver(sqlDrv))
// 	// add runtime hooks
// 	// client.Use(removeCache(cache))
// 	// Run the auto migration tool.
// 	if err := client.Schema.Create(context.Background()); err != nil {
// 		log.Errorf("failed creating schema resources: %v", err)
// 		return err
// 	}
// 	d.db = client
// 	return nil
// }
//
// func (d *Data) initRedis() error {
// 	rdbCfg := d.conf.Redis
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:         rdbCfg.Addr,
// 		Password:     rdbCfg.Password,
// 		DB:           int(rdbCfg.Db),
// 		DialTimeout:  utils.Float2Duration(rdbCfg.DialTimeout),
// 		WriteTimeout: utils.Float2Duration(rdbCfg.WriteTimeout),
// 		ReadTimeout:  utils.Float2Duration(rdbCfg.ReadTimeout),
// 	})
// 	rdb.AddHook(redisotel.TracingHook{})
// 	d.rdb = rdb
// 	return nil
// }

func WithOneTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

type entTxKey struct{}

// WithTx 保证多层事务嵌套时，使用的是同一个事务
func WithTx(ctx context.Context, client *ent.Client, fn func(ctx context.Context, tx *ent.Tx) error) error {
	tx, ok := ctx.Value(entTxKey{}).(*ent.Tx)
	txCtx := ctx

	if !ok {
		var err error
		tx, err = client.Tx(ctx)
		if err != nil {
			return err
		}
		txCtx = context.WithValue(ctx, entTxKey{}, tx)
	}

	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(txCtx, tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// RunAfterTxCommit 在事务提交后执行，用于事物中的异步操作，保证如果事务失败，不执行异步操作
func RunAfterTxCommit(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, ok := ctx.Value(entTxKey{}).(*ent.Tx)
	// 不在事务中，直接执行
	if !ok {
		return fn(ctx)
	}

	hook := func(next ent.Committer) ent.Committer {
		return ent.CommitFunc(func(ctx context.Context, tx *ent.Tx) error {
			err := next.Commit(ctx, tx)
			if err != nil {
				return err
			}
			return fn(ctx)
		})
	}

	tx.OnCommit(hook)
	return nil
}
