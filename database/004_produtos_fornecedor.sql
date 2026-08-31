BEGIN;

-- Para bases já existentes, cadastre/associe o fornecedor de cada produto pela
-- tela antes de tornar a coluna obrigatória diretamente no banco.
ALTER TABLE tb_produtos ADD COLUMN IF NOT EXISTS id_fornecedor BIGINT REFERENCES tb_fornecedores(id);
DO $$ BEGIN
    ALTER TABLE tb_produtos ADD CONSTRAINT ck_produtos_fornecedor_obrigatorio CHECK (id_fornecedor IS NOT NULL) NOT VALID;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
CREATE INDEX IF NOT EXISTS idx_produtos_fornecedor ON tb_produtos(id_fornecedor);

COMMIT;
