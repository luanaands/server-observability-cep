package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/luanaands/server-validation-cep/configs"
	_ "github.com/luanaands/server-validation-cep/docs"
	"github.com/luanaands/server-validation-cep/internal/dto"
	"github.com/luanaands/server-validation-cep/internal/infra/service"
	"github.com/luanaands/server-validation-cep/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// @title Desafio CEP API - golang
// @version 1.0
// @description API para consulta do tempo real de um CEP
// @termsOfService http://swagger.io/terms/

// @contact.name Luana Andrade
// @contact.email luanaands@gmail.com

// @host localhost:8082
// @schemes https
// @basePath /
func main() {
	configs, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider(configs.OtelServiceName, configs.OtelExporterOtlpEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	tracer := otel.Tracer("service-tracer")

	templateData := &dto.TemplateData{
		OTELTracer:      tracer,
		Title:           configs.Title,
		ExternalCallURL: configs.MyCoreHost,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("MyCoreHost", configs.MyCoreHost))

	var cepService service.CepDetailsInterface = service.NewCepDetailsService()
	handlerCep := handlers.NewCepHandler(cepService, templateData)

	r.Post("/cep", handlerCep.GetCep)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8082/docs/doc.json")))

	go func() {
		log.Println("Server is running on port 8082")
		if err = http.ListenAndServe(":8082", r); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason")
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}

func initProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
	ctx := context.Background()

	res, erro := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if erro != nil {
		return nil, fmt.Errorf("failed to create resource: %w", erro)
	}

	// Aumentar timeout para permitir conexão com o collector
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	normalizedCollectorURL := strings.TrimPrefix(strings.TrimPrefix(strings.TrimSpace(collectorURL), "http://"), "https://")

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(normalizedCollectorURL),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		// Log warning mas não falha - permite que a app rode sem OTEL
		log.Printf("Warning: failed to create OTLP exporter: %v. Traces will not be exported.\n", err)
		// Retorna um no-op shutdown function
		return func(context.Context) error { return nil }, nil
	}

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tracerProvider.Shutdown, nil
}
