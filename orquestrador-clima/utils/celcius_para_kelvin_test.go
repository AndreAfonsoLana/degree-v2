package utils

import "testing"

func TestCelsiusParaKelvin(t *testing.T) {
	celsius := 29.0
	esperadoKelvin := 302.15
	result := CelsiusParaKelvin(celsius)

	if result != esperadoKelvin {
		t.Errorf("Esperado %f, mas recebeu %f", esperadoKelvin, result)
	}
}
