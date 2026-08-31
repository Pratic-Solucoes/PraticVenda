BEGIN;
ALTER TABLE tb_formas_pagamento DROP CONSTRAINT IF EXISTS tb_formas_pagamento_tipo_check;
ALTER TABLE tb_formas_pagamento DROP CONSTRAINT IF EXISTS ck_formas_pagamento_tipo;
ALTER TABLE tb_formas_pagamento ADD CONSTRAINT ck_formas_pagamento_tipo CHECK (tipo IN ('DINHEIRO','CARTAO','PIX'));
ALTER TABLE tb_movimento_caixa DROP CONSTRAINT IF EXISTS tb_movimento_caixa_tipo_credito_check;
ALTER TABLE tb_movimento_caixa DROP CONSTRAINT IF EXISTS ck_movimento_caixa_tipo_credito;
ALTER TABLE tb_movimento_caixa ADD CONSTRAINT ck_movimento_caixa_tipo_credito CHECK (tipo_credito IN ('DINHEIRO','CARTAO','PIX'));
COMMIT;
