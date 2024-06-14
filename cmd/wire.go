//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"

	"purchase/infra/mq/kafka"

	"purchase/adapter/controller/rpc"

	"purchase/app/app_service"
	"purchase/app/assembler"
	"purchase/domain/service"
	"purchase/infra/acl"
	"purchase/infra/async_task"
	"purchase/infra/config"
	"purchase/infra/dlock"
	"purchase/infra/idempotent"
	"purchase/infra/logx"
	"purchase/infra/mq"
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
		app_service.NewPurchaseService,
		rpc.ProviderSet,
		data.ProviderSet,
		acl.ProviderSet,
		// service.ProviderSet,
		app_service.ProviderSet,
		// async_queue.ProviderSet,
		mq.ProviderSet,
		kafka.ProviderSet,
		tx.ProviderSet,
		assembler.ProviderSet,
		dal.ProviderSet,
		config.ProviderSet,
		logx.ProviderSet,
		async_task.ProviderSet,
		dlock.ProviderSet,
		request.ProviderSet,
		idempotent.ProviderSet,
		newApp,
	))
}
