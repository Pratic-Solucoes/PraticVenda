package model

type CondicaoPagamento struct {
	ID                uint64 `json:"id"`
	IDFormaPagamento  uint64 `json:"id_forma_pagamento"`
	Descricao         string `json:"descricao"`
	QtdParcelas       int64  `json:"qtd_parcelas"`
	DiasPrimeiroVenc  int64  `json:"dias_primeiro_venc"`
	IntervaloParcelas int64  `json:"intervalo_parcelas"`
}
