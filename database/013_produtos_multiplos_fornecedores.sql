BEGIN;
CREATE TABLE IF NOT EXISTS tb_produtos_fornecedores (
    id_produto BIGINT NOT NULL REFERENCES tb_produtos(id) ON DELETE CASCADE,
    id_fornecedor BIGINT NOT NULL REFERENCES tb_fornecedores(id) ON DELETE RESTRICT,
    PRIMARY KEY (id_produto, id_fornecedor)
);
INSERT INTO tb_produtos_fornecedores(id_produto,id_fornecedor)
SELECT id,id_fornecedor FROM tb_produtos WHERE id_fornecedor IS NOT NULL ON CONFLICT DO NOTHING;
CREATE INDEX IF NOT EXISTS idx_produtos_fornecedores_fornecedor ON tb_produtos_fornecedores(id_fornecedor);
COMMIT;
