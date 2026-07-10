package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type FormaPagamentoRepository struct {
	db *sql.DB
}

func (r *FormaPagamentoRepository) Criar(ctx context.Context, tx *sql.Tx, fp *model.FormaPagamento) (*model.FormaPagamento, error) {

	query := `insert into tb_formas_pagamento(descricao) values ($1)`

	result, err := tx.ExecContext(ctx, query, fp.Descricao)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	fp.ID = uint64(id)

	return fp, nil
}

func (r *FormaPagamentoRepository) Listar(ctx context.Context, tx *sql.Tx) ([]model.FormaPagamento, error) {
	query := `select id, descricao from tb_formas_pagamento order by id`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var formas []model.FormaPagamento

	for rows.Next() {
		var fp model.FormaPagamento
		if err := rows.Scan(&fp.ID, &fp.Descricao); err != nil {
			return nil, err
		}
		formas = append(formas, fp)
	}
	return formas, nil
}
