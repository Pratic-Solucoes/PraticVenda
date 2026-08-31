package model

import (
	"errors"
	"time"
)

type ProdutoEntradaEstoque struct {
	ID                     uint64  `json:"id,omitempty" db:"id"`
	IDProduto              uint64  `json:"id_produto,omitempty" db:"id_produto"`
	NomeProduto            string  `json:"nome_produto,omitempty" db:"nome_produto"`
	IDEntradaEstoque       uint64  `json:"id_entrada_estoque,omitempty"`
	ValorUnitario          float64 `json:"valor_unitario,omitempty" db:"valor_unitario"`
	ValorIcmsST            float64 `json:"valor_icms_st,omitempty" db:"valor_icms_st"`
	ValorIPI               float64 `json:"valor_ipi,omitempty" db:"valor_ipi"`
	ValorDesconto          float64 `json:"valor_desconto,omitempty"`
	RateioDespesaAdicional float64 `json:"rateio_despesa_adicional"`
	ValorCusto             float64 `json:"valor_custo,omitempty"`
	ValorTotal             float64 `json:"valor_total,omitempty"`
	Quantidade             float64 `json:"quantidade,omitempty" db:"quantidade"`
}

func (p *ProdutoEntradaEstoque) Validar() error {
	if p.IDProduto == 0 {
		return errors.New("ID do produto é obrigatório")
	}
	if p.ValorUnitario <= 0 {
		return errors.New("Valor unitário deve ser maior que zero")
	}
	if p.ValorCusto <= 0 {
		return errors.New("Valor de custo deve ser maior que zero")
	}
	if p.ValorTotal <= 0 {
		return errors.New("Valor total deve ser maior que zero")
	}
	if p.Quantidade <= 0 {
		return errors.New("Quantidade deve ser maior que zero")
	}
	return nil
}

type EntradaEstoque struct {
	ID                    uint64                  `json:"id,omitempty" db:"id"`
	IDEstoque             uint64                  `json:"id_estoque,omitempty"`
	Estoque               string                  `json:"estoque,omitempty"`
	IDFornecedor          uint64                  `json:"id_fornecedor,omitempty" db:"id_fornecedor"`
	Fornecedor            string                  `json:"fornecedor,omitempty"`
	Produtos              []ProdutoEntradaEstoque `json:"produtos,omitempty"`
	ValorDespesaAdicional float64                 `json:"despesa_adicional,omitempty" db:"despesa_adicional"`
	IDUsuario             int64                   `json:"usuario_id,omitempty" db:"usuario_id"`
	Usuario               string                  `json:"usuario,omitempty" db:"usuario"`
	ValorTotal            float64                 `json:"valor_total,omitempty" db:"valor_total"`
	Status                string                  `json:"status,omitempty" db:"status"`
	CriadoEm              time.Time               `json:"criado_em,omitempty" db:"criado_em"`
}

// FiltroEntradaEstoque reúne os filtros disponíveis na listagem de entradas.
type FiltroEntradaEstoque struct {
	ID         uint64
	Fornecedor string
	Data       string
	Status     string
}
