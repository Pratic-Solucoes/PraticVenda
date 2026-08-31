BEGIN;

ALTER TABLE tb_contas_receber ALTER COLUMN id_cliente DROP NOT NULL;

COMMIT;
