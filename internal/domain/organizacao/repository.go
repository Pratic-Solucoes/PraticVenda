package organizacao

import (
	"context"
	"database/sql"
)

type RepositoryInterface interface {
	Criar(ctx context.Context, tx *sql.Tx, organizacao *Organizacao) error
	Editar(ctx context.Context, tx *sql.Tx, organizacao *Organizacao) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RepositoryInterface {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Criar(ctx context.Context, tx *sql.Tx, organizacao *Organizacao) error {

	query := `
		insert into tb_organizacoes (id_dono, nome_organizacao)
		values($1, $2)
		returning id;
	`
	err := tx.QueryRowContext(
		ctx,
		query,
		organizacao.IDDono,
		organizacao.NomeOrganizacao,
	).Scan(&organizacao.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) BuscarPorID(ctx context.Context, tx *sql.Tx, id uint64) (*Organizacao, error) {
	return nil, nil
}

func (r *Repository) Editar(ctx context.Context, tx *sql.Tx, organizacao *Organizacao) error {
	return nil
}
