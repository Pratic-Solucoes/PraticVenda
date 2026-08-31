package service

import (
	"testing"

	"gestao/internal/model"
)

func TestValidarItensSaida(t *testing.T) {
	item := func(id uint64) model.ProdutoSaidaEstoque {
		return model.ProdutoSaidaEstoque{IDProduto: id, Quantidade: 1, ValorUnitario: 10, ValorCusto: 8, ValorTotal: 10}
	}

	tests := []struct {
		nome    string
		itens   []model.ProdutoSaidaEstoque
		temErro bool
	}{
		{nome: "produtos distintos", itens: []model.ProdutoSaidaEstoque{item(1), item(2)}},
		{nome: "produto repetido", itens: []model.ProdutoSaidaEstoque{item(1), item(1)}, temErro: true},
		{nome: "quantidade inválida", itens: []model.ProdutoSaidaEstoque{{IDProduto: 1, ValorUnitario: 10, ValorCusto: 8, ValorTotal: 10}}, temErro: true},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			err := validarItensSaida(tt.itens)
			if (err != nil) != tt.temErro {
				t.Fatalf("validarItensSaida() error = %v, temErro = %v", err, tt.temErro)
			}
		})
	}
}
