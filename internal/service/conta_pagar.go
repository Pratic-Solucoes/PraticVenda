package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"gestao/internal/repository"
	"gestao/pkg/helpers"
)

type ContaPagarService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *ContaPagarService) CriarContaPagar(ctx context.Context, contaPagar *model.ContaPagarCriar) error {
	if err := contaPagar.Validar(); err != nil {
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

	err = s.repository.ContasPagar.CriarContaPagar(ctx, tx, contaPagar)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *ContaPagarService) ListarContasPagar(ctx context.Context, busca, vencimento, status string) ([]*model.ContaPagar, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	contasPagar, err := s.repository.ContasPagar.ListarContasPagar(ctx, tx, busca, vencimento, status)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return contasPagar, nil
}

func (s *ContaPagarService) PagarContaPagar(ctx context.Context, id int64, valorPagamento float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return err
	}

	conta, err := s.repository.ContasPagar.BuscarPorID(ctx, tx, id)
	if err != nil {
		return err
	}

	if conta.Status == "PAGO" {
		return repository.CONTA_PAGAR_QUITADA
	}

	if valorPagamento <= 0 {
		return errors.New("valor do pagamento deve ser maior que zero")
	}

	if valorPagamento > conta.SaldoRestante {
		return errors.New("valor do pagamento não pode ser maior que o saldo restante")
	}

	err = s.repository.ContasPagar.PagarContaPagar(ctx, tx, id, valorPagamento)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *ContaPagarService) EditarContaPagar(ctx context.Context, id int64, contaPagar *model.ContaPagarCriar) error {
	if err := contaPagar.Validar(); err != nil {
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

	err = s.repository.ContasPagar.EditarContaPagar(ctx, tx, id, contaPagar)
	if err != nil {
		return err
	}
	return tx.Commit()
}
