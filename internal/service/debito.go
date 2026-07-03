package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
	"gestao/pkg/helpers"
)

type DebitoService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *DebitoService) LancarDebito(ctx context.Context, debito *model.DebitoAvulsoCriar) error {
	if err := debito.Validar(); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return err
	}

	err = s.repository.Debitos.LancarDebito(ctx, tx, debito)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *DebitoService) ListarDebitos(ctx context.Context, busca, vencimento, status string) ([]*model.Debito, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	debitos, err := s.repository.Debitos.ListarDebitos(ctx, tx, busca, vencimento, status)
	if err != nil {
		return nil, err
	}
	
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return debitos, nil
}

func (s *DebitoService) PagarDebito(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return err
	}

	err = s.repository.Debitos.PagarDebito(ctx, tx, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *DebitoService) EditarDebito(ctx context.Context, id int64, debito *model.DebitoAvulsoCriar) error {
	if err := debito.Validar(); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return err
	}

	err = s.repository.Debitos.EditarDebito(ctx, tx, id, debito)
	if err != nil {
		return err
	}
	return tx.Commit()
}
