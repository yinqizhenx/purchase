package request

import (
	"context"
	"net/http"

	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"purchase/infra/metric"
)

var ProviderSet = wire.NewSet(NewHttpClient)

type HttpClient struct {
	*khttp.Client
	// c resty.Client
}

func NewHttpClient(lg log.Logger) (*HttpClient, func()) {
	cli, err := khttp.NewClient(context.Background(), khttp.WithMiddleware(
		tracing.Client(),
		logging.Client(lg),
		metrics.Client(
			metrics.WithSeconds(prom.NewHistogram(metric.MetricSeconds)),
			metrics.WithRequests(prom.NewCounter(metric.MetricRequests)),
		),
	))
	if err != nil {
		panic(err)
	}
	c := &HttpClient{
		Client: cli,
	}
	return c, c.Close
}

// NewRequest resp must be a pointer
func (c *HttpClient) NewRequest(url string, req, resp interface{}) *Request {
	r := &Request{
		client: c,
		url:    url,
		req:    req,
		resp:   resp,
		header: make(map[string][]string),
	}
	return r
}

func (c *HttpClient) Close() {
	c.Client.Close()
}

type Request struct {
	client *HttpClient
	url    string
	header map[string][]string
	req    interface{}
	resp   interface{}
}

func (r *Request) SetHeader(key, val string) {
	r.header[key] = append(r.header[key], val)
}

func (r *Request) Get(ctx context.Context) error {
	header := http.Header(r.header)
	return r.client.Invoke(ctx, "GET", r.url, r.req, r.resp, khttp.Header(&header))
}

func (r *Request) Post(ctx context.Context) error {
	header := http.Header(r.header)
	return r.client.Invoke(ctx, "POST", r.url, r.req, r.resp, khttp.Header(&header))
}
