package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func NewCEPClientForTest(baseUrl string, client *http.Client) *CEPClient {
	return &CEPClient{
		baseURL:    baseUrl,
		httpClient: client,
	}
}

func TestGetCidadeByCep_Sucesso(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")                          // Resposta em JSON
		w.WriteHeader(http.StatusOK)                                                // Retorna o status HTTP 200 OK
		w.Write([]byte(`{"cep":"13183310", "localidade":"Hortolândia","uf":"SP"}`)) // Retorna o corpo (payload)
	}))
	defer mockServer.Close() // Gararante que o servidor falso criado na memória será desligado

	cliente := NewCEPClientForTest(mockServer.URL, mockServer.Client())

	resultadoChannel := cliente.GetCidadeByCep("13183310")
	select {

	case resposta := <-resultadoChannel:

		if resposta.Err != nil {
			t.Fatalf("esperava err nil, mas recebeu o erro: %v", resposta.Err)
		}

		if resposta.Cidade != "Hortolândia" {
			t.Errorf("esperava cidade 'Hortolândia', recebeu: %s", resposta.Cidade)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout: a goroutine demorou muito para responder no channel")
	}
}
