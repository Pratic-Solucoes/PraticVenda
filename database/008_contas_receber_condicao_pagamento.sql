BEGIN;

-- Complementa instalações existentes para que cada crédito registre a condição
-- e a forma de pagamento escolhidas no momento do lançamento.
ALTER TABLE tb_contas_receber
    ADD COLUMN IF NOT EXISTS id_condicao_pagamento BIGINT REFERENCES tb_condicoes_pagamento(id),
    ADD COLUMN IF NOT EXISTS id_forma_pagamento BIGINT REFERENCES tb_formas_pagamento(id);

CREATE INDEX IF NOT EXISTS idx_contas_receber_vencimento
    ON tb_contas_receber (dt_vencimento, status);

CREATE INDEX IF NOT EXISTS idx_contas_receber_cliente
    ON tb_contas_receber (id_cliente);

COMMIT;
