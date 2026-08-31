package model

import (
	"errors"
	"strings"
)

type ItemVendaPDV struct {
	IDProduto     int64    `json:"id_produto"`
	Quantidade    float64  `json:"quantidade"`
	ValorUnitario *float64 `json:"valor_unitario,omitempty"`
}
type VendaPDV struct {
	ID                  int64          `json:"id,omitempty"`
	Itens               []ItemVendaPDV `json:"itens"`
	IDCliente           *int64         `json:"id_cliente,omitempty"`
	IDFormaPagamento    int64          `json:"id_forma_pagamento"`
	IDCondicaoPagamento int64          `json:"id_condicao_pagamento"`
	ValorDesconto       float64        `json:"valor_desconto"`
	ApelidoConsumidor   string         `json:"apelido_consumidor,omitempty"`
}

type PreVendaPDV struct {
	ID                int64          `json:"id"`
	IDCliente         *int64         `json:"id_cliente,omitempty"`
	Cliente           string         `json:"cliente,omitempty"`
	ApelidoConsumidor string         `json:"apelido_consumidor,omitempty"`
	ValorTotal        float64        `json:"valor_total"`
	Itens             []ItemVendaPDV `json:"itens,omitempty"`
}

func (v VendaPDV) Validar() error {
	if len(v.Itens) == 0 {
		return errors.New("a venda não possui itens")
	}
	if v.IDFormaPagamento <= 0 || v.IDCondicaoPagamento <= 0 {
		return errors.New("a forma e a condição de pagamento são obrigatórias")
	}
	if v.ValorDesconto < 0 {
		return errors.New("o desconto não pode ser negativo")
	}
	for _, i := range v.Itens {
		if i.IDProduto <= 0 || i.Quantidade <= 0 {
			return errors.New("item de venda inválido")
		}
	}
	return nil
}

func (v *VendaPDV) ValidarItens() error {
	v.ApelidoConsumidor = strings.TrimSpace(v.ApelidoConsumidor)
	if len(v.Itens) == 0 {
		return errors.New("a venda não possui itens")
	}
	if v.ValorDesconto < 0 {
		return errors.New("o desconto não pode ser negativo")
	}
	for _, i := range v.Itens {
		if i.IDProduto <= 0 || i.Quantidade <= 0 {
			return errors.New("item de venda inválido")
		}
	}
	return nil
}
