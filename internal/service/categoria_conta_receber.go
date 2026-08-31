package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type CategoriaContaReceberService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *CategoriaContaReceberService) CriarCategoria(ctx context.Context, c *model.CategoriaContaReceber) (*model.CategoriaContaReceber, error) {
	if !c.Validar() {
		return nil, errors.New("o nome da categoria é obrigatório")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	resultado, err := s.repository.CategoriasContasReceber.CriarCategoria(ctx, tx, c)
	if err != nil {
		return nil, err
	}
	return resultado, tx.Commit()
}
func (s *CategoriaContaReceberService) ListarCategorias(ctx context.Context) ([]model.CategoriaContaReceber, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	resultado, err := s.repository.CategoriasContasReceber.ListarCategorias(ctx, tx)
	if err != nil {
		return nil, err
	}
	return resultado, tx.Commit()
}
