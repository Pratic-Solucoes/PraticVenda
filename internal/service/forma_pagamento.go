package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
	"gestao/pkg/helpers"
)

type FormaPagamentoService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *FormaPagamentoService) Criar(ctx context.Context, fp *model.FormaPagamento) (*model.FormaPagamento, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	fpCriado, err := s.repository.FormasPagamento.Criar(ctx, tx, fp)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return fpCriado, nil
}
