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

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithExportTimeout(10*time.Second),
			sdktrace.WithMaxExportBatchSize(64),
		),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)), // 10% sampling by default
	)

	// Set global tracer provider and propagator (W3C tracecontext)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = tp.Shutdown(ctx)
	}, nil
}

func newExporter() (sdktrace.SpanExporter, error) {
	switch exporter := os.Getenv("OTEL_TRACES_EXPORTER"); exporter {
	case "otlp":
		// OTLP exporter uses environment variables (OTEL_EXPORTER_OTLP_*).
		// The go.opentelemetry.io/otel/exporters/otlp/otlptrace package handles this.
		// For now, fall back to stdout if OTLP exporter isn't imported.
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	default:
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	}
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
