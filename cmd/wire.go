//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"

	"purchase/adapter/handler/rpc"
	"purchase/app"
	"purchase/app/assembler"
	"purchase/app/event_handler"
	"purchase/cmd/server"
	"purchase/domain/factory"
	"purchase/domain/service"
	"purchase/infra/acl"
	"purchase/infra/async_task"
	"purchase/infra/config"
	"purchase/infra/dlock"
	"purchase/infra/idempotent"
	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/infra/mq/kafka_sa"
	"purchase/infra/persistence/convertor"
	"purchase/infra/persistence/dal"
	"purchase/infra/persistence/data"
	"purchase/infra/persistence/repo_impl"
	"purchase/infra/persistence/tx"
	"purchase/infra/request"
)

// initApp init kratos application.
func initApp() (*App, func(), error) {
	panic(wire.Build(
		service.ProviderSet,
		repo_impl.ProviderSet,
		server.ProviderSet,
		data.ProviderSet,
		acl.ProviderSet,
		// service.ProviderSet,
		app.ProviderSet,
		rpc.ProviderSet,
		// async_queue.ProviderSet,
		mq.ProviderSet,
		kafka_sa.ProviderSet,
		tx.ProviderSet,
		assembler.ProviderSet,
		dal.ProviderSet,
		config.ProviderSet,
		logx.ProviderSet,
		async_task.ProviderSet,
		dlock.ProviderSet,
		request.ProviderSet,
		idempotent.ProviderSet,
		factory.ProviderSet,
		convertor.ProviderSet,
		event_handler.ProviderSet,
		newApp,
	))
}
