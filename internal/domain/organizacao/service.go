package organizacao

import (
	"context"
	"database/sql"
)

type ServiceInterface interface {
	Criar(ctx context.Context, organizacao Organizacao) error
	Editar(ctx context.Context, organizacao Organizacao) error
}

type Service struct {
	repository RepositoryInterface
	db         *sql.DB
}

func NewService(repository RepositoryInterface, db *sql.DB) ServiceInterface {
	return &Service{
		repository: repository,
		db:         db,
	}
}

func (s *Service) Criar(ctx context.Context, organizacao Organizacao) error {

	if err := organizacao.Validar(); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := s.repository.Criar(ctx, tx, &organizacao); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Editar(ctx context.Context, organizacao Organizacao) error {
	return nil
}
