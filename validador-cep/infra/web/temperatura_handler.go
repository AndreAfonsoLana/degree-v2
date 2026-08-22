package web

import (
	"encoding/json"
	"net/http"

	"github.com/AndreAfonsoLana/go-degree-validador-cep/internal/usecase"
	"github.com/AndreAfonsoLana/go-degree-validador-cep/internal/usecase/dto"
	"go.opentelemetry.io/otel"
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
	tracer := otel.Tracer("validador-cep")
	ctx, span := tracer.Start(r.Context(), "validador-cep-handler")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input dto.GetTemperaturaInputDTIO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) // 422 para JSON malformado
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid JSON"})
		return
	}

	output, err := h.UseCase.GetTemperaturaUseCase(ctx, input.CEP)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
