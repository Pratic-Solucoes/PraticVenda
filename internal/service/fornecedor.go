package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
	"gestao/pkg/helpers"
)

type FornecedorService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *FornecedorService) ListarFornecedores(ctx context.Context, busca string) ([]*model.Fornecedor, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	fornecedores, err := s.repository.Fornecedores.ListarFornecedores(ctx, tx, busca)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return fornecedores, nil
}

// CriarFornecedor gerencia a regra de negócio para criar um fornecedor.
// Valida dados (se necessário) e abre a transação para o repositório.
func (s *FornecedorService) CriarFornecedor(ctx context.Context, f *model.Fornecedor) (*model.Fornecedor, error) {
	if err := f.Validar(); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	fornecedorCriado, err := s.repository.Fornecedores.CriarFornecedor(ctx, tx, f)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return fornecedorCriado, nil
}

func (s *FornecedorService) ObterFornecedorPorID(ctx context.Context, id int64) (*model.Fornecedor, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	fornecedor, err := s.repository.Fornecedores.ObterFornecedorPorID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return fornecedor, nil
}

func (s *FornecedorService) AtualizarFornecedor(ctx context.Context, id int64, f *model.Fornecedor) error {
	if err := f.Validar(); err != nil {
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

	err = s.repository.Fornecedores.AtualizarFornecedor(ctx, tx, id, f)
	if err != nil {
		return err
	}

	return tx.Commit()
}
