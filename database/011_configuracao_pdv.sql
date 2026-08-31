BEGIN;

CREATE TABLE IF NOT EXISTS tb_configuracoes_pdv (
    id SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    id_estoque_padrao BIGINT REFERENCES tb_estoques(id),
    id_categoria_credito BIGINT REFERENCES tb_categorias_contas_receber(id),
    exigir_cliente_prazo BOOLEAN NOT NULL DEFAULT TRUE,
    permitir_desconto_manual BOOLEAN NOT NULL DEFAULT FALSE,
    permitir_alterar_preco BOOLEAN NOT NULL DEFAULT FALSE,
    limite_diferenca_caixa DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (limite_diferenca_caixa >= 0),
    exigir_justificativa_diferenca BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_configuracoes_pdv_formas_pagamento (
    id_forma_pagamento BIGINT PRIMARY KEY REFERENCES tb_formas_pagamento(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tb_configuracoes_pdv_condicoes_pagamento (
    id_condicao_pagamento BIGINT PRIMARY KEY REFERENCES tb_condicoes_pagamento(id) ON DELETE CASCADE
);

ALTER TABLE tb_vendas_pdv
    ADD COLUMN IF NOT EXISTS id_condicao_pagamento BIGINT REFERENCES tb_condicoes_pagamento(id),
    ADD COLUMN IF NOT EXISTS valor_desconto DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (valor_desconto >= 0);

ALTER TABLE tb_movimento_financeiro
    ADD COLUMN IF NOT EXISTS id_venda_pdv BIGINT REFERENCES tb_vendas_pdv(id),
    ADD COLUMN IF NOT EXISTS id_categoria_credito BIGINT REFERENCES tb_categorias_contas_receber(id);

ALTER TABLE tb_movimento_financeiro DROP CONSTRAINT IF EXISTS chk_movimento_financeiro_origem;
ALTER TABLE tb_movimento_financeiro DROP CONSTRAINT IF EXISTS tb_movimento_financeiro_tipo_movimento_check;
ALTER TABLE tb_movimento_financeiro ADD CONSTRAINT tb_movimento_financeiro_tipo_movimento_check CHECK (tipo_movimento IN ('CONTA_RECEBER', 'CONTA_PAGAR', 'VENDA_PDV'));
ALTER TABLE tb_movimento_financeiro ADD CONSTRAINT chk_movimento_financeiro_origem CHECK (
    (tipo_movimento = 'CONTA_PAGAR' AND id_conta_pagar IS NOT NULL AND id_conta_receber IS NULL AND id_venda_pdv IS NULL)
    OR (tipo_movimento = 'CONTA_RECEBER' AND id_conta_receber IS NOT NULL AND id_conta_pagar IS NULL AND id_venda_pdv IS NULL)
    OR (tipo_movimento = 'VENDA_PDV' AND id_venda_pdv IS NOT NULL AND id_conta_pagar IS NULL AND id_conta_receber IS NULL AND id_categoria_credito IS NOT NULL)
);

ALTER TABLE tb_controle_caixa ADD COLUMN IF NOT EXISTS justificativa_diferenca TEXT;

COMMIT;
