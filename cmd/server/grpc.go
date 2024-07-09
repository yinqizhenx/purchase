package server

import (
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	"purchase/adapter/handler/rpc"
	pbpc "purchase/idl/payment_center"
	"purchase/infra/metric"
	"purchase/infra/middleware"
	trac "purchase/infra/tracing"
	"purchase/pkg/utils"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c config.Config, svc *rpc.Server) *grpc.Server {
	type GrpcServerConfig struct {
		Network string  `json:"network"`
		Addr    string  `json:"addr"`
		Timeout float64 `json:"timeout"`
	}
	conf := &GrpcServerConfig{}
	err := utils.DecodeConfig(c, "server.grpc", conf)
	if err != nil {
		panic(err)
	}
	tp, err := trac.GetTracerProvider(c)
	if err != nil {
		panic(err)
	}
	opts := []grpc.ServerOption{
		grpc.Middleware(
			middleware.InjectContext(),
			recovery.Recovery(),
			tracing.Server(tracing.WithTracerProvider(tp)),
			middleware.ServerLogging(),
			// validate.Validator(),
			middleware.WithValidator(),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(metric.MetricSeconds)),
				metrics.WithRequests(prom.NewCounter(metric.MetricRequests)),
			),
		),
	}
	if conf.Network != "" {
		opts = append(opts, grpc.Network(conf.Network))
	}
	if conf.Addr != "" {
		opts = append(opts, grpc.Address(conf.Addr))
	}
	if conf.Timeout != 0 {
		opts = append(opts, grpc.Timeout(utils.Float2Duration(conf.Timeout)))
	}
	srv := grpc.NewServer(opts...)
	registerServer(srv, svc)
	return srv
}

// registerServer 注册grpc服务，新加的功能需要在此注册
func registerServer(s ggrpc.ServiceRegistrar, srv *rpc.Server) {
	pbpc.RegisterPaymentCenterServer(s, srv)
}
