package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"time"
)

type VendaPDVService struct{ db *sql.DB }
type linhaVenda struct {
	id       int64
	q, preco float64
}

func (s *VendaPDVService) SalvarPreVenda(ctx context.Context, v model.VendaPDV, usuario int64) (int64, error) {
	if err := v.ValidarItens(); err != nil {
		return 0, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	estoque, _, _, _, _, _, err := s.configuracao(ctx, tx)
	if err != nil {
		return 0, err
	}
	linhas, total, err := s.calcularItens(ctx, tx, v, false)
	if err != nil {
		return 0, err
	}
	id, err := s.gravarVenda(ctx, tx, v, usuario, estoque, nil, nil, nil, total, "ABERTA", linhas)
	if err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (s *VendaPDVService) ListarPreVendas(ctx context.Context) ([]model.PreVendaPDV, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT v.id,v.id_cliente,COALESCE(c.nome,''),COALESCE(v.apelido_consumidor,''),v.valor_total FROM tb_vendas_pdv v LEFT JOIN tb_clientes c ON c.id=v.id_cliente WHERE v.status='ABERTA' ORDER BY v.criado_em DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.PreVendaPDV{}
	for rows.Next() {
		var p model.PreVendaPDV
		if err := rows.Scan(&p.ID, &p.IDCliente, &p.Cliente, &p.ApelidoConsumidor, &p.ValorTotal); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}
func (s *VendaPDVService) ObterPreVenda(ctx context.Context, id int64) (*model.PreVendaPDV, error) {
	p := &model.PreVendaPDV{ID: id}
	err := s.db.QueryRowContext(ctx, `SELECT v.id_cliente,COALESCE(c.nome,''),COALESCE(v.apelido_consumidor,''),v.valor_total FROM tb_vendas_pdv v LEFT JOIN tb_clientes c ON c.id=v.id_cliente WHERE v.id=$1 AND v.status='ABERTA'`, id).Scan(&p.IDCliente, &p.Cliente, &p.ApelidoConsumidor, &p.ValorTotal)
	if err == sql.ErrNoRows {
		return nil, errors.New("pré-venda não encontrada ou não está aberta")
	}
	if err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id_produto,quantidade,valor_unitario FROM tb_itens_vendas_pdv WHERE id_venda=$1 ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i model.ItemVendaPDV
		if err = rows.Scan(&i.IDProduto, &i.Quantidade, &i.ValorUnitario); err != nil {
			return nil, err
		}
		p.Itens = append(p.Itens, i)
	}
	return p, rows.Err()
}

func (s *VendaPDVService) Finalizar(ctx context.Context, v model.VendaPDV, usuario int64) (int64, error) {
	if err := v.Validar(); err != nil {
		return 0, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var controle int64
	if err = tx.QueryRowContext(ctx, `SELECT id FROM tb_controle_caixa WHERE id_usuario=$1 AND status='ABERTO' ORDER BY aberto_em DESC LIMIT 1`, usuario).Scan(&controle); err != nil {
		return 0, errors.New("abra o caixa antes de finalizar uma venda")
	}
	estoque, categoria, exigir, pDesconto, pPreco, recebido, err := s.configuracao(ctx, tx)
	if err != nil {
		return 0, err
	}
	tipo, parcelas, permitida, err := s.pagamento(ctx, tx, v)
	if err != nil || !permitida {
		return 0, errors.New("forma ou condição de pagamento não está liberada para o PDV")
	}
	if exigir && parcelas > 1 && (v.IDCliente == nil || *v.IDCliente <= 0) {
		return 0, errors.New("cliente é obrigatório para vendas parceladas/a prazo")
	}
	linhas, total, err := s.calcularItens(ctx, tx, v, pPreco)
	if err != nil {
		return 0, err
	}
	if v.ValorDesconto > 0 && !pDesconto {
		return 0, errors.New("desconto manual não está permitido no PDV")
	}
	var id int64
	if v.ID > 0 {
		err = tx.QueryRowContext(ctx, `UPDATE tb_vendas_pdv SET id_cliente=$1,apelido_consumidor=$2,id_controle_caixa=$3,id_forma_pagamento=$4,id_condicao_pagamento=$5,valor_desconto=$6,valor_total=$7,status='CONCLUIDA',concluido_em=NOW() WHERE id=$8 AND status='ABERTA' RETURNING id`, v.IDCliente, v.ApelidoConsumidor, controle, v.IDFormaPagamento, v.IDCondicaoPagamento, v.ValorDesconto, total, v.ID).Scan(&id)
		if err == sql.ErrNoRows {
			return 0, errors.New("pré-venda não encontrada ou já processada")
		}
	} else {
		id, err = s.gravarVenda(ctx, tx, v, usuario, estoque, &controle, &v.IDFormaPagamento, &v.IDCondicaoPagamento, total, "CONCLUIDA", linhas)
	}
	if err != nil {
		return 0, err
	}
	if err = s.baixarEstoque(ctx, tx, linhas, estoque, usuario, id, "VENDA PDV"); err != nil {
		return 0, err
	}
	if recebido {
		if _, err = tx.ExecContext(ctx, `INSERT INTO tb_movimento_caixa(id_controle_caixa,tipo_movimento,id_venda,id_cliente,tipo_credito,valor_credito) VALUES($1,'VENDA',$2,$3,$4,$5)`, controle, id, v.IDCliente, tipo, total); err != nil {
			return 0, err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO tb_movimento_financeiro(tipo_movimento,id_venda_pdv,id_categoria_credito,dt_movimento,valor_movimento,forma_pagamento) SELECT 'VENDA_PDV',$1,$2,CURRENT_DATE,$3,descricao FROM tb_formas_pagamento WHERE id=$4`, id, categoria, total, v.IDFormaPagamento); err != nil {
			return 0, err
		}
	}
	if err = s.registrarContaReceberPDV(ctx, tx, id, v, categoria, total, recebido); err != nil {
		return 0, err
	}
	return id, tx.Commit()
}
func (s *VendaPDVService) Cancelar(ctx context.Context, id, usuario int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var estoque int64
	var status string
	if err = tx.QueryRowContext(ctx, `SELECT id_estoque,status FROM tb_vendas_pdv WHERE id=$1 FOR UPDATE`, id).Scan(&estoque, &status); err == sql.ErrNoRows {
		return errors.New("venda não encontrada")
	}
	if err != nil {
		return err
	}
	if status == "CANCELADA" {
		return errors.New("venda já está cancelada")
	}
	if status == "ABERTA" {
		_, err = tx.ExecContext(ctx, `UPDATE tb_vendas_pdv SET status='CANCELADA',cancelado_em=NOW(),id_usuario_cancelamento=$1 WHERE id=$2`, usuario, id)
		if err != nil {
			return err
		}
		return tx.Commit()
	}
	rows, err := tx.QueryContext(ctx, `SELECT id_produto,quantidade FROM tb_itens_vendas_pdv WHERE id_venda=$1`, id)
	if err != nil {
		return err
	}
	var linhas []linhaVenda
	for rows.Next() {
		var l linhaVenda
		if err = rows.Scan(&l.id, &l.q); err != nil {
			rows.Close()
			return err
		}
		linhas = append(linhas, l)
	}
	rows.Close()
	if err = s.reporEstoque(ctx, tx, linhas, estoque, usuario, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO tb_movimento_caixa(id_controle_caixa,tipo_movimento,id_venda,id_cliente,tipo_credito,valor_credito) SELECT id_controle_caixa,'ESTORNO',id,id_cliente,fp.tipo,valor_total FROM tb_vendas_pdv v JOIN tb_formas_pagamento fp ON fp.id=v.id_forma_pagamento WHERE v.id=$1 AND EXISTS(SELECT 1 FROM tb_movimento_caixa m WHERE m.id_venda=v.id AND m.tipo_movimento='VENDA')`, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO tb_movimento_financeiro(tipo_movimento,id_venda_pdv,id_categoria_credito,dt_movimento,valor_movimento,forma_pagamento) SELECT 'ESTORNO_VENDA_PDV',m.id_venda_pdv,m.id_categoria_credito,CURRENT_DATE,m.valor_movimento,m.forma_pagamento FROM tb_movimento_financeiro m WHERE m.id_venda_pdv=$1 AND m.tipo_movimento='VENDA_PDV'`, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE tb_contas_receber SET status='CANCELADO',saldo_restante=0,updated_at=NOW() WHERE tipo_origem='VENDA_PDV' AND id_origem=$1 AND status IN ('PENDENTE','PAGO_PARCIAL','PAGO')`, id); err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE tb_vendas_pdv SET status='CANCELADA',cancelado_em=NOW(),id_usuario_cancelamento=$1 WHERE id=$2`, usuario, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *VendaPDVService) registrarContaReceberPDV(ctx context.Context, tx *sql.Tx, venda int64, v model.VendaPDV, categoria int64, total float64, recebido bool) error {
	status := "PENDENTE"
	saldo := total
	var pagamento any
	if recebido {
		status = "PAGO"
		saldo = 0
		pagamento = time.Now()
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO tb_contas_receber(id_cliente,id_categoria,id_condicao_pagamento,id_forma_pagamento,descricao,valor_original,saldo_restante,dt_vencimento,nr_parcela,nr_total_parcelas,status,dt_emissao,dt_pagamento,tipo_origem,id_origem) SELECT $1,$2,$3,$4,'VENDA PDV - ' || $5::text,$6,$7,CURRENT_DATE + (cp.dias_primeiro_venc || ' days')::interval,1,cp.qtd_parcelas,$8,CURRENT_TIMESTAMP,$9,'VENDA_PDV',$5 FROM tb_condicoes_pagamento cp WHERE cp.id=$3`, v.IDCliente, categoria, v.IDCondicaoPagamento, v.IDFormaPagamento, venda, total, saldo, status, pagamento)
	return err
}
func (s *VendaPDVService) configuracao(ctx context.Context, tx *sql.Tx) (int64, int64, bool, bool, bool, bool, error) {
	var a, b int64
	var c, d, e, f bool
	err := tx.QueryRowContext(ctx, `SELECT id_estoque_padrao,id_categoria_credito,exigir_cliente_prazo,permitir_desconto_manual,permitir_alterar_preco,gerar_financeiro_recebido FROM tb_configuracoes_pdv WHERE id=1`).Scan(&a, &b, &c, &d, &e, &f)
	if err == sql.ErrNoRows {
		return 0, 0, false, false, false, false, errors.New("configure o PDV antes de realizar vendas")
	}
	return a, b, c, d, e, f, err
}
func (s *VendaPDVService) pagamento(ctx context.Context, tx *sql.Tx, v model.VendaPDV) (string, int, bool, error) {
	var t string
	var n int
	var ok bool
	err := tx.QueryRowContext(ctx, `SELECT fp.tipo,cp.qtd_parcelas,EXISTS(SELECT 1 FROM tb_configuracoes_pdv_formas_pagamento cf WHERE cf.id_forma_pagamento=fp.id) AND EXISTS(SELECT 1 FROM tb_configuracoes_pdv_condicoes_pagamento cc WHERE cc.id_condicao_pagamento=cp.id) AND EXISTS(SELECT 1 FROM tb_condicao_forma_pagamento cfp WHERE cfp.id_condicao=cp.id AND cfp.id_forma_pagamento=fp.id) FROM tb_formas_pagamento fp CROSS JOIN tb_condicoes_pagamento cp WHERE fp.id=$1 AND cp.id=$2`, v.IDFormaPagamento, v.IDCondicaoPagamento).Scan(&t, &n, &ok)
	return t, n, ok, err
}
func (s *VendaPDVService) calcularItens(ctx context.Context, tx *sql.Tx, v model.VendaPDV, editar bool) ([]linhaVenda, float64, error) {
	var linhas []linhaVenda
	var total float64
	for _, it := range v.Itens {
		var p float64
		if err := tx.QueryRowContext(ctx, `SELECT preco_venda FROM tb_produtos WHERE id=$1 AND ativo=true AND materia_prima=false`, it.IDProduto).Scan(&p); err != nil {
			return nil, 0, errors.New("produto inválido ou matéria-prima não vendável")
		}
		if it.ValorUnitario != nil {
			if !editar {
				return nil, 0, errors.New("alteração de preço não está permitida no PDV")
			}
			if *it.ValorUnitario < 0 {
				return nil, 0, errors.New("preço do item inválido")
			}
			p = *it.ValorUnitario
		}
		linhas = append(linhas, linhaVenda{it.IDProduto, it.Quantidade, p})
		total += p * it.Quantidade
	}
	if v.ValorDesconto > total {
		return nil, 0, errors.New("o desconto não pode ser maior que o total da venda")
	}
	return linhas, total - v.ValorDesconto, nil
}
func (s *VendaPDVService) gravarVenda(ctx context.Context, tx *sql.Tx, v model.VendaPDV, u, estoque int64, controle, forma, condicao *int64, total float64, status string, linhas []linhaVenda) (int64, error) {
	var id int64
	err := tx.QueryRowContext(ctx, `INSERT INTO tb_vendas_pdv(id_estoque,id_usuario,id_controle_caixa,id_cliente,id_forma_pagamento,id_condicao_pagamento,apelido_consumidor,valor_desconto,valor_total,status,concluido_em) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,CASE WHEN $11 THEN NOW() END) RETURNING id`, estoque, u, controle, v.IDCliente, forma, condicao, v.ApelidoConsumidor, v.ValorDesconto, total, status, status == "CONCLUIDA").Scan(&id)
	if err != nil {
		return 0, err
	}
	for _, l := range linhas {
		if _, err = tx.ExecContext(ctx, `INSERT INTO tb_itens_vendas_pdv(id_venda,id_produto,quantidade,valor_unitario,valor_total) VALUES($1,$2,$3,$4,$5)`, id, l.id, l.q, l.preco, l.q*l.preco); err != nil {
			return 0, err
		}
	}
	return id, nil
}
func (s *VendaPDVService) baixarEstoque(ctx context.Context, tx *sql.Tx, ls []linhaVenda, estoque, u, venda int64, categoria string) error {
	for _, l := range ls {
		var composto bool
		if err := tx.QueryRowContext(ctx, `SELECT composto FROM tb_produtos WHERE id=$1`, l.id).Scan(&composto); err != nil {
			return errors.New("produto inválido")
		}
		rows, err := tx.QueryContext(ctx, `SELECT COALESCE(c.id_produto_componente,$1),COALESCE(c.quantidade,1)*$2 FROM (SELECT 1) x LEFT JOIN tb_composicoes_produtos c ON c.id_produto_composto=$1`, l.id, l.q)
		if err != nil {
			return err
		}
		componentes := []linhaVenda{}
		for rows.Next() {
			var p int64
			var q float64
			if err = rows.Scan(&p, &q); err != nil {
				rows.Close()
				return err
			}
			componentes = append(componentes, linhaVenda{id: p, q: q})
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			return err
		}
		rows.Close()

		for _, componente := range componentes {
			p, q := componente.id, componente.q
			if composto {
				// Produtos compostos não têm saldo próprio. Seus componentes são baixados
				// para manter o custo/histórico, mas a venda não é bloqueada por saldo zero.
				if _, err = tx.ExecContext(ctx, `INSERT INTO tb_produtos_estoque(id_produto,id_estoque,quantidade,estoque_minimo,atualizado_em) VALUES($1,$2,-($3::numeric),0,NOW()) ON CONFLICT(id_produto,id_estoque) DO UPDATE SET quantidade=tb_produtos_estoque.quantidade-($3::numeric),atualizado_em=NOW()`, p, estoque, q); err != nil {
					return err
				}
			} else {
				var saldo float64
				if err = tx.QueryRowContext(ctx, `SELECT quantidade FROM tb_produtos_estoque WHERE id_produto=$1 AND id_estoque=$2 FOR UPDATE`, p, estoque).Scan(&saldo); err != nil || saldo < q {
					return errors.New("saldo insuficiente para concluir a venda")
				}
				if _, err = tx.ExecContext(ctx, `UPDATE tb_produtos_estoque SET quantidade=quantidade-$1,atualizado_em=NOW() WHERE id_produto=$2 AND id_estoque=$3`, q, p, estoque); err != nil {
					return err
				}
			}
			if _, err = tx.ExecContext(ctx, `INSERT INTO tb_movimento_estoque(id_produto,id_estoque,id_usuario,quantidade,tipo_movimento,id_categoria_movimento,id_origem) SELECT $1,$2,$3,$4,'SAIDA',id,$5 FROM tb_categoria_movimento_estoque WHERE nome=$6`, p, estoque, u, q, venda, categoria); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *VendaPDVService) reporEstoque(ctx context.Context, tx *sql.Tx, ls []linhaVenda, estoque, u, venda int64) error {
	for _, l := range ls {
		rows, err := tx.QueryContext(ctx, `SELECT COALESCE(c.id_produto_componente,$1),COALESCE(c.quantidade,1)*$2 FROM (SELECT 1) x LEFT JOIN tb_composicoes_produtos c ON c.id_produto_composto=$1`, l.id, l.q)
		if err != nil {
			return err
		}
		componentes := []linhaVenda{}
		for rows.Next() {
			var p int64
			var q float64
			if err = rows.Scan(&p, &q); err != nil {
				rows.Close()
				return err
			}
			componentes = append(componentes, linhaVenda{id: p, q: q})
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			return err
		}
		rows.Close()

		for _, componente := range componentes {
			p, q := componente.id, componente.q
			if _, err = tx.ExecContext(ctx, `UPDATE tb_produtos_estoque SET quantidade=quantidade+$1,atualizado_em=NOW() WHERE id_produto=$2 AND id_estoque=$3`, q, p, estoque); err != nil {
				return err
			}
			if _, err = tx.ExecContext(ctx, `INSERT INTO tb_movimento_estoque(id_produto,id_estoque,id_usuario,quantidade,tipo_movimento,id_categoria_movimento,id_origem) SELECT $1,$2,$3,$4,'ENTRADA',id,$5 FROM tb_categoria_movimento_estoque WHERE nome='CANCELAMENTO VENDA PDV'`, p, estoque, u, q, venda); err != nil {
				return err
			}
		}
	}
	return nil
}
