package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
)

type ConfiguracaoPDVService struct{ db *sql.DB }

func (s *ConfiguracaoPDVService) Obter(ctx context.Context) (*model.ConfiguracaoPDV, error) {
	c := &model.ConfiguracaoPDV{FormasPagamento: []int64{}, CondicoesPagamento: []int64{}}
	err := s.db.QueryRowContext(ctx, `SELECT id_estoque_padrao,id_categoria_credito,exigir_cliente_prazo,permitir_desconto_manual,permitir_alterar_preco,gerar_financeiro_recebido FROM tb_configuracoes_pdv WHERE id=1`).Scan(&c.IDEstoquePadrao, &c.IDCategoriaCredito, &c.ExigirClientePrazo, &c.PermitirDescontoManual, &c.PermitirAlterarPreco, &c.GerarFinanceiroRecebido)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	formas, err := s.ids(ctx, `SELECT id_forma_pagamento FROM tb_configuracoes_pdv_formas_pagamento ORDER BY id_forma_pagamento`)
	if err != nil {
		return nil, err
	}
	condicoes, err := s.ids(ctx, `SELECT id_condicao_pagamento FROM tb_configuracoes_pdv_condicoes_pagamento ORDER BY id_condicao_pagamento`)
	if err != nil {
		return nil, err
	}
	c.FormasPagamento, c.CondicoesPagamento = formas, condicoes
	return c, nil
}

func (s *ConfiguracaoPDVService) ids(ctx context.Context, query string) ([]int64, error) {
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *ConfiguracaoPDVService) Salvar(ctx context.Context, c model.ConfiguracaoPDV) error {
	if c.IDEstoquePadrao <= 0 || c.IDCategoriaCredito <= 0 || len(c.FormasPagamento) == 0 || len(c.CondicoesPagamento) == 0 {
		return errors.New("preencha o estoque, categoria, formas e condições de pagamento")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var valido bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tb_estoques WHERE id=$1 AND ativo=true) AND EXISTS(SELECT 1 FROM tb_categorias_contas_receber WHERE id=$2)`, c.IDEstoquePadrao, c.IDCategoriaCredito).Scan(&valido)
	if err != nil || !valido {
		if err != nil {
			return err
		}
		return errors.New("estoque ou categoria de crédito inválidos")
	}
	if err = validarIDs(ctx, tx, `tb_formas_pagamento`, `id`, c.FormasPagamento); err != nil {
		return err
	}
	if err = validarIDs(ctx, tx, `tb_condicoes_pagamento`, `id`, c.CondicoesPagamento); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO tb_configuracoes_pdv(id,id_estoque_padrao,id_categoria_credito,exigir_cliente_prazo,permitir_desconto_manual,permitir_alterar_preco,gerar_financeiro_recebido) VALUES(1,$1,$2,$3,$4,$5,$6) ON CONFLICT(id) DO UPDATE SET id_estoque_padrao=EXCLUDED.id_estoque_padrao,id_categoria_credito=EXCLUDED.id_categoria_credito,exigir_cliente_prazo=EXCLUDED.exigir_cliente_prazo,permitir_desconto_manual=EXCLUDED.permitir_desconto_manual,permitir_alterar_preco=EXCLUDED.permitir_alterar_preco,gerar_financeiro_recebido=EXCLUDED.gerar_financeiro_recebido,updated_at=NOW()`, c.IDEstoquePadrao, c.IDCategoriaCredito, c.ExigirClientePrazo, c.PermitirDescontoManual, c.PermitirAlterarPreco, c.GerarFinanceiroRecebido); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM tb_configuracoes_pdv_formas_pagamento`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM tb_configuracoes_pdv_condicoes_pagamento`); err != nil {
		return err
	}
	for _, id := range c.FormasPagamento {
		if _, err = tx.ExecContext(ctx, `INSERT INTO tb_configuracoes_pdv_formas_pagamento(id_forma_pagamento) VALUES($1)`, id); err != nil {
			return err
		}
	}
	for _, id := range c.CondicoesPagamento {
		if _, err = tx.ExecContext(ctx, `INSERT INTO tb_configuracoes_pdv_condicoes_pagamento(id_condicao_pagamento) VALUES($1)`, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func validarIDs(ctx context.Context, tx *sql.Tx, tabela, coluna string, ids []int64) error {
	for _, id := range ids {
		var existe bool
		if id <= 0 {
			return errors.New("item de configuração inválido")
		}
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM `+tabela+` WHERE `+coluna+`=$1)`, id).Scan(&existe); err != nil {
			return err
		}
		if !existe {
			return errors.New("item de configuração não encontrado")
		}
	}
	return nil
}
