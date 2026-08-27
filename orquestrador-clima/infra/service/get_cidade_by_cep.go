package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/infra/service/dto"
	"go.opentelemetry.io/otel"
)

type CEPClient struct {
	baseURL    string
	httpClient *http.Client
}

type CidadeResult struct {
	Cidade string
	Err    error
}

func NewCidadeService(baseURL string) *CEPClient {
	return &CEPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (v *CEPClient) GetCidadeByCep(contexto context.Context, cep string) <-chan CidadeResult {
	channelResultado := make(chan CidadeResult, 1)

	tr := otel.Tracer("orquestrador-clima-usecase")
	contexto, span := tr.Start(contexto, "GetCidadeByCep")

	defer span.End()

	go func() {
		ctxTimeout, cancelar := context.WithTimeout(contexto, 10*time.Second)
		defer cancelar() // Fechar contexto após uso
		defer close(channelResultado)

		base := v.baseURL
		if base == "" {
			base = "https://viacep.com.br"
		}

		url := fmt.Sprintf("%s/ws/%s/json/", base, cep)

		requisicao, erro := http.NewRequestWithContext(ctxTimeout, http.MethodGet, url, nil)

		if erro != nil {
			channelResultado <- CidadeResult{Err: erro}
			return
		}
		resposta, erro := v.httpClient.Do(requisicao)

		if erro != nil {
			channelResultado <- CidadeResult{Err: erro}
			return
		}
		defer resposta.Body.Close()

		if resposta.StatusCode == http.StatusNotFound {
			channelResultado <- CidadeResult{Err: fmt.Errorf("can not find zipcode")}
		}

		var payload dto.CepResponseDTO
		if erro := json.NewDecoder(resposta.Body).Decode(&payload); erro != nil {
			channelResultado <- CidadeResult{
				Err: erro,
			}
			return
		}
		if payload.Localidade == "" {
			channelResultado <- CidadeResult{Err: fmt.Errorf("can not find zipcode")}
			return
		}
		channelResultado <- CidadeResult{
			Cidade: payload.Localidade,
			Err:    nil,
		}

	}() // Função que executa automaticamente
	return channelResultado
}
