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

	// Lógica realizar entrada de estoque

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
