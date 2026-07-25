// Package tracing provides OpenTelemetry tracing setup for the Muse backend.
//
// Usage:
//
//	import "music/internal/platform/tracing"
//
//	shutdown, err := tracing.Init(ctx, tracing.Config{
//	    ServiceName: "music-api",
//	    Environment: "production",
//	})
//	if err != nil { ... }
//	defer shutdown()
//
// The tracer exports spans to stdout by default. To send to an OTLP collector,
// set the standard OTEL environment variables:
//
//	OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4318
//	OTEL_EXPORTER_OTLP_HEADERS=api-key=xxx
package tracing

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Config holds the minimal configuration to bootstrap the tracer provider.
type Config struct {
	ServiceName string
	Environment string
}

// Init sets up the global OpenTelemetry tracer provider with a stdout exporter
// and W3C tracecontext propagation. Returns a shutdown function that must be
// called on process exit to flush remaining spans.
//
// The exporter can be overridden via the OTEL_TRACES_EXPORTER env var:
//   - "stdout" (default): write spans to stdout as JSON
//   - "otlp": export via OTLP (requires OTEL_EXPORTER_OTLP_ENDPOINT)
//
// If the environment variable is empty or unset, "stdout" is used.
func Init(ctx context.Context, cfg Config) (func(), error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("deployment.environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("tracing: create resource: %w", err)
	}

	exporter, err := newExporter()
	if err != nil {
		return nil, fmt.Errorf("tracing: create exporter: %w", err)
	}

	var shutdown func()
	if exporter != nil {
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
			sdktrace.WithBatcher(exporter,
				sdktrace.WithBatchTimeout(5*time.Second),
				sdktrace.WithExportTimeout(10*time.Second),
				sdktrace.WithMaxExportBatchSize(64),
			),
			sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)), // 10% sampling by default
		)
		otel.SetTracerProvider(tp)
		shutdown = func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = tp.Shutdown(ctx)
		}
	} else {
		shutdown = func() {}
	}

	// Set global tracer provider and propagator (W3C tracecontext)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return shutdown, nil
}

func newExporter() (sdktrace.SpanExporter, error) {
	switch exporter := os.Getenv("OTEL_TRACES_EXPORTER"); exporter {
	case "otlp":
		return newOTLPExporter()
	case "stdout":
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	default:
		return nil, nil
	}
}

func newOTLPExporter() (sdktrace.SpanExporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithTimeout(10 * time.Second),
	}

	// Env-based endpoint (OTEL_EXPORTER_OTLP_ENDPOINT) is the default,
	// but allow explicit override via OTEL_EXPORTER_OTLP_TRACES_ENDPOINT.
	if e := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"); e != "" {
		opts = append(opts, otlptracehttp.WithEndpointURL(e))
	}

	return otlptracehttp.New(context.Background(), opts...)
}

// Tracer returns the named tracer for the given package/component name.
// This is the primary way to create spans in business logic.
//
//	func MyHandler(ctx context.Context) {
//	    ctx, span := tracing.Tracer("my-module").Start(ctx, "MyHandler")
//	    defer span.End()
//	    // ... do work
//	}
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// DBSpan wraps a database operation in a tracing span. Use it to trace
// individual SQL queries or batches of queries.
//
//	rows, err := tracing.DBSpan(ctx, "db", "SELECT tracks", func(ctx context.Context) ([]Track, error) {
//	    return repo.FindAll(ctx)
//	})
func DBSpan[T any](ctx context.Context, tracerName, spanName string, fn func(context.Context) (T, error)) (T, error) {
	_, span := Tracer(tracerName).Start(ctx, spanName)
	defer span.End()

	result, err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return result, err
}

// DBSpanNoResult wraps a database operation that returns no value in a tracing span.
func DBSpanNoResult(ctx context.Context, tracerName, spanName string, fn func(context.Context) error) error {
	_, span := Tracer(tracerName).Start(ctx, spanName)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// ConfigureFromEnv reads standard OTEL environment variables and returns a Config.
// Supported variables:
//   - OTEL_SERVICE_NAME
//   - OTEL_EXPORTER_OTLP_ENDPOINT
//   - OTEL_TRACES_EXPORTER (stdout | otlp)
//   - OTEL_RESOURCE_ATTRIBUTES
func ConfigureFromEnv() Config {
	svc := os.Getenv("OTEL_SERVICE_NAME")
	if svc == "" {
		svc = "muse-api"
	}
	env := os.Getenv("OTEL_RESOURCE_ATTRIBUTES")
	if env == "" {
		env = "development"
	}
	return Config{
		ServiceName: svc,
		Environment: env,
	}
}
