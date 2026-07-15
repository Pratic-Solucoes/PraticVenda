package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
	"gestao/pkg/helpers"
)

type EntradaEstoqueService struct {
	db          *sql.DB
	repositorio *repository.Repository
}

func (s *EntradaEstoqueService) RegistrarEntrada(ctx context.Context, entrada *model.EntradaEstoque) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verifica se algum produto da entrada de estoque está com dado incorreto
	for _, produto := range entrada.Produtos {
		if err := produto.Validar(); err != nil {
			return err
		}
	}
	// Seta o schema do usuario
	if err := helpers.SetSchema(ctx, tx); err != nil {
		return err
	}

	// Lógica para criar a entrada de estoque no BD
	if err := s.repositorio.EntradaEstoque.RegistrarEntrada(ctx, tx, entrada); err != nil {
		return err
	}

	if entrada.Status == "CONCLUIDA" {
		// Lógica para inserir os produtos, movimento de estoque e atualizar o saldo do produto no estoque
		// Só haverá movimentação de estoque caso a entrada seja concluída
		if err := s.repositorio.EntradaEstoque.RegistrarProdutosEntrada(ctx, tx, entrada); err != nil {
			return err
		}
	}

	return tx.Commit()
}
