package data

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"purchase/infra/persistence/dal/db/ent"
	"purchase/pkg/utils"
)

type dbConfig struct {
	Driver string `json:"driver"`
	Source string `json:"source"`
}

func NewDB(c config.Config) (*ent.Client, error) {
	conf := &dbConfig{}
	err := utils.DecodeConfig(c, "data.database", conf)
	if err != nil {
		return nil, err
	}
	drv, err := sql.Open(
		conf.Driver,
		conf.Source,
	)
	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		return nil, err
	}
	sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		// slog.InfoContext(ctx, "ent sql", i...)
		tracer := otel.Tracer("ent.")
		kind := trace.SpanKindServer
		_, span := tracer.Start(ctx,
			"Query",
			trace.WithAttributes(
				attribute.String("sql", fmt.Sprint(i...)),
			),
			trace.WithSpanKind(kind),
		)
		span.End()
	})
	client := ent.NewClient(ent.Driver(sqlDrv))
	// add runtime hooks
	// client.Use(removeCache(cache))
	// Run the auto migration tool.
	if err = client.Schema.Create(context.Background()); err != nil {
		log.Errorf("failed creating schema resources: %v", err)
		return nil, err
	}
	return client, nil
}
