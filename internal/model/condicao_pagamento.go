package model

import (
	"errors"
	"strings"
)

type CondicaoPagamento struct {
	ID                uint64   `json:"id"`
	FormasPagamento   []uint64 `json:"formas_pagamento"`
	Descricao         string   `json:"descricao"`
	QtdParcelas       int64    `json:"qtd_parcelas"`
	DiasPrimeiroVenc  int64    `json:"dias_primeiro_venc"`
	IntervaloParcelas int64    `json:"intervalo_parcelas"`
}

func (cp *CondicaoPagamento) Validar() error {
	var erros []error
	cp.Descricao = strings.TrimSpace(cp.Descricao)
	if cp.Descricao == "" {
		erros = append(erros, errors.New("necessário informar a descrição da condição de pagamento"))
	}
	if cp.QtdParcelas <= 0 {
		erros = append(erros, errors.New("necessário informar uma quantidade de parcelas maior que zero"))
	}
	if cp.DiasPrimeiroVenc < 0 {
		erros = append(erros, errors.New("os dias para o primeiro vencimento não podem ser negativos"))
	}
	if cp.IntervaloParcelas < 0 {
		erros = append(erros, errors.New("o intervalo entre parcelas não pode ser negativo"))
	}
	if len(cp.FormasPagamento) == 0 {
		erros = append(erros, errors.New("selecione ao menos uma forma de pagamento"))
	}
	formasSelecionadas := make(map[uint64]struct{}, len(cp.FormasPagamento))
	for _, idFormaPagamento := range cp.FormasPagamento {
		if idFormaPagamento == 0 {
			erros = append(erros, errors.New("forma de pagamento inválida"))
			continue
		}
		if _, existe := formasSelecionadas[idFormaPagamento]; existe {
			erros = append(erros, errors.New("uma forma de pagamento foi informada mais de uma vez"))
			continue
		}
		formasSelecionadas[idFormaPagamento] = struct{}{}
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}
	return nil
}
