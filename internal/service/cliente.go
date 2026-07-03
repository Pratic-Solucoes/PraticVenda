package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
	"gestao/pkg/helpers"
)

type ClienteService struct {
	db         *sql.DB
	repository *repository.Repository
}

func (s *ClienteService) CriarCliente(ctx context.Context, c *model.Cliente) (*model.Cliente, error) {

	if err := c.Validar(); err != nil {
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

	clienteCriado, err := s.repository.Clientes.CriarCliente(ctx, tx, c)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return clienteCriado, nil
}

func (s *ClienteService) ListarClientes(ctx context.Context, busca string) ([]model.Cliente, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	clientes, err := s.repository.Clientes.ListarClientes(ctx, tx, busca)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return clientes, nil
}

func (s *ClienteService) ObterClientePorID(ctx context.Context, id int64) (*model.Cliente, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	cliente, err := s.repository.Clientes.ObterClientePorID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return cliente, nil
}

func (s *ClienteService) AtualizarCliente(ctx context.Context, id int64, c *model.Cliente) error {
	if err := c.Validar(); err != nil {
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

	err = s.repository.Clientes.AtualizarCliente(ctx, tx, id, c)
	if err != nil {
		return err
	}

	return tx.Commit()
}
