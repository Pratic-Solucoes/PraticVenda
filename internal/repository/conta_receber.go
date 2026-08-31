package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gestao/internal/model"
)

var (
	CONTA_RECEBER_QUITADA        = errors.New("conta a receber já quitada")
	CONTA_RECEBER_NAO_ENCONTRADA = errors.New("conta a receber não encontrada")
)

type ContaReceberRepository struct{ db *sql.DB }

func (r *ContaReceberRepository) CriarParcelas(ctx context.Context, tx *sql.Tx, conta *model.ContaReceberCriar, grupo string, valores []float64, vencimentos []string) error {
	const query = `INSERT INTO tb_contas_receber (id_cliente,id_categoria,id_condicao_pagamento,id_forma_pagamento,id_grupo_parcelas,descricao,valor_original,saldo_restante,dt_vencimento,nr_parcela,nr_total_parcelas,status,dt_emissao,tipo_origem)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$7,$8,$9,$10,'PENDENTE',$11,'LANCAMENTO_MANUAL')`
	for i, valor := range valores {
		if _, err := tx.ExecContext(ctx, query, conta.IDCliente, conta.IDCategoria, conta.IDCondicaoPagamento, conta.IDFormaPagamento, grupo, conta.Descricao, valor, vencimentos[i], i+1, len(valores), conta.DtEmissao); err != nil {
			return err
		}
	}
	return nil
}

func (r *ContaReceberRepository) Listar(ctx context.Context, tx *sql.Tx, busca, vencimento, status string) ([]*model.ContaReceber, error) {
	// As contas a receber são o registro principal desta tela. As referências de
	// cadastro podem estar ausentes em lançamentos legados, portanto não podem
	// removê-los da listagem.
	query := `SELECT cr.id,COALESCE(cr.id_cliente,0),COALESCE(cr.id_categoria,0),COALESCE(cr.id_condicao_pagamento,0),COALESCE(cr.id_forma_pagamento,0),COALESCE(cr.id_grupo_parcelas::text,''),COALESCE(cr.descricao,''),cr.valor_original,cr.saldo_restante,cr.dt_vencimento,cr.nr_parcela,cr.nr_total_parcelas,cr.status,cr.dt_emissao,cr.dt_pagamento,
	COALESCE(c.id,0),COALESCE(c.nome,'Consumidor final'),COALESCE(cat.id,0),COALESCE(cat.descricao,''),COALESCE(cp.id,0),COALESCE(cp.descricao,''),COALESCE(cp.qtd_parcelas,0),COALESCE(cp.dias_primeiro_venc,0),COALESCE(cp.intervalo_parcelas,0),COALESCE(fp.id,0),COALESCE(fp.descricao,'')
	FROM tb_contas_receber cr
	LEFT JOIN tb_clientes c ON c.id=cr.id_cliente
	LEFT JOIN tb_categorias_contas_receber cat ON cat.id=cr.id_categoria
	LEFT JOIN tb_condicoes_pagamento cp ON cp.id=cr.id_condicao_pagamento
	LEFT JOIN tb_formas_pagamento fp ON fp.id=cr.id_forma_pagamento
	WHERE 1=1`
	args := []any{}
	if busca != "" {
		query += fmt.Sprintf(" AND (c.nome ILIKE $%d OR cr.id::text=$%d)", len(args)+1, len(args)+2)
		args = append(args, "%"+busca+"%", busca)
	}
	if vencimento != "" {
		query += fmt.Sprintf(" AND cr.dt_vencimento=$%d", len(args)+1)
		args = append(args, vencimento)
	}
	if status != "" {
		query += fmt.Sprintf(" AND cr.status=$%d", len(args)+1)
		args = append(args, status)
	}
	query += " ORDER BY cr.dt_vencimento, cr.id"
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contas []*model.ContaReceber
	for rows.Next() {
		d := &model.ContaReceber{}
		cliente := &model.Cliente{}
		categoria := &model.CategoriaContaReceber{}
		condicao := &model.CondicaoPagamento{}
		forma := &model.FormaPagamento{}
		if err := rows.Scan(&d.ID, &d.IDCliente, &d.IDCategoria, &d.IDCondicaoPagamento, &d.IDFormaPagamento, &d.IDGrupoParcelas, &d.Descricao, &d.ValorOriginal, &d.SaldoRestante, &d.DtVencimento, &d.NrParcela, &d.NrTotalParcelas, &d.Status, &d.DtEmissao, &d.DtPagamento, &cliente.ID, &cliente.Nome, &categoria.ID, &categoria.Nome, &condicao.ID, &condicao.Descricao, &condicao.QtdParcelas, &condicao.DiasPrimeiroVenc, &condicao.IntervaloParcelas, &forma.ID, &forma.Descricao); err != nil {
			return nil, err
		}
		d.Cliente = cliente
		d.Categoria = categoria
		d.CondicaoPagamento = condicao
		d.FormaPagamento = forma
		contas = append(contas, d)
	}
	return contas, rows.Err()
}

func (r *ContaReceberRepository) BuscarPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.ContaReceber, error) {
	d := &model.ContaReceber{}
	err := tx.QueryRowContext(ctx, `SELECT id,id_cliente,id_categoria,id_condicao_pagamento,id_forma_pagamento,COALESCE(id_grupo_parcelas::text,''),descricao,valor_original,saldo_restante,dt_vencimento,nr_parcela,nr_total_parcelas,status,dt_emissao,dt_pagamento FROM tb_contas_receber WHERE id=$1`, id).Scan(&d.ID, &d.IDCliente, &d.IDCategoria, &d.IDCondicaoPagamento, &d.IDFormaPagamento, &d.IDGrupoParcelas, &d.Descricao, &d.ValorOriginal, &d.SaldoRestante, &d.DtVencimento, &d.NrParcela, &d.NrTotalParcelas, &d.Status, &d.DtEmissao, &d.DtPagamento)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, CONTA_RECEBER_NAO_ENCONTRADA
	}
	return d, err
}

func (r *ContaReceberRepository) Receber(ctx context.Context, tx *sql.Tx, id int64, valor float64, idFormaPagamentoReal int64) error {
	conta, err := r.BuscarPorID(ctx, tx, id)
	if err != nil {
		return err
	}
	if conta.Status == "PAGO" {
		return CONTA_RECEBER_QUITADA
	}
	valorRestante := conta.SaldoRestante - valor
	if _, err = tx.ExecContext(ctx, `INSERT INTO tb_movimento_financeiro(tipo_movimento,id_conta_receber,dt_movimento,valor_movimento,forma_pagamento) SELECT 'CONTA_RECEBER',$1,CURRENT_DATE,$2,descricao FROM tb_formas_pagamento WHERE id=$3`, id, valor, idFormaPagamentoReal); err != nil {
		return err
	}
	// A parcela original é quitada. Quando o recebimento é parcial, o saldo é
	// preservado em uma nova parcela do mesmo grupo para manter o histórico da baixa.
	if _, err = tx.ExecContext(ctx, `UPDATE tb_contas_receber SET status='PAGO', saldo_restante=0, dt_pagamento=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id); err != nil {
		return err
	}
	if valorRestante <= 0 {
		return nil
	}
	novoTotal := conta.NrTotalParcelas + 1
	if conta.IDGrupoParcelas != "" {
		if _, err = tx.ExecContext(ctx, `UPDATE tb_contas_receber SET nr_total_parcelas=$1 WHERE id_grupo_parcelas=$2`, novoTotal, conta.IDGrupoParcelas); err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO tb_contas_receber (id_cliente,id_categoria,id_condicao_pagamento,id_forma_pagamento,id_grupo_parcelas,descricao,valor_original,saldo_restante,dt_vencimento,nr_parcela,nr_total_parcelas,status,dt_emissao,tipo_origem)
	VALUES ($1,$2,$3,$4,NULLIF($5,'')::uuid,$6,$7,$7,$8,$9,$10,'PENDENTE',$11,'SALDO_REMANESCENTE')`, conta.IDCliente, conta.IDCategoria, conta.IDCondicaoPagamento, conta.IDFormaPagamento, conta.IDGrupoParcelas, conta.Descricao, valorRestante, conta.DtVencimento, novoTotal, novoTotal, conta.DtEmissao)
	return err
}
