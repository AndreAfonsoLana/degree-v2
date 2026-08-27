package usecase

import (
	"context"
	"errors"

	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/infra/service"
	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/internal/usecase/dto"
	"go.opentelemetry.io/otel"
)

type CEPProvedor interface {
	GetCidadeByCep(ctx context.Context, cep string) <-chan service.CidadeResult
}

// 👇 MUDANÇA 1: Adicionamos o 'ctx context.Context' aqui na interface
type TemperaturaProvedor interface {
	GetTemperaturaByCidade(ctx context.Context, cidade string) <-chan service.TemperaturaResultado
}

type ConsultaClimaUseCase struct {
	cepProvedor         CEPProvedor
	temperaturaProvedor TemperaturaProvedor
}

func NewConsultaClimaUseCase(
	cep CEPProvedor,
	temperatura TemperaturaProvedor,
) *ConsultaClimaUseCase {
	return &ConsultaClimaUseCase{
		cepProvedor:         cep,
		temperaturaProvedor: temperatura,
	}
}

func (c *ConsultaClimaUseCase) ConsultarClimaPorCEP(ctx context.Context, cep string) (dto.GetTemperaturaOutputDTO, error) {
	tr := otel.Tracer("orquestrador-clima-service")
	ctx, span := tr.Start(ctx, "ConsultarClimaPorCEP")
	defer span.End()

	if len(cep) != 8 {
		return dto.GetTemperaturaOutputDTO{}, errors.New("invalid zipcode")
	}

	cidadeResult := <-c.cepProvedor.GetCidadeByCep(ctx, cep)
	if cidadeResult.Err != nil {
		return dto.GetTemperaturaOutputDTO{}, errors.New(cidadeResult.Err.Error())
	}

	temperaturas := <-c.temperaturaProvedor.GetTemperaturaByCidade(ctx, cidadeResult.Cidade)

	return dto.GetTemperaturaOutputDTO{
		Cidade: cidadeResult.Cidade,
		TempC:  temperaturas.Temp_c,
		TempF:  temperaturas.Temp_F,
		TempK:  temperaturas.Temp_K,
	}, nil
}
