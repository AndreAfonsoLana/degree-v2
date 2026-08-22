package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer(serviceName, endpoint string) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	if endpoint == "" {
		endpoint = "otel-collector:4317"
	}

	// Tenta criar o exportador gRPC
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		fmt.Printf("❌ ERRO AO CRIAR OTLP EXPORTER para %s: %v\n", endpoint, err)
		return nil, err
	}

	// ADICIONE ESTA VERIFICAÇÃO PARA FORÇAR O TESTE DE CONEXÃO:
	// O Start costuma disparar a tentativa real de conexão se configurado,
	// mas vamos ver se o erro estoura aqui ou nos logs da aplicação.

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	fmt.Printf("✅ Tracer inicializado com sucesso para o serviço: %s\n", serviceName)
	return tp, nil
}
