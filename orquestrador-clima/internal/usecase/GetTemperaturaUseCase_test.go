package usecase

import (
	"testing"

	"os"

	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/infra/service"
	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/internal/usecase/dto"
	"github.com/joho/godotenv"
)

func TestGetTemperatura(t *testing.T) {
	cepTeste := "13015904" // Marcer o CEP de teste

	err := godotenv.Load("../../.env")

	if err != nil {
		t.Fatalf("Erro ao carregar .env: %v", err)
	}

	err = godotenv.Load("../../.env.teste")

	if err != nil {
		t.Fatalf("Erro ao carregar .env: %v", err)
	}

	URL_TEMPERATURA := os.Getenv("URL_WEATHER")
	//TESTE_LANA := os.Getenv("TESTE")

	URL_CEP := os.Getenv("URL_CEP")
	TEMPERATURA_API_KEY := os.Getenv("WEATHER_API_KEY")

	//t.Logf("URL_WEATHER: %s | URL_CEP: %s", URL_TEMPERATURA, URL_CEP)
	//t.Logf("TESTE_LANA: %s", TESTE_LANA)
	//t.Logf("TEMPERATURA_API_KEY: %s ", TEMPERATURA_API_KEY)

	cepService := service.NewCidadeService(URL_CEP)
	tempService := service.NewTemperaturaService(URL_TEMPERATURA, TEMPERATURA_API_KEY)

	usecaseNew := NewConsultaClimaUseCase(
		cepService,
		tempService,
	)

	resultado, erro := usecaseNew.ConsultarClimaPorCEP(cepTeste)

	if erro != nil {
		t.Errorf("Erro ao consultar clima: %v", erro)
	}
	if resultado == (dto.GetTemperaturaOutputDTO{}) {
		t.Errorf("Resultado vazio")
	}

	if resultado != (dto.GetTemperaturaOutputDTO{}) {
		t.Logf("SUCESSO! <-> Resultado: %+v\n", resultado)
	}
}
