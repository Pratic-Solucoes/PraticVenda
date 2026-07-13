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

// Produto representa os campos básicos de tb_produtos necessários para o estoque
type Produto struct {
	ID                 int64     `json:"id" db:"id"`
	CodigoBarras       *string   `json:"codigo_barras,omitempty" db:"codigo_barras"`
	CodigoInternoLoja  *string   `json:"codigo_interno_loja,omitempty" db:"codigo_interno_loja"`
	Nome               string    `json:"nome" db:"nome"`
	Descricao          *string   `json:"descricao,omitempty" db:"descricao"`
	PrecoCusto         float64   `json:"preco_custo" db:"preco_custo"`
	PrecoVenda         float64   `json:"preco_venda" db:"preco_venda"`
	UnidadeEstoque     string    `json:"unidade_estoque" db:"unidade_estoque"`
	UnidadeVenda       string    `json:"unidade_venda" db:"unidade_venda"`
	PesoBruto          float64   `json:"peso_bruto" db:"peso_bruto"`
	PesoLiquido        float64   `json:"peso_liquido" db:"peso_liquido"`
	Ativo              bool      `json:"ativo" db:"ativo"`
	CriadoEm           time.Time `json:"criado_em" db:"criado_em"`
	AtualizadoEm       time.Time `json:"atualizado_em" db:"atualizado_em"`
}

// ProdutoEstoque representa a tabela tb_produtos_estoque
type ProdutoEstoque struct {
	ID             int64     `json:"id" db:"id"`
	IDProduto      int64     `json:"id_produto" db:"id_produto"`
	IDEstoque      int64     `json:"id_estoque" db:"id_estoque"`
	Quantidade     float64   `json:"quantidade" db:"quantidade"`
	EstoqueMinimo  float64   `json:"estoque_minimo" db:"estoque_minimo"`
	AtualizadoEm   time.Time `json:"atualizado_em" db:"atualizado_em"`

	// Relacionamentos
	Produto *Produto `json:"produto,omitempty"`
}
