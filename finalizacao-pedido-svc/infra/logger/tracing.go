package logger

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
)

var Tracer = otel.Tracer("finalizacao-pedido-svc")

func InitTracer() func() {
	ctx := context.Background()

	// Configurar o exporter OTLP HTTP para o Collector
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("otel-collector:4318"),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithURLPath("/v1/traces"),
	)
	if err != nil {
		Log.Fatal("Erro ao criar exporter OTLP", zap.Error(err))
	}

	// Configurar o recurso com nome do serviço
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(os.Getenv("OPEN_TELEMETRY_NAME")),
		),
	)
	if err != nil {
		Log.Fatal("Erro ao criar resource", zap.Error(err))
	}

	// Configurar o TracerProvider com Batch Span Processor
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	// Configurar o propagador global para extrair e injetar o contexto
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			Log.Error("Erro ao desligar tracer provider", zap.Error(err))
		}
	}
}
