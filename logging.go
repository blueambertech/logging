package logging

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer
var traceProvider *sdktrace.TracerProvider

func Setup(ctx context.Context, svcName string) {
	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	traceProvider = newTraceProvider(exp, svcName)

	otel.SetTracerProvider(traceProvider)
	Tracer = traceProvider.Tracer(svcName)
}

func DeferredCleanup(ctx context.Context) {
	_ = traceProvider.Shutdown(ctx)
}

func newExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return stdouttrace.New()
}

func newTraceProvider(exp sdktrace.SpanExporter, svcName string) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(svcName),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
