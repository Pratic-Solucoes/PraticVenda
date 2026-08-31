package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
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

func (s *FormaPagamentoService) Listar(ctx context.Context) ([]model.FormaPagamento, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	fps, err := s.repository.FormasPagamento.Listar(ctx, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return fps, nil
}

func (s *FormaPagamentoService) BuscarPorID(ctx context.Context, idFp int64) (*model.FormaPagamento, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	fp, err := s.repository.FormasPagamento.BuscarPorID(ctx, tx, idFp)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return fp, nil
}

func (s *FormaPagamentoService) Atualizar(ctx context.Context, fp *model.FormaPagamento) (*model.FormaPagamento, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	fpAtualizado, err := s.repository.FormasPagamento.Atualizar(ctx, tx, fp)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return fpAtualizado, nil
}
