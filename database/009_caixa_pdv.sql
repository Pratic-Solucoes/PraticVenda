BEGIN;

ALTER TABLE tb_formas_pagamento ADD COLUMN IF NOT EXISTS tipo VARCHAR(20) NOT NULL DEFAULT 'CARTAO' CONSTRAINT ck_formas_pagamento_tipo CHECK (tipo IN ('DINHEIRO','CARTAO','PIX'));
UPDATE tb_formas_pagamento SET tipo = 'DINHEIRO' WHERE LOWER(descricao) LIKE '%dinheiro%';

CREATE TABLE IF NOT EXISTS tb_caixas (
    id BIGSERIAL PRIMARY KEY,
    id_usuario BIGINT NOT NULL REFERENCES public.tb_usuarios_gestao(id),
    nome VARCHAR(100) NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (id_usuario, nome)
);

CREATE TABLE IF NOT EXISTS tb_controle_caixa (
    id BIGSERIAL PRIMARY KEY,
    id_caixa BIGINT NOT NULL REFERENCES tb_caixas(id),
    id_usuario BIGINT NOT NULL REFERENCES public.tb_usuarios_gestao(id),
    status VARCHAR(10) NOT NULL DEFAULT 'ABERTO' CHECK (status IN ('ABERTO','FECHADO')),
    aberto_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valor_abertura DECIMAL(15,2) NOT NULL CHECK (valor_abertura >= 0),
    fechado_em TIMESTAMPTZ,
    valor_dinheiro_informado DECIMAL(15,2),
    valor_cartao_informado DECIMAL(15,2),
    valor_dinheiro_esperado DECIMAL(15,2),
    valor_cartao_esperado DECIMAL(15,2),
    diferenca_dinheiro DECIMAL(15,2),
    diferenca_cartao DECIMAL(15,2)
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_controle_caixa_aberto ON tb_controle_caixa (id_caixa) WHERE status = 'ABERTO';

CREATE TABLE IF NOT EXISTS tb_movimento_caixa (
    id BIGSERIAL PRIMARY KEY,
    id_controle_caixa BIGINT NOT NULL REFERENCES tb_controle_caixa(id),
    tipo_movimento VARCHAR(20) NOT NULL CHECK (tipo_movimento IN ('ABERTURA','VENDA')),
    id_venda BIGINT,
    id_cliente BIGINT REFERENCES tb_clientes(id),
    tipo_credito VARCHAR(20) NOT NULL CONSTRAINT ck_movimento_caixa_tipo_credito CHECK (tipo_credito IN ('DINHEIRO','CARTAO','PIX')),
    valor_credito DECIMAL(15,2) NOT NULL CHECK (valor_credito >= 0),
    criado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE tb_vendas_pdv
    ADD COLUMN IF NOT EXISTS id_controle_caixa BIGINT REFERENCES tb_controle_caixa(id),
    ADD COLUMN IF NOT EXISTS id_cliente BIGINT REFERENCES tb_clientes(id),
    ADD COLUMN IF NOT EXISTS id_forma_pagamento BIGINT REFERENCES tb_formas_pagamento(id);

COMMIT;
