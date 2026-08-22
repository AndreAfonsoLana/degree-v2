package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Função auxiliar exclusiva para os testes, permitindo injetar o http.Client do mock
func NewGetTempeturaServiceForTest(baseURL string, client *http.Client) *CEPClient {
	return &CEPClient{
		baseURL:    baseURL,
		httpClient: client,
	}
}

func TestGetTemperaturaByCEP_Sucesso(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"cidade": "Hortolândia", "temp_c": 28.5, "temp_f": 83.3, "temp_k": 301.65}`))
	}))
	defer mockServer.Close()

	cliente := NewGetTempeturaServiceForTest(mockServer.URL, mockServer.Client())

	resultadoChannel := cliente.GetTemperaturaByCep("13183310")
	select {
	case resposta := <-resultadoChannel:
		if resposta.Err != nil {
			t.Fatalf("não esperava erro, mas recebeu: %v", resposta.Err)
		}

		t.Logf("SUCESSO! Resultado lido do canal: %+v", resposta)
		if resposta.City != "Hortolândia" {
			t.Errorf("esperava cidade 'Hortolândia', mas recebeu: %s", resposta.City)
		}

		if resposta.Temp_C != 28.5 {
			t.Errorf("esperava temp_C 28.5, mas recebeu: %.1f", resposta.Temp_C)
		}

	case <-time.After(2 * time.Second):
		t.Fatal("timeout: a goroutine demorou muito para responder")
	}
}
