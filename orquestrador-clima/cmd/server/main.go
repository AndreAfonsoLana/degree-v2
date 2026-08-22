package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/infra/service"
	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/infra/web"
	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/internal/telemetry"
	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/internal/usecase"
	"github.com/joho/godotenv"
)

// Helper para ler variável com valor padrão (fallback)
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	tp, err := telemetry.InitTracer("orquestrador-clima", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		fmt.Printf("Erro ao iniciar o tracer: %v\n", err)
	}
	defer tp.Shutdown(context.Background())

	fmt.Println("Starting main orquestrador-clima...")

	_ = godotenv.Load()

	URL_WEATHER := getEnv("URL_WEATHER", "https://api.weatherapi.com")
	TEMPERATURA_API_KEY := getEnv("WEATHER_API_KEY", "e3af607e22e94e6e828115439260508")
	URL_CEP := getEnv("URL_CEP", "https://viacep.com.br")
	portaServer := getEnv("PORTA_SERVER", "8081")

	fmt.Printf("Configurações carregadas: URL_WEATHER=%s, PORTA=%s\n", URL_WEATHER, portaServer)

	temperaturaService := service.NewTemperaturaService(URL_WEATHER, TEMPERATURA_API_KEY)

	temperaturaClimaUseCase := usecase.NewConsultaClimaUseCase(
		service.NewCidadeService(URL_CEP),
		temperaturaService,
	)

	temperaturaHandler := web.NewTemperaturaHandler(temperaturaClimaUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/temperatura", temperaturaHandler.HandleTemperatura)

	fmt.Printf("Server is running on port %s ....\n", portaServer)

	if err := http.ListenAndServe(":"+portaServer, mux); err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		os.Exit(1)
	}
}
