package model

// MovimentoFinanceiro representa uma entrada ou saída no fluxo de caixa
// CRIADO NA REFATORAÇÃO: Registrar histórico de pagamentos (parciais ou totais)
type MovimentoFinanceiro struct {
	ID             int64   `json:"id" db:"id"`
	TipoMovimento  string  `json:"tipo_movimento" db:"tipo_movimento"` // ENTRADA ou SAIDA
	IDOrigem       *int64  `json:"id_origem,omitempty" db:"id_origem"`
	DtMovimento    string  `json:"dt_movimento" db:"dt_movimento"`
	ValorMovimento float64 `json:"valor_movimento" db:"valor_movimento"`
	ValorAcrescimo float64 `json:"valor_acrescimo" db:"valor_acrescimo"`
	ValorDesconto  float64 `json:"valor_desconto" db:"valor_desconto"`
	FormaPagamento *string `json:"forma_pagamento,omitempty" db:"forma_pagamento"`
	CreatedAt      string  `json:"created_at" db:"created_at"`
}
