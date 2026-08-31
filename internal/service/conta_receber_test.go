package service

import (
	"reflect"
	"testing"
)

func TestDividirEmParcelasPreservaValorTotal(t *testing.T) {
	parcelas := dividirEmParcelas(100.00, 3)
	esperado := []float64{33.34, 33.33, 33.33}
	if !reflect.DeepEqual(parcelas, esperado) {
		t.Fatalf("parcelas = %v, esperado %v", parcelas, esperado)
	}

	var total float64
	for _, valor := range parcelas {
		total += valor
	}
	if total != 100.00 {
		t.Fatalf("total das parcelas = %.2f, esperado 100.00", total)
	}
}
