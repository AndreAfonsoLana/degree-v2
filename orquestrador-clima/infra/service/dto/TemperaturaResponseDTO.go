package dto

type TemperaturaResponseDTO struct {
	Current struct {
		Temp_c float64 `json:"temp_c"`
	} `json:"current"`
}
