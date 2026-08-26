package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
	"gestao/utils/helpers"
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

func (s *ClienteService) CriarEndereco(ctx context.Context, idCliente int64, e *model.EnderecoCliente) (*model.EnderecoCliente, error) {
	if err := e.Validar(); err != nil {
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

	enderecoCriado, err := s.repository.Clientes.CriarEndereco(ctx, tx, idCliente, e)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return enderecoCriado, nil
}

func (s *ClienteService) EditarEndereco(ctx context.Context, idCliente int64, idEndereco int64, e *model.EnderecoCliente) error {
	if err := e.Validar(); err != nil {
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

	err = s.repository.Clientes.EditarEndereco(ctx, tx, idCliente, idEndereco, e)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ClienteService) BuscarEnderecoByID(ctx context.Context, idCliente int64, idEndereco int64) (*model.EnderecoCliente, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	endereco, err := s.repository.Clientes.BuscarEnderecoByID(ctx, tx, idCliente, idEndereco)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return endereco, nil
}

func (s *ClienteService) CriarTelefone(ctx context.Context, idCliente int64, t *model.TelefoneCliente) (*model.TelefoneCliente, error) {
	if err := t.Validar(); err != nil {
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

	telefoneCriado, err := s.repository.Clientes.CriarTelefone(ctx, tx, idCliente, t)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return telefoneCriado, nil
}

func (s *ClienteService) EditarTelefone(ctx context.Context, idCliente int64, idTelefone int64, t *model.TelefoneCliente) error {
	if err := t.Validar(); err != nil {
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

	err = s.repository.Clientes.EditarTelefone(ctx, tx, idCliente, idTelefone, t)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ClienteService) BuscarTelefoneByID(ctx context.Context, idCliente int64, idTelefone int64) (*model.TelefoneCliente, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	telefone, err := s.repository.Clientes.BuscarTelefoneByID(ctx, tx, idCliente, idTelefone)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return telefone, nil
}
