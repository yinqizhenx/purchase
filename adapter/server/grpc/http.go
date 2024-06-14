package grpc

import (
	"github.com/go-kratos/kratos/v2/transport/http/pprof"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"purchase/pkg/utils"
)

func NewHttpServer(c config.Config) *http.Server {
	type Conf struct {
		Address      string `yaml:"address"`
		StartMetric  bool   `yaml:"start_metric"`
		StartProfile bool   `yaml:"start_profile"`
	}
	conf := &Conf{}
	utils.DecodeConfigX(c, "server.http", conf)
	httpAddr := conf.Address
	if httpAddr == "" {
		httpAddr = ":18080"
	}
	httpSrv := http.NewServer(
		http.Address(httpAddr),
	)
	if conf.StartMetric {
		httpSrv.Handle("/metrics", promhttp.Handler())
	}
	if conf.StartProfile {
		httpSrv.HandlePrefix("/", pprof.NewHandler())
	}
	return httpSrv
}
