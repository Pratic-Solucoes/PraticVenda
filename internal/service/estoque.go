package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type EstoqueService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *EstoqueService) CriarEstoque(ctx context.Context, input *model.EstoqueCriar) (*model.Estoque, error) {
	if err := input.Validar(); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	e := &model.Estoque{
		Nome:      input.Nome,
		Descricao: input.Descricao,
		Ativo:     true,
	}

	eCriado, err := s.repository.Estoques.CriarEstoque(ctx, tx, e)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return eCriado, nil
}

func (s *EstoqueService) ListarEstoques(ctx context.Context) ([]*model.Estoque, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	estoques, err := s.repository.Estoques.ListarEstoques(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return estoques, nil
}

func (s *EstoqueService) ListarProdutosDoEstoque(ctx context.Context, idEstoque int64) ([]*model.ProdutoEstoque, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	produtos, err := s.repository.Estoques.ListarProdutosDoEstoque(ctx, tx, idEstoque)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return produtos, nil
}
