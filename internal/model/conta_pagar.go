package model

import (
	"errors"
)

// ContaPagar representa a entidade principal de conta a pagar no banco de dados
// REATORAÇÃO: Refletindo as mudanças para dar suporte a agrupamento e saldo_restante
type ContaPagar struct {
	ID              int64   `json:"id" db:"id"`
	IDFornecedor    int64   `json:"id_fornecedor" db:"id_fornecedor"`
	IDCategoria     int64   `json:"id_categoria" db:"id_categoria"`                     // NOT NULL no banco
	IDGrupoParcelas *string `json:"id_grupo_parcelas,omitempty" db:"id_grupo_parcelas"` // UUID do grupo
	Descricao       *string `json:"descricao,omitempty" db:"descricao"`
	NrDocumento     *string `json:"nr_documento,omitempty" db:"nr_documento"`
	NrNotaFiscal    *string `json:"nr_nota_fiscal,omitempty" db:"nr_nota_fiscal"`
	ValorOriginal   float64 `json:"valor_original" db:"valor_original"` // Antigo campo "valor"
	SaldoRestante   float64 `json:"saldo_restante" db:"saldo_restante"` // Saldo após pagamentos parciais
	DtEntrada       string  `json:"dt_entrada" db:"dt_entrada"`
	DtVencimento    string  `json:"dt_vencimento" db:"dt_vencimento"`
	NrParcela       int     `json:"nr_parcela" db:"nr_parcela"`
	NrTotalParcelas int     `json:"nr_total_parcelas" db:"nr_total_parcelas"`
	Status          string  `json:"status" db:"status"` // Ex: PENDENTE, PAGO_PARCIAL, PAGO, CANCELADO, ATRASADO
	CreatedAt       string  `json:"created_at" db:"created_at"`
	DtPagamento     *string `json:"dt_pagamento,omitempty" db:"dt_pagamento"`
	UpdatedAt       string  `json:"updated_at" db:"updated_at"`

	// Relacionamentos opcionais
	Fornecedor *Fornecedor          `json:"fornecedor,omitempty"`
	Categoria  *CategoriaContaPagar `json:"categoria,omitempty"`
}

// ContaPagarCriar é um DTO (Data Transfer Object) usado apenas para receber dados de criação
type ContaPagarCriar struct {
	IDFornecedor    int64   `json:"id_fornecedor"`
	IDCategoria     int64   `json:"id_categoria"`
	IDGrupoParcelas *string `json:"id_grupo_parcelas,omitempty"`
	Descricao       string  `json:"descricao"`
	NrDocumento     string  `json:"nr_documento,omitempty"`
	NrNotaFiscal    string  `json:"nr_nota_fiscal,omitempty"`
	ValorOriginal   float64 `json:"valor_original"` // Refatorado
	DtEntrada       string  `json:"dt_entrada"`     // Pode receber string "YYYY-MM-DD" no JSON
	DtVencimento    string  `json:"dt_vencimento"`  // Pode receber string "YYYY-MM-DD" no JSON
	NrParcela       int     `json:"nr_parcela"`
	NrTotalParcelas int     `json:"nr_total_parcelas"`
}

// Validar verifica se os dados obrigatórios para criação foram preenchidos corretamente
func (d *ContaPagarCriar) Validar() error {
	var erros []error
	if d.IDFornecedor == 0 {
		erros = append(erros, errors.New("o id do fornecedor não foi informado"))
	}
	if d.IDCategoria == 0 {
		erros = append(erros, errors.New("o id da categoria não foi informado"))
	}
	if d.Descricao == "" {
		erros = append(erros, errors.New("a descrição não foi informada"))
	}
	if d.ValorOriginal <= 0 {
		erros = append(erros, errors.New("o valor original deve ser maior que zero"))
	}
	if d.DtEntrada == "" {
		erros = append(erros, errors.New("a data de entrada não foi informada"))
	}
	if d.DtVencimento == "" {
		erros = append(erros, errors.New("a data de vencimento não foi informada"))
	}
	if d.NrParcela == 0 {
		erros = append(erros, errors.New("o número da parcela não foi informado"))
	}
	if d.NrTotalParcelas == 0 {
		erros = append(erros, errors.New("o número total de parcelas não foi informado"))
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}
	return nil
}
