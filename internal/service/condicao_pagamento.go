package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type CondicaoPagamentoService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *CondicaoPagamentoService) Criar(ctx context.Context, cp *model.CondicaoPagamento) error {
	if err := cp.Validar(); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repository.CondicoesPagamento.Criar(ctx, tx, cp); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *CondicaoPagamentoService) Listar(ctx context.Context) ([]model.CondicaoPagamento, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	condicoes, err := s.repository.CondicoesPagamento.Listar(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return condicoes, nil
}

func (s *CondicaoPagamentoService) BuscarPorID(ctx context.Context, id int64) (*model.CondicaoPagamento, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	condicao, err := s.repository.CondicoesPagamento.BuscarPorID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return condicao, nil
}

func (s *CondicaoPagamentoService) Atualizar(ctx context.Context, cp *model.CondicaoPagamento) error {
	if err := cp.Validar(); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repository.CondicoesPagamento.Atualizar(ctx, tx, cp); err != nil {
		return err
	}

	return tx.Commit()
}
