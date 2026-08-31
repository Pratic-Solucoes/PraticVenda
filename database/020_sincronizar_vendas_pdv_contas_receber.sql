BEGIN;

-- Cria contas a receber para vendas de PDV já concluídas antes da implantação
-- do espelho financeiro. Vendas com movimento financeiro de PDV já existente
-- são registradas como pagas, preservando valor original e saldo zero.
INSERT INTO tb_contas_receber (
    id_cliente, id_categoria, id_condicao_pagamento, id_forma_pagamento,
    descricao, valor_original, saldo_restante, dt_vencimento,
    nr_parcela, nr_total_parcelas, status, dt_emissao, dt_pagamento,
    tipo_origem, id_origem
)
SELECT
    v.id_cliente,
    cfg.id_categoria_credito,
    v.id_condicao_pagamento,
    v.id_forma_pagamento,
    'VENDA PDV - ' || v.id::text,
    v.valor_total,
    CASE WHEN EXISTS (
        SELECT 1 FROM tb_movimento_financeiro mf
        WHERE mf.id_venda_pdv = v.id AND mf.tipo_movimento = 'VENDA_PDV'
    ) THEN 0 ELSE v.valor_total END,
    COALESCE(v.concluido_em::date, v.criado_em::date),
    1,
    cp.qtd_parcelas,
    CASE WHEN EXISTS (
        SELECT 1 FROM tb_movimento_financeiro mf
        WHERE mf.id_venda_pdv = v.id AND mf.tipo_movimento = 'VENDA_PDV'
    ) THEN 'PAGO' ELSE 'PENDENTE' END,
    COALESCE(v.concluido_em, v.criado_em),
    CASE WHEN EXISTS (
        SELECT 1 FROM tb_movimento_financeiro mf
        WHERE mf.id_venda_pdv = v.id AND mf.tipo_movimento = 'VENDA_PDV'
    ) THEN COALESCE(v.concluido_em, v.criado_em) ELSE NULL END,
    'VENDA_PDV',
    v.id
FROM tb_vendas_pdv v
JOIN tb_configuracoes_pdv cfg ON cfg.id = 1
JOIN tb_condicoes_pagamento cp ON cp.id = v.id_condicao_pagamento
WHERE v.status = 'CONCLUIDA'
  AND NOT EXISTS (
      SELECT 1 FROM tb_contas_receber cr
      WHERE cr.tipo_origem = 'VENDA_PDV' AND cr.id_origem = v.id
  );

COMMIT;
