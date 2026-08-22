package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndreAfonsoLana/go-degree-validador-cep/infra/service"
	"github.com/AndreAfonsoLana/go-degree-validador-cep/internal/usecase"
	"github.com/AndreAfonsoLana/go-degree-validador-cep/internal/usecase/dto"
)

type mockProvedor struct {
	Result service.TemperaturaResult
}

func (m *mockProvedor) GetTemperaturaByCep(cep string) <-chan service.TemperaturaResult {
	ch := make(chan service.TemperaturaResult, 1)
	ch <- m.Result
	close(ch)
	return ch
}

func TestHandleTemperatura_Sucesso(t *testing.T) {
	mockService := &mockProvedor{
		Result: service.TemperaturaResult{
			City: "Campinas", Temp_C: 25.0, Err: nil,
		},
	}
	uc := usecase.NewGetTemperaturaUsecase(mockService)
	handler := NewTemperaturaHandler(uc)

	req, _ := http.NewRequest(http.MethodGet, "/temperatura?cep=13015904", nil)

	recorder := httptest.NewRecorder()

	handler.HandleTemperatura(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Esperava status 200, recebeu %d", recorder.Code)
	}

	var resposta dto.GetTemperaturaOutputDTO
	json.NewDecoder(recorder.Body).Decode(&resposta)

	if resposta.Cidade != "Campinas" {
		t.Errorf("Esperava cidade Campinas, recebeu %s", resposta.Cidade)
	}
}
