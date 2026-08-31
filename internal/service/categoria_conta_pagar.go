package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type CategoriaContaPagarService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *CategoriaContaPagarService) CriarCategoria(ctx context.Context, c *model.CategoriaContaPagar) (*model.CategoriaContaPagar, error) {
	if c.Nome == "" {
		return nil, errors.New("o nome da categoria é obrigatório")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	categoriaCriada, err := s.repository.CategoriasContasPagar.CriarCategoria(ctx, tx, c)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return categoriaCriada, nil
}

func (s *CategoriaContaPagarService) ListarCategorias(ctx context.Context) ([]model.CategoriaContaPagar, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	categorias, err := s.repository.CategoriasContasPagar.ListarCategorias(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return categorias, nil
}
