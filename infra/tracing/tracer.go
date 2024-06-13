package tracing

import (
	"errors"

	"github.com/go-kratos/kratos/v2/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"

	"purchase/pkg/utils"
)

func GetTracerProvider(c config.Config) (trace.TracerProvider, error) {
	var tracerUrl string
	err := utils.DecodeConfig(c, "trace.endpoint", &tracerUrl)
	if err != nil {
		return nil, err
	}
	if tracerUrl == "" {
		return nil, errors.New("no tracer found in config")
	}

	var appName string
	err = utils.DecodeConfig(c, "app.name", &appName)
	if err != nil {
		return nil, err
	}
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(tracerUrl)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(appName),
			attribute.String("env", "dev"),
		)),
	)
	return tp, nil
}
