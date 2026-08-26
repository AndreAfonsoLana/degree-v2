package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"github.com/AndreAfonsoLana/go-degree-validador-cep/infra/service/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type CEPClient struct {
	baseURL    string
	httpClient *http.Client
}
type TemperaturaResult struct {
	City   string
	Temp_C float32
	Temp_F float32
	Temp_K float32
	Err    error
}

func NewGetTempeturaService(baseURL string) *CEPClient {
	return &CEPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}
func (t *CEPClient) GetTemperaturaByCep(ctx context.Context, cep string) <-chan TemperaturaResult {
	channelResultado := make(chan TemperaturaResult, 1)

	go func() {
		//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		reqCtx, cancel := context.WithTimeout(ctx, 20*time.Second)

		defer cancel() // Garante que o contexto será cancelado após o uso
		defer close(channelResultado)

		url := fmt.Sprintf("%s/temperatura?cep=%s", t.baseURL, cep)
		fmt.Printf(" URL SEND REQUEST %s\n", url)

		requisicao, erro := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)

		if erro != nil {
			channelResultado <- TemperaturaResult{Err: erro}
			return
		}

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(requisicao.Header))

		resposta, erro := t.httpClient.Do(requisicao)
		// Envia a requisição

		if erro != nil {
			channelResultado <- TemperaturaResult{Err: erro}
			return
		}
		defer resposta.Body.Close()

		var payload dto.GetTemperaturaResponse
		if erro := json.NewDecoder(resposta.Body).Decode(&payload); erro != nil {
			channelResultado <- TemperaturaResult{
				Err: erro,
			}
			return
		}
		if payload.Cidade == "" {
			channelResultado <- TemperaturaResult{Err: fmt.Errorf("can not find zipcode")}
			return
		}

		channelResultado <- TemperaturaResult{
			City:   payload.Cidade,
			Temp_C: float32(payload.TempC),
			Temp_F: float32(payload.TempF),
			Temp_K: float32(payload.TempK),
			Err:    nil,
		}
	}()

	return channelResultado
}
