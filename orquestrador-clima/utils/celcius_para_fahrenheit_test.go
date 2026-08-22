package utils

import (
	"testing"
)

func TestCelciusParaFahrenheit(t *testing.T) {
	celsius := 30.0
	esperadoFahrenheit := 86.0
	resulto := CelsiusParaFahrenheit(celsius)

	if resulto != esperadoFahrenheit {
		t.Errorf("Esperado %f, mas recebeu %f", esperadoFahrenheit, resulto)
	}
}
