package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"

	"purchase/adapter/consumer"
	"purchase/infra/async_task"
	"purchase/infra/dtx"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
)

// var Provider = wire.NewSet(newApp)

type App struct {
	*kratos.App
	// asyncq   *asyncq.AsyncQueueServer
	// eventApp *event.DomainEventApp
}

func newApp(logger log.Logger, gs *kgrpc.Server, hs *khttp.Server, ams *async_task.AsyncTaskMux, ec *consumer.EventConsumer, dtm *dtx.DistributeTxManager) *App {
	a := &App{
		// asyncq:   aq,
		// eventApp: eventApp,
	}
	options := []kratos.Option{
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
			ams,
			ec,
			dtm,
		),
	}
	// if a.asyncq != nil {
	// 	options = append(options, kratos.BeforeStart(func(ctx context.Context) error {
	// 		return a.asyncq.Run()
	// 	}))
	// }
	// if a.eventApp != nil {
	// 	options = append(options, kratos.BeforeStart(func(ctx context.Context) error {
	// 		go a.eventApp.Start(ctx)
	// 		return nil
	// 	}))
	// }

	kApp := kratos.New(options...)
	a.App = kApp
	return a
}

func (a *App) Run() {
	if err := a.App.Run(); err != nil {
		panic(err)
	}
}
