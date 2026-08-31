package model

import (
	"errors"
	"strings"
)

type ContaReceber struct {
	ID                  int64                  `json:"id" db:"id"`
	IDCliente           int64                  `json:"id_cliente" db:"id_cliente"`
	IDCategoria         int64                  `json:"id_categoria" db:"id_categoria"`
	IDCondicaoPagamento int64                  `json:"id_condicao_pagamento" db:"id_condicao_pagamento"`
	IDFormaPagamento    int64                  `json:"id_forma_pagamento" db:"id_forma_pagamento"`
	IDGrupoParcelas     string                 `json:"id_grupo_parcelas" db:"id_grupo_parcelas"`
	Descricao           string                 `json:"descricao" db:"descricao"`
	ValorOriginal       float64                `json:"valor_original" db:"valor_original"`
	SaldoRestante       float64                `json:"saldo_restante" db:"saldo_restante"`
	DtVencimento        string                 `json:"dt_vencimento" db:"dt_vencimento"`
	NrParcela           int                    `json:"nr_parcela" db:"nr_parcela"`
	NrTotalParcelas     int                    `json:"nr_total_parcelas" db:"nr_total_parcelas"`
	Status              string                 `json:"status" db:"status"`
	DtEmissao           string                 `json:"dt_emissao" db:"dt_emissao"`
	DtPagamento         *string                `json:"dt_pagamento,omitempty" db:"dt_pagamento"`
	Cliente             *Cliente               `json:"cliente,omitempty"`
	Categoria           *CategoriaContaReceber `json:"categoria,omitempty"`
	CondicaoPagamento   *CondicaoPagamento     `json:"condicao_pagamento,omitempty"`
	FormaPagamento      *FormaPagamento        `json:"forma_pagamento,omitempty"`
}

// ContaReceberCriar representa o lançamento inteiro. A condição define as parcelas;
// o valor informado é distribuído entre elas com ajuste de centavos na última parcela.
type ContaReceberCriar struct {
	IDCliente           int64   `json:"id_cliente"`
	IDCategoria         int64   `json:"id_categoria"`
	IDCondicaoPagamento int64   `json:"id_condicao_pagamento"`
	IDFormaPagamento    int64   `json:"id_forma_pagamento"`
	Descricao           string  `json:"descricao"`
	ValorTotal          float64 `json:"valor_total"`
	DtEmissao           string  `json:"dt_emissao"`
}

func (d *ContaReceberCriar) Validar() error {
	d.Descricao = strings.TrimSpace(d.Descricao)
	var erros []error
	if d.IDCliente <= 0 {
		erros = append(erros, errors.New("o cliente é obrigatório"))
	}
	if d.IDCategoria <= 0 {
		erros = append(erros, errors.New("a categoria de crédito é obrigatória"))
	}
	if d.IDCondicaoPagamento <= 0 {
		erros = append(erros, errors.New("a condição de pagamento é obrigatória"))
	}
	if d.IDFormaPagamento <= 0 {
		erros = append(erros, errors.New("a forma de pagamento é obrigatória"))
	}
	if d.Descricao == "" {
		erros = append(erros, errors.New("a descrição é obrigatória"))
	}
	if d.ValorTotal <= 0 {
		erros = append(erros, errors.New("o valor total deve ser maior que zero"))
	}
	if d.DtEmissao == "" {
		erros = append(erros, errors.New("a data de emissão é obrigatória"))
	}
	return errors.Join(erros...)
}
