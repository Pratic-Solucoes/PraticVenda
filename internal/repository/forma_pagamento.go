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

	query := `insert into tb_formas_pagamento(descricao) values ($1) returning id`

	var id uint64
	err := tx.QueryRowContext(ctx, query, fp.Descricao).Scan(&id)
	if err != nil {
		return nil, err
	}

	fp.ID = id

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

func (r *FormaPagamentoRepository) BuscarPorID(ctx context.Context, tx *sql.Tx, idFp int64) (*model.FormaPagamento, error) {

	query := `select id, descricao from tb_formas_pagamento where id = $1`

	var fp model.FormaPagamento
	err := tx.QueryRowContext(ctx, query, idFp).Scan(&fp.ID, &fp.Descricao)
	if err != nil {
		return nil, err
	}

	return &fp, nil
}

func (r *FormaPagamentoRepository) Atualizar(ctx context.Context, tx *sql.Tx, fp *model.FormaPagamento) (*model.FormaPagamento, error) {

	query := `update tb_formas_pagamento set descricao = $1 where id = $2 returning id, descricao`

	var id uint64
	err := tx.QueryRowContext(ctx, query, fp.Descricao, fp.ID).Scan(&id, &fp.Descricao)
	if err != nil {
		return nil, err
	}

	fp.ID = id

	return fp, nil
}
