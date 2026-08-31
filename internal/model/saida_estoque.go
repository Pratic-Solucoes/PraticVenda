package model

import (
	"errors"
	"time"
)

type ProdutoSaidaEstoque struct {
	ID             uint64  `json:"id,omitempty" db:"id"`
	IDProduto      uint64  `json:"id_produto,omitempty" db:"id_produto"`
	NomeProduto    string  `json:"nome_produto,omitempty" db:"nome_produto"`
	IDSaidaEstoque uint64  `json:"id_saida_estoque,omitempty"`
	ValorUnitario  float64 `json:"valor_unitario,omitempty" db:"valor_unitario"`
	ValorCusto     float64 `json:"valor_custo,omitempty"`
	ValorTotal     float64 `json:"valor_total,omitempty"`
	Quantidade     float64 `json:"quantidade,omitempty" db:"quantidade"`
}

func (p *ProdutoSaidaEstoque) Validar() error {
	if p.IDProduto == 0 {
		return errors.New("ID do produto é obrigatório")
	}
	if p.ValorUnitario <= 0 {
		return errors.New("Valor unitário deve ser maior que zero")
	}
	if p.ValorCusto < 0 {
		return errors.New("Valor de custo não pode ser negativo")
	}
	if p.ValorTotal <= 0 {
		return errors.New("Valor total deve ser maior que zero")
	}
	if p.Quantidade <= 0 {
		return errors.New("Quantidade deve ser maior que zero")
	}
	return nil
}

type SaidaEstoque struct {
	ID         uint64                `json:"id,omitempty" db:"id"`
	IDEstoque  uint64                `json:"id_estoque,omitempty"`
	Estoque    string                `json:"estoque,omitempty"`
	Produtos   []ProdutoSaidaEstoque `json:"produtos,omitempty"`
	IDUsuario  int64                 `json:"usuario_id,omitempty" db:"usuario_id"`
	Usuario    string                `json:"usuario,omitempty" db:"usuario"`
	ValorTotal float64               `json:"valor_total,omitempty" db:"valor_total"`
	Status     string                `json:"status,omitempty" db:"status"`
	CriadoEm   time.Time             `json:"criado_em,omitempty" db:"criado_em"`
}

// FiltroSaidaEstoque reúne os filtros disponíveis na listagem de saídas.
type FiltroSaidaEstoque struct {
	ID     uint64
	Data   string
	Status string
}
