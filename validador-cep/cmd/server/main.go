package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/AndreAfonsoLana/go-degree-validador-cep/infra/service"
	"github.com/AndreAfonsoLana/go-degree-validador-cep/infra/web"
	"github.com/AndreAfonsoLana/go-degree-validador-cep/internal/telemetry"
	"github.com/AndreAfonsoLana/go-degree-validador-cep/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	tp, err := telemetry.InitTracer("validador-cep", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		fmt.Printf("Erro ao iniciar o tracer: %v\n", err)
	}
	defer tp.Shutdown(context.Background())

	fmt.Println("Starting main validador-cep...")

	errou := godotenv.Load(".env")

	if errou != nil {
		fmt.Printf("Erro ao carregar .env: %v\n", errou)
	}

	URL_TEMPERATURA_API := os.Getenv("URL_API_TEMPERATURA")

	temperaturaService := service.NewGetTempeturaService(URL_TEMPERATURA_API)

	temperaturaUseCase := usecase.NewGetTemperaturaUsecase(temperaturaService)

	temperaturaHandler := web.NewTemperaturaHandler(temperaturaUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/temperatura", temperaturaHandler.HandleTemperatura)

	fmt.Printf("Server is running on port 8080 ...")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
