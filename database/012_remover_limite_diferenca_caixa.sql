BEGIN;

ALTER TABLE tb_configuracoes_pdv
    DROP COLUMN IF EXISTS limite_diferenca_caixa,
    DROP COLUMN IF EXISTS exigir_justificativa_diferenca;

ALTER TABLE tb_controle_caixa
    DROP COLUMN IF EXISTS justificativa_diferenca;

COMMIT;
