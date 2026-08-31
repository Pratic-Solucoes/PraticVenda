BEGIN;

ALTER TABLE tb_produtos ADD COLUMN IF NOT EXISTS materia_prima BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE tb_produtos DROP CONSTRAINT IF EXISTS ck_produtos_tipo_operacional;
ALTER TABLE tb_produtos ADD CONSTRAINT ck_produtos_tipo_operacional CHECK (NOT (composto AND materia_prima));

COMMIT;
