package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type CategoriaRepository struct {
	db *sql.DB
}

func NovoCategoriaRepository(db *sql.DB) *CategoriaRepository {
	return &CategoriaRepository{db: db}
}

func (r *CategoriaRepository) CriarCategoria(ctx context.Context, tx *sql.Tx, c *model.CategoriaDebito) (*model.CategoriaDebito, error) {
	query := `INSERT INTO tb_categorias_debito (nome) VALUES ($1) RETURNING id;`

	err := tx.QueryRowContext(ctx, query, c.Nome).Scan(&c.ID)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CategoriaRepository) ListarCategorias(ctx context.Context, tx *sql.Tx) ([]*model.CategoriaDebito, error) {
	query := `SELECT id, nome FROM tb_categorias_debito ORDER BY nome ASC`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categorias []*model.CategoriaDebito
	for rows.Next() {
		c := &model.CategoriaDebito{}
		if err := rows.Scan(&c.ID, &c.Nome); err != nil {
			return nil, err
		}
		categorias = append(categorias, c)
	}
	return categorias, nil
}
