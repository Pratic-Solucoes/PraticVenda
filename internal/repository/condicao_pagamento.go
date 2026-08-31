package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"

	"github.com/lib/pq"
)

type CondicaoPagamentoRepository struct {
	db *sql.DB
}

func (r *CondicaoPagamentoRepository) Criar(ctx context.Context, tx *sql.Tx, cp *model.CondicaoPagamento) error {

	query := `
		insert into tb_condicoes_pagamento(descricao, qtd_parcelas, dias_primeiro_venc, intervalo_parcelas)
		values($1, $2, $3, $4)
		returning id
	`

	if err := tx.QueryRowContext(ctx, query, cp.Descricao, cp.QtdParcelas, cp.DiasPrimeiroVenc, cp.IntervaloParcelas).Scan(&cp.ID); err != nil {
		return err
	}

	return r.atualizarFormasPagamento(ctx, tx, cp.ID, cp.FormasPagamento)
}

func (r *CondicaoPagamentoRepository) Listar(ctx context.Context, tx *sql.Tx) ([]model.CondicaoPagamento, error) {
	rows, err := tx.QueryContext(ctx, consultaCondicoesPagamento+" ORDER BY cp.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	condicoes := make([]model.CondicaoPagamento, 0)
	for rows.Next() {
		condicao, err := lerCondicaoPagamento(rows)
		if err != nil {
			return nil, err
		}
		condicoes = append(condicoes, *condicao)
	}

	return condicoes, rows.Err()
}

func (r *CondicaoPagamentoRepository) BuscarPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.CondicaoPagamento, error) {
	row := tx.QueryRowContext(ctx, consultaCondicoesPagamento+" HAVING cp.id = $1", id)
	return lerCondicaoPagamento(row)
}

func (r *CondicaoPagamentoRepository) Atualizar(ctx context.Context, tx *sql.Tx, cp *model.CondicaoPagamento) error {
	const query = `
		UPDATE tb_condicoes_pagamento
		SET descricao = $1, qtd_parcelas = $2, dias_primeiro_venc = $3, intervalo_parcelas = $4
		WHERE id = $5
	`
	result, err := tx.ExecContext(ctx, query, cp.Descricao, cp.QtdParcelas, cp.DiasPrimeiroVenc, cp.IntervaloParcelas, cp.ID)
	if err != nil {
		return err
	}

	afetadas, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if afetadas == 0 {
		return sql.ErrNoRows
	}

	return r.atualizarFormasPagamento(ctx, tx, cp.ID, cp.FormasPagamento)
}

func (r *CondicaoPagamentoRepository) atualizarFormasPagamento(ctx context.Context, tx *sql.Tx, idCondicao uint64, formas []uint64) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM tb_condicao_forma_pagamento WHERE id_condicao = $1`, idCondicao); err != nil {
		return err
	}

	const query = `
		INSERT INTO tb_condicao_forma_pagamento (id_condicao, id_forma_pagamento)
		SELECT $1, forma_pagamento_id
		FROM unnest($2::bigint[]) AS forma_pagamento_id
	`
	idsFormas := make([]int64, len(formas))
	for indice, idForma := range formas {
		idsFormas[indice] = int64(idForma)
	}
	result, err := tx.ExecContext(ctx, query, idCondicao, pq.Array(idsFormas))
	if err != nil {
		return err
	}

	afetadas, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if afetadas != int64(len(formas)) {
		return sql.ErrNoRows
	}
	return nil
}

const consultaCondicoesPagamento = `
	SELECT cp.id, cp.descricao, cp.qtd_parcelas, cp.dias_primeiro_venc, cp.intervalo_parcelas,
		COALESCE(array_agg(cfp.id_forma_pagamento ORDER BY cfp.id_forma_pagamento)
		FILTER (WHERE cfp.id_forma_pagamento IS NOT NULL), '{}')
	FROM tb_condicoes_pagamento cp
	LEFT JOIN tb_condicao_forma_pagamento cfp ON cfp.id_condicao = cp.id
	GROUP BY cp.id, cp.descricao, cp.qtd_parcelas, cp.dias_primeiro_venc, cp.intervalo_parcelas
`

type scannerCondicaoPagamento interface {
	Scan(dest ...any) error
}

func lerCondicaoPagamento(scanner scannerCondicaoPagamento) (*model.CondicaoPagamento, error) {
	var condicao model.CondicaoPagamento
	var formas []int64
	if err := scanner.Scan(
		&condicao.ID,
		&condicao.Descricao,
		&condicao.QtdParcelas,
		&condicao.DiasPrimeiroVenc,
		&condicao.IntervaloParcelas,
		pq.Array(&formas),
	); err != nil {
		return nil, err
	}

	condicao.FormasPagamento = make([]uint64, len(formas))
	for indice, id := range formas {
		condicao.FormasPagamento[indice] = uint64(id)
	}
	return &condicao, nil
}
