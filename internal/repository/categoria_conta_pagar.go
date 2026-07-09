package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type CategoriaContaPagarRepository struct {
	db *sql.DB
}

func NovoCategoriaContaPagarRepository(db *sql.DB) *CategoriaContaPagarRepository {
	return &CategoriaContaPagarRepository{db: db}
}

func (r *CategoriaContaPagarRepository) CriarCategoria(ctx context.Context, tx *sql.Tx, c *model.CategoriaContaPagar) (*model.CategoriaContaPagar, error) {
	query := `INSERT INTO tb_categorias_contas_pagar (nome) VALUES ($1) RETURNING id;`

	err := tx.QueryRowContext(ctx, query, c.Nome).Scan(&c.ID)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CategoriaContaPagarRepository) ListarCategorias(ctx context.Context, tx *sql.Tx) ([]*model.CategoriaContaPagar, error) {
	query := `SELECT id, nome FROM tb_categorias_contas_pagar ORDER BY nome ASC`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categorias []*model.CategoriaContaPagar
	for rows.Next() {
		c := &model.CategoriaContaPagar{}
		if err := rows.Scan(&c.ID, &c.Nome); err != nil {
			return nil, err
		}
		categorias = append(categorias, c)
	}
	return categorias, nil
}
