package usecase

import (
	"context"
	"errors"

	"github.com/AndreAfonsoLana/go-degree-validador-cep/infra/service"
	"github.com/AndreAfonsoLana/go-degree-validador-cep/internal/usecase/dto"
)

type CEPProvedor interface {
	GetTemperaturaByCep(ctx context.Context, cep string) <-chan service.TemperaturaResult
}

type ConsultaClimaUseCase struct {
	cepProvedor CEPProvedor
}

func NewGetTemperaturaUsecase(cep CEPProvedor) *ConsultaClimaUseCase {
	return &ConsultaClimaUseCase{
		cepProvedor: cep,
	}
}

func (t *ConsultaClimaUseCase) GetTemperaturaUseCase(ctx context.Context, cep string) (dto.GetTemperaturaOutputDTO, error) {
	if len(cep) != 8 {
		return dto.GetTemperaturaOutputDTO{}, errors.New("invalid zipcode")
	}

	temperaturaResult := <-t.cepProvedor.GetTemperaturaByCep(ctx, cep)

	if temperaturaResult.Err != nil {
		return dto.GetTemperaturaOutputDTO{}, errors.New(temperaturaResult.Err.Error())
	}
	return dto.GetTemperaturaOutputDTO{
		Cidade: temperaturaResult.City,
		TempC:  float64(temperaturaResult.Temp_C),
		TempF:  float64(temperaturaResult.Temp_F),
		TempK:  float64(temperaturaResult.Temp_K),
	}, nil
}
