package usecase

import (
	"errors"
	"testing"

	"github.com/AndreAfonsoLana/go-degree-validador-cep/infra/service"
)

type mockCEPProvedor struct {
	Result service.TemperaturaResult
}

func (m *mockCEPProvedor) GetTemperaturaByCep(cep string) <-chan service.TemperaturaResult {
	ch := make(chan service.TemperaturaResult, 1)
	ch <- m.Result
	close(ch)
	return ch
}

func TestGetTemperaturaUseCase(t *testing.T) {

	cenarios := []struct {
		nome           string
		cepInput       string
		mockResult     service.TemperaturaResult
		esperaErro     bool
		erroEsperado   string
		cidadeEsperada string
	}{
		{
			nome:     "Sucesso - Deve retornar os dados corretamente",
			cepInput: "13183310", // CEP válido com 8 dígitos
			mockResult: service.TemperaturaResult{
				City:   "Hortolândia",
				Temp_C: 28.5,
				Temp_F: 83.3,
				Temp_K: 301.65,
				Err:    nil,
			},
			esperaErro:     false,
			cidadeEsperada: "Hortolândia",
		},
		{
			nome:           "Erro - CEP com tamanho inválido",
			cepInput:       "123",                       // Inválido
			mockResult:     service.TemperaturaResult{}, // Mock nem será chamado
			esperaErro:     true,
			erroEsperado:   "invalid zipcode",
			cidadeEsperada: "",
		},
		{
			nome:     "Erro - Erro retornado pelo serviço (API fora do ar)",
			cepInput: "13183310",
			mockResult: service.TemperaturaResult{
				Err: errors.New("timeout na api"),
			},
			esperaErro:     true,
			erroEsperado:   "timeout na api",
			cidadeEsperada: "",
		},
	}

	for _, cenario := range cenarios {
		t.Run(cenario.nome, func(t *testing.T) {
			mockProvedor := &mockCEPProvedor{Result: cenario.mockResult}
			useCase := NewGetTemperaturaUsecase(mockProvedor)

			resultado, erro := useCase.GetTemperaturaUseCase(cenario.cepInput)

			if cenario.esperaErro {
				if erro == nil {
					t.Fatalf("Esperava um erro, mas não recebeu nenhum")
				}
				if erro.Error() != cenario.erroEsperado {
					t.Errorf("Esperava o erro '%s', recebeu '%s'", cenario.erroEsperado, erro.Error())
				}
			} else {
				if erro != nil {
					t.Fatalf("Não esperava erro, mas recebeu: %v", erro)
				}
				if resultado.Cidade != cenario.cidadeEsperada {
					t.Errorf("Esperava a cidade '%s', recebeu '%s'", cenario.cidadeEsperada, resultado.Cidade)
				}
				t.Logf("Teste Sucesso: %+v", resultado)
			}
		})
	}
}
