create table if not exists tb_categorias_contas_pagar(
    id bigserial primary key,
    descricao varchar(255) not null unique
);

create table if not exists tb_formas_pagamento(
    id bigserial primary key,
    descricao varchar(100) not null unique
)

create table if not exists tb_condicoes_pagamento(
    id bigserial primary key,
    descricao varchar(100) not null,
    qtd_parcelas bigint not null,
    dias_primeiro_venc bigint not null,
    intervalo_parcelas bigint not null
);

create table if not exists tb_condicao_forma_pagamento(
    id_condicao bigint not null,
    id_forma_pagamento bigint not null,
    primary key (id_condicao, id_forma_pagamento),
    foreign key (id_condicao) references tb_condicoes_pagamento(id) on delete cascade,
    foreign key (id_forma_pagamento) references tb_formas_pagamento(id) on delete cascade
);

create table if not exists tb_contas_receber(
    id bigserial primary key,
    id_cliente bigint not null,
    id_categoria bigint not null,
    tipo_origem varchar(50) not null,
    id_origem bigint,
    id_grupo_parcelas UUID,
    descricao varchar(255),
    valor_original DECIMAL(15,2) NOT NULL,
    saldo_restante DECIMAL(15,2) NOT NULL,
    dt_vencimento DATE NOT NULL,
    nr_parcela INT NOT NULL,
    nr_total_parcelas INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE' CHECK (status IN ('PENDENTE', 'PAGO_PARCIAL', 'PAGO', 'CANCELADO')),
    dt_emissao TIMESTAMP NOT NULL,
    dt_pagamento TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_cliente) REFERENCES tb_clientes(id),
    FOREIGN KEY (id_categoria) REFERENCES tb_categorias_contas_receber(id)
)

CREATE TABLE IF NOT EXISTS tb_contas_pagar (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL,
    id_categoria BIGINT NOT NULL,
    id_grupo_parcelas UUID,
    descricao VARCHAR(255),
    nr_documento VARCHAR(255),
    nr_nota_fiscal VARCHAR(255),
    valor_original DECIMAL(15,2) NOT NULL,
    saldo_restante DECIMAL(15,2) NOT NULL,
    dt_entrada DATE NOT NULL,
    dt_vencimento DATE NOT NULL,
    nr_parcela INT NOT NULL,
    nr_total_parcelas INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE' CHECK (status IN ('PENDENTE', 'PAGO_PARCIAL', 'PAGO', 'CANCELADO')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    dt_pagamento TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_fornecedor) REFERENCES tb_fornecedores(id),
    FOREIGN KEY (id_categoria) REFERENCES tb_categorias_contas_pagar(id)
);

CREATE TABLE IF NOT EXISTS tb_movimento_financeiro (
    id BIGSERIAL PRIMARY KEY,
    tipo_movimento VARCHAR(20) NOT NULL CHECK (tipo_movimento IN ('CONTA_RECEBER', 'CONTA_PAGAR')),
    id_conta_pagar BIGINT,
    id_conta_receber BIGINT,
    dt_movimento DATE NOT NULL,
    valor_movimento DECIMAL(15,2) NOT NULL,
    valor_acrescimo DECIMAL(15,2) DEFAULT 0.00,
    valor_desconto DECIMAL(15,2) DEFAULT 0.00,
    forma_pagamento VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_conta_pagar) REFERENCES tb_contas_pagar(id)
    -- FOREIGN KEY (id_conta_receber) REFERENCES tb_contas_receber(id) -- Descomentar quando existir
);
