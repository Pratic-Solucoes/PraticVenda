BEGIN;

ALTER TABLE tb_produtos DROP CONSTRAINT IF EXISTS ck_produtos_fornecedor_obrigatorio;
ALTER TABLE tb_produtos ADD CONSTRAINT ck_produtos_fornecedor_obrigatorio
    CHECK (composto OR id_fornecedor IS NOT NULL);

COMMIT;
