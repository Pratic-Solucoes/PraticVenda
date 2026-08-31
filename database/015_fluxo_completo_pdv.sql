BEGIN;

ALTER TABLE tb_configuracoes_pdv
    ADD COLUMN IF NOT EXISTS gerar_financeiro_recebido BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE tb_vendas_pdv
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'CONCLUIDA',
    ADD COLUMN IF NOT EXISTS apelido_consumidor VARCHAR(100),
    ADD COLUMN IF NOT EXISTS concluido_em TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS cancelado_em TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS id_usuario_cancelamento BIGINT REFERENCES public.tb_usuarios_gestao(id);

ALTER TABLE tb_vendas_pdv DROP CONSTRAINT IF EXISTS ck_vendas_pdv_status;
ALTER TABLE tb_vendas_pdv ADD CONSTRAINT ck_vendas_pdv_status CHECK (status IN ('ABERTA','CONCLUIDA','CANCELADA'));
UPDATE tb_vendas_pdv SET concluido_em = criado_em WHERE status = 'CONCLUIDA' AND concluido_em IS NULL;

ALTER TABLE tb_movimento_caixa DROP CONSTRAINT IF EXISTS tb_movimento_caixa_tipo_movimento_check;
ALTER TABLE tb_movimento_caixa ADD CONSTRAINT tb_movimento_caixa_tipo_movimento_check CHECK (tipo_movimento IN ('ABERTURA','VENDA','ESTORNO'));

ALTER TABLE tb_movimento_financeiro DROP CONSTRAINT IF EXISTS chk_movimento_financeiro_origem;
ALTER TABLE tb_movimento_financeiro DROP CONSTRAINT IF EXISTS tb_movimento_financeiro_tipo_movimento_check;
ALTER TABLE tb_movimento_financeiro ADD CONSTRAINT tb_movimento_financeiro_tipo_movimento_check CHECK (tipo_movimento IN ('CONTA_RECEBER','CONTA_PAGAR','VENDA_PDV','ESTORNO_VENDA_PDV'));
ALTER TABLE tb_movimento_financeiro ADD CONSTRAINT chk_movimento_financeiro_origem CHECK (
    (tipo_movimento = 'CONTA_PAGAR' AND id_conta_pagar IS NOT NULL AND id_conta_receber IS NULL AND id_venda_pdv IS NULL)
    OR (tipo_movimento = 'CONTA_RECEBER' AND id_conta_receber IS NOT NULL AND id_conta_pagar IS NULL AND id_venda_pdv IS NULL)
    OR (tipo_movimento IN ('VENDA_PDV','ESTORNO_VENDA_PDV') AND id_venda_pdv IS NOT NULL AND id_conta_pagar IS NULL AND id_conta_receber IS NULL AND id_categoria_credito IS NOT NULL)
);

INSERT INTO tb_categoria_movimento_estoque(nome) VALUES ('CANCELAMENTO VENDA PDV') ON CONFLICT(nome) DO NOTHING;
CREATE INDEX IF NOT EXISTS idx_vendas_pdv_status_criado_em ON tb_vendas_pdv(status, criado_em DESC);

COMMIT;
