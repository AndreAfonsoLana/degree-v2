package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"net/url"

	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/infra/service/dto"
	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/utils"
	"go.opentelemetry.io/otel"
)

type CidadeClinte struct {
	baseUrl    string
	apiKey     string
	httpClient *http.Client
}
type TemperaturaResultado struct {
	Temp_c float64
	Temp_F float64
	Temp_K float64
	Err    error
}

func NewTemperaturaService(baseURL string, apiKey string) *CidadeClinte {
	return &CidadeClinte{
		baseUrl: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}
func (c *CidadeClinte) GetTemperaturaByCidade(ctx context.Context, cidade string) <-chan TemperaturaResultado {
	channelResultado := make(chan TemperaturaResultado, 1)
	tr := otel.Tracer("orquestrador-clima-usecase")
	ctx, span := tr.Start(ctx, "GetTemperaturaByCidadE")
	defer span.End()

	go func() {
		contexto, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelar()
		defer close(channelResultado)

		cidadeFormata := utils.RemoverAcentos(cidade)
		base := c.baseUrl

		if base == "" {
			base = "https://api.weatherapi.com"
		}

		cidadeEscapada := url.QueryEscape(cidadeFormata)

		endpoint := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s", c.baseUrl, c.apiKey, cidadeEscapada)

		requisicao, erro := http.NewRequestWithContext(contexto, http.MethodGet, endpoint, nil)

		if erro != nil {
			channelResultado <- TemperaturaResultado{Err: erro}
			return
		}
		resposta, erro := c.httpClient.Do(requisicao)

		if erro != nil {
			channelResultado <- TemperaturaResultado{Err: erro}
			return
		}

		var payload dto.TemperaturaResponseDTO

		if error := json.NewDecoder(resposta.Body).Decode(&payload); error != nil {
			channelResultado <- TemperaturaResultado{Err: erro}
			return
		}
		channelResultado <- TemperaturaResultado{
			Temp_c: payload.Current.Temp_c,
			Temp_F: utils.CelsiusParaFahrenheit(payload.Current.Temp_c),
			Temp_K: utils.CelsiusParaKelvin(payload.Current.Temp_c),
			Err:    nil,
		}

		defer resposta.Body.Close()
	}()
	return channelResultado
}
