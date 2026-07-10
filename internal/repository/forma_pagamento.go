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
