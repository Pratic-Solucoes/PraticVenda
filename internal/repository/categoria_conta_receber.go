package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type CategoriaContaReceberRepository struct{ db *sql.DB }

func (r *CategoriaContaReceberRepository) CriarCategoria(ctx context.Context, tx *sql.Tx, c *model.CategoriaContaReceber) (*model.CategoriaContaReceber, error) {
	err := tx.QueryRowContext(ctx, `INSERT INTO tb_categorias_contas_receber (descricao) VALUES ($1) RETURNING id`, c.Nome).Scan(&c.ID)
	return c, err
}

func (r *CategoriaContaReceberRepository) ListarCategorias(ctx context.Context, tx *sql.Tx) ([]model.CategoriaContaReceber, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id, descricao FROM tb_categorias_contas_receber ORDER BY descricao`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categorias []model.CategoriaContaReceber
	for rows.Next() {
		var c model.CategoriaContaReceber
		if err := rows.Scan(&c.ID, &c.Nome); err != nil {
			return nil, err
		}
		categorias = append(categorias, c)
	}
	return categorias, rows.Err()
}
