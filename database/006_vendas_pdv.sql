BEGIN;
CREATE TABLE IF NOT EXISTS tb_vendas_pdv (
 id BIGSERIAL PRIMARY KEY, id_estoque BIGINT NOT NULL REFERENCES tb_estoques(id), id_usuario BIGINT NOT NULL REFERENCES public.tb_usuarios_gestao(id), valor_total DECIMAL(12,2) NOT NULL, criado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS tb_itens_vendas_pdv (
 id BIGSERIAL PRIMARY KEY, id_venda BIGINT NOT NULL REFERENCES tb_vendas_pdv(id) ON DELETE CASCADE, id_produto BIGINT NOT NULL REFERENCES tb_produtos(id), quantidade DECIMAL(12,3) NOT NULL CHECK(quantidade>0), valor_unitario DECIMAL(12,2) NOT NULL, valor_total DECIMAL(12,2) NOT NULL
);
INSERT INTO tb_categoria_movimento_estoque(nome) VALUES ('VENDA PDV') ON CONFLICT(nome) DO NOTHING;
COMMIT;
