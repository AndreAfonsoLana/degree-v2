package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AndreAfonsoLana/go-degree-orquestrador-clima/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TemperaturaHandler struct {
	UseCase *usecase.ConsultaClimaUseCase
}

func NewTemperaturaHandler(uc *usecase.ConsultaClimaUseCase) *TemperaturaHandler {
	return &TemperaturaHandler{
		UseCase: uc,
	}
}

func (h *TemperaturaHandler) HandleTemperatura(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("orquestrador-clima")
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	ctx, span := tracer.Start(ctx, "HandleTemperatura")

	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("Recebi request na rota: %s\n", r.URL.Path)

	cep := r.URL.Query().Get("cep")

	output, err := h.UseCase.ConsultarClimaPorCEP(ctx, cep)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
