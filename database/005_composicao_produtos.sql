BEGIN;

CREATE TABLE IF NOT EXISTS tb_composicoes_produtos (
    id BIGSERIAL PRIMARY KEY,
    id_produto_composto BIGINT NOT NULL REFERENCES tb_produtos(id) ON DELETE CASCADE,
    id_produto_componente BIGINT NOT NULL REFERENCES tb_produtos(id) ON DELETE RESTRICT,
    quantidade DECIMAL(12,3) NOT NULL CHECK (quantidade > 0),
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ck_composicao_produto_distinto CHECK (id_produto_composto <> id_produto_componente),
    CONSTRAINT uq_composicao_produto_componente UNIQUE (id_produto_composto, id_produto_componente)
);

CREATE INDEX IF NOT EXISTS idx_composicoes_produto_composto ON tb_composicoes_produtos(id_produto_composto);
COMMIT;
