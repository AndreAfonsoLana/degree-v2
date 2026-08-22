package dto

type GetTemperaturaResponse struct {
	Cidade string  `json:"cidade"`
	TempC  float64 `json:"temp_C"`
	TempF  float64 `json:"temp_F"`
	TempK  float64 `json:"temp_K"`
}
