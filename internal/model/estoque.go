package model

import (
	"errors"
	"time"
)

// Estoque representa a tabela tb_estoques
type Estoque struct {
	ID        int64     `json:"id" db:"id"`
	Nome      string    `json:"nome" db:"nome"`
	Descricao *string   `json:"descricao,omitempty" db:"descricao"`
	Ativo     bool      `json:"ativo" db:"ativo"`
	CriadoEm  time.Time `json:"criado_em" db:"criado_em"`
}

// EstoqueCriar é o DTO para criação de estoque
type EstoqueCriar struct {
	Nome      string  `json:"nome"`
	Descricao *string `json:"descricao,omitempty"`
}

// Validar valida se os campos do DTO estão corretos
func (e *EstoqueCriar) Validar() error {
	if e.Nome == "" {
		return errors.New("o nome do estoque é obrigatório")
	}
	return nil
}

// ProdutoEstoque representa a tabela tb_produtos_estoque
type ProdutoEstoque struct {
	ID            int64     `json:"id" db:"id"`
	IDProduto     int64     `json:"id_produto" db:"id_produto"`
	IDEstoque     int64     `json:"id_estoque" db:"id_estoque"`
	Quantidade    float64   `json:"quantidade" db:"quantidade"`
	EstoqueMinimo float64   `json:"estoque_minimo" db:"estoque_minimo"`
	AtualizadoEm  time.Time `json:"atualizado_em" db:"atualizado_em"`

	// Relacionamentos
	Produto *Produto `json:"produto,omitempty"`
}

type ProdutoEntradaEstoque struct {
	ID                     uint64  `json:"id,omitempty" db:"id"`
	IDProduto              uint64  `json:"id_produto,omitempty" db:"id_produto"`
	ValorUnitario          float64 `json:"valor_unitario,omitempty" db:"valor_unitario"`
	ValorIcmsST            float64 `json:"valor_icms_st,omitempty" db:"valor_icms_st"`
	ValorIPI               float64 `json:"valor_ipi,omitempty" db:"valor_ipi"`
	ValorDesconto          float64 `json:"valor_desconto,omitempty"`
	RateioDespesaAdicional float64 `json:"rateio_despesa_adicional"`
	ValorCusto             float64 `json:"valor_custo,omitempty"`
	ValorTotal             float64 `json:"valor_total,omitempty"`
	Quantidade             float64 `json:"quantidade,omitempty" db:"quantidade"`
}

type EntradaEstoque struct {
	ID               uint64                  `json:"id,omitempty" db:"id"`
	IDFornecedor     uint64                  `json:"id_fornecedor,omitempty" db:"id_fornecedor"`
	Produtos         []ProdutoEntradaEstoque `json:"produtos,omitempty"`
	DespesaAdicional float64                 `json:"despesa_adicional,omitempty" db:"despesa_adicional"`
	UsuarioID        int64                   `json:"usuario_id,omitempty" db:"usuario_id"`
	ValorTotal       float64                 `json:"valor_total,omitempty" db:"valor_total"`
	Status           string                  `json:"status,omitempty" db:"status"`
	CriadoEm         time.Time               `json:"criado_em,omitempty" db:"criado_em"`
}
