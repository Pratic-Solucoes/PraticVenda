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
