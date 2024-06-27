// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"purchase/adapter/scheduler"
	"purchase/app"
	"purchase/app/assembler"
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
	"purchase/infra/mq/kafka"
	"purchase/infra/persistence/dal"
	"purchase/infra/persistence/data"
	"purchase/infra/persistence/repo_impl"
	"purchase/infra/persistence/tx"
	"purchase/infra/request"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp() (*App, func(), error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}
	logger := logx.NewLogger(configConfig)
	client, err := data.NewDB(configConfig)
	if err != nil {
		return nil, nil, err
	}
	paDal := dal.NewPADal(client)
	unboundedChan, cleanup := async_task.NewMessageChan()
	asyncTaskDal := dal.NewAsyncTaskDal(client)
	paymentCenterRepo := repo_impl.NewPARepository(paDal, unboundedChan, asyncTaskDal)
	httpClient, cleanup2 := request.NewHttpClient(logger)
	mdmService := acl.NewMDMService(httpClient)
	eventRepo := repo_impl.NewEventRepository(unboundedChan, asyncTaskDal)
	pcFactory := factory.NewPCFactory(mdmService)
	paDomainService := service.NewPADomainService(paymentCenterRepo, mdmService, eventRepo, pcFactory)
	assemblerAssembler := assembler.NewAssembler()
	transactionManager := tx.NewTransactionManager(client)
	paymentCenterAppService := app.NewPaymentCenterAppService(paDomainService, paymentCenterRepo, assemblerAssembler, transactionManager)
	suDal := dal.NewSUDal(client)
	idGenFunc := mq.NewIDGenFunc()
	publisher, err := kafka.NewKafkaPublisher(configConfig, idGenFunc)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	suRepo := repo_impl.NewSURepository(suDal, publisher)
	suDomainService := service.NewSUDomainService(suRepo)
	suAppService := app.NewSuAppService(suDomainService, suRepo)
	appService := app.NewPurchaseService(paymentCenterAppService, suAppService)
	grpcServer := server.NewGRPCServer(configConfig, logger, appService)
	httpServer := server.NewHttpServer(configConfig)
	redisClient, err := data.NewRedis(configConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	lockBuilder := dlock.NewRedisLock(redisClient)
	asyncTaskMux := scheduler.NewAsyncTaskServer(publisher, asyncTaskDal, transactionManager, unboundedChan, lockBuilder)
	idempotentIdempotent := idempotent.NewIdempotentImpl(redisClient)
	subscriber, err := kafka.NewKafkaSubscriber(configConfig, idempotentIdempotent, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	domainEventAppService := app.NewDomainEventAppService()
	eventConsumer := server.NewEventConsumerServer(subscriber, domainEventAppService)
	mainApp := newApp(logger, grpcServer, httpServer, asyncTaskMux, eventConsumer)
	return mainApp, func() {
		cleanup2()
		cleanup()
	}, nil
}
