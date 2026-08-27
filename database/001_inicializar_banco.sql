-- PraticVenda: estrutura única e idempotente do banco PostgreSQL.
-- Execute em um banco novo com:
--   psql "$DATABASE_URL" -f database/001_inicializar_banco.sql
--
-- Estrutura:
--   * public: empresas e usuários usados na autenticação.
--   * schemas schema_pao_quente e schema_mercadinho: dados operacionais isolados.
-- Não há funções, procedures ou SQL dinâmico neste arquivo.

BEGIN;

CREATE TABLE IF NOT EXISTS public.tb_empresas_gestao (
    id BIGSERIAL PRIMARY KEY,
    nome_fantasia VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    telefone VARCHAR(20) NOT NULL UNIQUE,
    schema VARCHAR(50) NOT NULL UNIQUE,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_empresas_gestao_schema
        CHECK (schema ~ '^[a-z][a-z0-9_]{0,48}$')
);

CREATE TABLE IF NOT EXISTS public.tb_usuarios_admin(
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    celular VARCHAR(20) NOT NULL UNIQUE,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.tb_usuarios_gestao (
    id BIGSERIAL PRIMARY KEY,
    id_empresa BIGINT NOT NULL DEFAULT 1 REFERENCES public.tb_empresas_gestao(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) UNIQUE,
    telefone VARCHAR(20) UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ativo BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_usuarios_gestao_empresa ON public.tb_usuarios_gestao(id_empresa);

-- Schemas operacionais dos tenants de desenvolvimento.
CREATE SCHEMA IF NOT EXISTS schema_pao_quente;
CREATE SCHEMA IF NOT EXISTS schema_mercadinho;

-- Estrutura operacional: Padaria Pão Quente.
SET search_path TO schema_pao_quente, public;
CREATE TABLE IF NOT EXISTS tb_empresas (
    id BIGSERIAL PRIMARY KEY,
    razao_social VARCHAR(100) NOT NULL,
    nome_fantasia VARCHAR(100) NOT NULL,
    cnpj VARCHAR(20) NOT NULL UNIQUE,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    data_criacao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_config_fiscais_empresa (
    id BIGSERIAL PRIMARY KEY,
    id_empresa BIGINT NOT NULL UNIQUE REFERENCES tb_empresas(id) ON DELETE CASCADE,
    cd_regime_tributario INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tb_credenciais_empresa (
    id BIGSERIAL PRIMARY KEY,
    id_empresa BIGINT NOT NULL REFERENCES tb_empresas(id) ON DELETE CASCADE,
    tp_ambiente SMALLINT NOT NULL CHECK (tp_ambiente IN (1, 2)),
    certificado_digital BYTEA,
    senha_criptografada VARCHAR(255),
    id_csc VARCHAR(6),
    csc_nfe VARCHAR(255),
    CONSTRAINT uq_credenciais_empresa_ambiente UNIQUE (id_empresa, tp_ambiente)
);

CREATE TABLE IF NOT EXISTS tb_endereco_empresa (
    id BIGSERIAL PRIMARY KEY,
    id_empresa BIGINT NOT NULL REFERENCES tb_empresas(id) ON DELETE CASCADE,
    logradouro VARCHAR(100) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    bairro VARCHAR(100) NOT NULL,
    cep VARCHAR(9) NOT NULL,
    cd_cidade BIGINT NOT NULL,
    nome_cidade VARCHAR(100) NOT NULL,
    estado VARCHAR(2) NOT NULL,
    nome_pais VARCHAR(50) NOT NULL DEFAULT 'Brasil',
    cd_pais BIGINT NOT NULL DEFAULT 1058,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    data_criacao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_clientes (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(150) NOT NULL UNIQUE,
    tipo VARCHAR(50) NOT NULL,
    email VARCHAR(200) UNIQUE,
    telefone VARCHAR(20) UNIQUE,
    cpf VARCHAR(14) UNIQUE,
    cnpj VARCHAR(18) UNIQUE,
    contribuinte VARCHAR(50),
    is_consumidor_final BOOLEAN NOT NULL DEFAULT FALSE,
    ie VARCHAR(14) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_enderecos_clientes (
    id BIGSERIAL PRIMARY KEY,
    id_cliente BIGINT NOT NULL REFERENCES tb_clientes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    cep VARCHAR(9) NOT NULL,
    logradouro VARCHAR(255) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    bairro VARCHAR(100) NOT NULL,
    municipio VARCHAR(100) NOT NULL,
    uf VARCHAR(2) NOT NULL,
    codigo_municipio VARCHAR(7) NOT NULL,
    is_principal BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_telefones_clientes (
    id BIGSERIAL PRIMARY KEY,
    id_cliente BIGINT NOT NULL REFERENCES tb_clientes(id) ON DELETE CASCADE,
    ddd CHAR(2) NOT NULL,
    numero VARCHAR(15) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    razao_social VARCHAR(255) NOT NULL,
    cnpj VARCHAR(14) NOT NULL UNIQUE,
    inscricao_estadual VARCHAR(20),
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_enderecos_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL REFERENCES tb_fornecedores(id) ON DELETE CASCADE,
    cep VARCHAR(9) NOT NULL,
    logradouro VARCHAR(255) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    bairro VARCHAR(100) NOT NULL,
    municipio VARCHAR(100) NOT NULL,
    uf CHAR(2) NOT NULL,
    codigo_municipio VARCHAR(7) NOT NULL,
    is_principal BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_telefones_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL REFERENCES tb_fornecedores(id) ON DELETE CASCADE,
    ddd CHAR(2) NOT NULL,
    numero VARCHAR(15) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_categorias_contas_pagar (
    id BIGSERIAL PRIMARY KEY,
    descricao VARCHAR(255) NOT NULL UNIQUE,
    nome VARCHAR(255) GENERATED ALWAYS AS (descricao) STORED
);

CREATE TABLE IF NOT EXISTS tb_categorias_contas_receber (
    id BIGSERIAL PRIMARY KEY,
    descricao VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tb_formas_pagamento (
    id BIGSERIAL PRIMARY KEY,
    descricao VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tb_condicoes_pagamento (
    id BIGSERIAL PRIMARY KEY,
    descricao VARCHAR(100) NOT NULL,
    qtd_parcelas BIGINT NOT NULL CHECK (qtd_parcelas > 0),
    dias_primeiro_venc BIGINT NOT NULL CHECK (dias_primeiro_venc >= 0),
    intervalo_parcelas BIGINT NOT NULL CHECK (intervalo_parcelas >= 0)
);

CREATE TABLE IF NOT EXISTS tb_condicao_forma_pagamento (
    id_condicao BIGINT NOT NULL REFERENCES tb_condicoes_pagamento(id) ON DELETE CASCADE,
    id_forma_pagamento BIGINT NOT NULL REFERENCES tb_formas_pagamento(id) ON DELETE CASCADE,
    PRIMARY KEY (id_condicao, id_forma_pagamento)
);

CREATE TABLE IF NOT EXISTS tb_contas_pagar (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL REFERENCES tb_fornecedores(id),
    id_categoria BIGINT NOT NULL REFERENCES tb_categorias_contas_pagar(id),
    id_grupo_parcelas UUID,
    descricao VARCHAR(255),
    nr_documento VARCHAR(255),
    nr_nota_fiscal VARCHAR(255),
    valor_original DECIMAL(15,2) NOT NULL CHECK (valor_original > 0),
    saldo_restante DECIMAL(15,2) NOT NULL CHECK (saldo_restante >= 0),
    valor DECIMAL(15,2) GENERATED ALWAYS AS (valor_original) STORED,
    dt_entrada DATE NOT NULL,
    dt_vencimento DATE NOT NULL,
    nr_parcela INTEGER NOT NULL CHECK (nr_parcela > 0),
    nr_total_parcelas INTEGER NOT NULL CHECK (nr_total_parcelas > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE'
        CHECK (status IN ('PENDENTE', 'PAGO_PARCIAL', 'PAGO', 'CANCELADO')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    dt_pagamento TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_contas_pagar_vencimento
    ON tb_contas_pagar (dt_vencimento, status);

CREATE TABLE IF NOT EXISTS tb_contas_receber (
    id BIGSERIAL PRIMARY KEY,
    id_cliente BIGINT NOT NULL REFERENCES tb_clientes(id),
    id_categoria BIGINT NOT NULL REFERENCES tb_categorias_contas_receber(id),
    tipo_origem VARCHAR(50) NOT NULL,
    id_origem BIGINT,
    id_grupo_parcelas UUID,
    descricao VARCHAR(255),
    valor_original DECIMAL(15,2) NOT NULL CHECK (valor_original > 0),
    saldo_restante DECIMAL(15,2) NOT NULL CHECK (saldo_restante >= 0),
    dt_vencimento DATE NOT NULL,
    nr_parcela INTEGER NOT NULL CHECK (nr_parcela > 0),
    nr_total_parcelas INTEGER NOT NULL CHECK (nr_total_parcelas > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE'
        CHECK (status IN ('PENDENTE', 'PAGO_PARCIAL', 'PAGO', 'CANCELADO')),
    dt_emissao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    dt_pagamento TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_movimento_financeiro (
    id BIGSERIAL PRIMARY KEY,
    tipo_movimento VARCHAR(20) NOT NULL
        CHECK (tipo_movimento IN ('CONTA_RECEBER', 'CONTA_PAGAR')),
    id_conta_pagar BIGINT REFERENCES tb_contas_pagar(id),
    id_conta_receber BIGINT REFERENCES tb_contas_receber(id),
    dt_movimento DATE NOT NULL,
    valor_movimento DECIMAL(15,2) NOT NULL CHECK (valor_movimento > 0),
    valor_acrescimo DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (valor_acrescimo >= 0),
    valor_desconto DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (valor_desconto >= 0),
    forma_pagamento VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_movimento_financeiro_origem CHECK (
        (tipo_movimento = 'CONTA_PAGAR' AND id_conta_pagar IS NOT NULL AND id_conta_receber IS NULL)
        OR (tipo_movimento = 'CONTA_RECEBER' AND id_conta_receber IS NOT NULL AND id_conta_pagar IS NULL)
    )
);

CREATE TABLE IF NOT EXISTS tb_grupos_tributarios (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    cfop_padrao VARCHAR(4) NOT NULL,
    origem_mercadoria INTEGER NOT NULL DEFAULT 0,
    csosn VARCHAR(4),
    icms_cst VARCHAR(3),
    icms_aliquota DECIMAL(5,2) NOT NULL DEFAULT 0,
    icms_mva_st DECIMAL(5,2) NOT NULL DEFAULT 0,
    icms_aliquota_st DECIMAL(5,2) NOT NULL DEFAULT 0,
    ipi_cst VARCHAR(2),
    ipi_aliquota DECIMAL(5,2) NOT NULL DEFAULT 0,
    pis_cst VARCHAR(2) NOT NULL DEFAULT '07',
    pis_aliquota DECIMAL(5,2) NOT NULL DEFAULT 0,
    cofins_cst VARCHAR(2) NOT NULL DEFAULT '07',
    cofins_aliquota DECIMAL(5,2) NOT NULL DEFAULT 0,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_produtos (
    id BIGSERIAL PRIMARY KEY,
    codigo_barras VARCHAR(14) UNIQUE,
    codigo_interno_loja VARCHAR(50),
    nome VARCHAR(150) NOT NULL,
    descricao TEXT,
    preco_custo DECIMAL(12,4) NOT NULL DEFAULT 0,
    preco_venda DECIMAL(12,2) NOT NULL DEFAULT 0,
    unidade_estoque VARCHAR(10) NOT NULL DEFAULT 'UN',
    unidade_venda VARCHAR(10) NOT NULL DEFAULT 'UN',
    peso_bruto DECIMAL(12,2) NOT NULL DEFAULT 0,
    peso_liquido DECIMAL(12,2) NOT NULL DEFAULT 0,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_produtos_codigo_interno ON tb_produtos(codigo_interno_loja);
CREATE INDEX IF NOT EXISTS idx_produtos_nome_ativos ON tb_produtos(nome) WHERE ativo = TRUE;

CREATE TABLE IF NOT EXISTS tb_produtos_fiscal (
    id_produto BIGINT PRIMARY KEY REFERENCES tb_produtos(id) ON DELETE CASCADE,
    ncm VARCHAR(8) NOT NULL,
    cest VARCHAR(7),
    id_grupo_tributario BIGINT NOT NULL REFERENCES tb_grupos_tributarios(id) ON DELETE RESTRICT,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_estoques (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    descricao VARCHAR(255),
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_produtos_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_produto BIGINT NOT NULL REFERENCES tb_produtos(id) ON DELETE CASCADE,
    id_estoque BIGINT NOT NULL REFERENCES tb_estoques(id) ON DELETE RESTRICT,
    quantidade DECIMAL(10,3) NOT NULL DEFAULT 0,
    estoque_minimo DECIMAL(10,3) NOT NULL DEFAULT 0,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_produto_por_estoque UNIQUE (id_produto, id_estoque)
);

CREATE INDEX IF NOT EXISTS idx_produtos_estoque_produto ON tb_produtos_estoque(id_produto);

CREATE TABLE IF NOT EXISTS tb_categoria_movimento_estoque (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tb_movimento_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_produto BIGINT NOT NULL REFERENCES tb_produtos(id),
    id_estoque BIGINT NOT NULL REFERENCES tb_estoques(id),
    id_usuario BIGINT NOT NULL REFERENCES public.tb_usuarios_gestao(id),
    quantidade DECIMAL(10,3) NOT NULL CHECK (quantidade > 0),
    tipo_movimento VARCHAR(10) NOT NULL CHECK (tipo_movimento IN ('ENTRADA', 'SAIDA')),
    id_categoria_movimento BIGINT NOT NULL REFERENCES tb_categoria_movimento_estoque(id),
    id_origem VARCHAR(50),
    observacao TEXT,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_movimento_estoque_produto ON tb_movimento_estoque(id_produto, id_estoque);
CREATE INDEX IF NOT EXISTS idx_movimento_estoque_criado_em ON tb_movimento_estoque(criado_em);

CREATE TABLE IF NOT EXISTS tb_entradas_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_estoque BIGINT NOT NULL REFERENCES tb_estoques(id),
    id_fornecedor BIGINT NOT NULL REFERENCES tb_fornecedores(id),
    valor_despesa_adicional DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_despesa_adicional >= 0),
    id_usuario BIGINT NOT NULL REFERENCES public.tb_usuarios_gestao(id),
    valor_total DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_total >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'ABERTO'
        CHECK (status IN ('ABERTO', 'CONCLUIDA', 'CANCELADA')),
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_produtos_entradas_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_entrada_estoque BIGINT NOT NULL REFERENCES tb_entradas_estoque(id) ON DELETE CASCADE,
    id_produto BIGINT NOT NULL REFERENCES tb_produtos(id),
    valor_unitario DECIMAL(10,2) NOT NULL CHECK (valor_unitario > 0),
    valor_icms_st DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_icms_st >= 0),
    valor_ipi DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_ipi >= 0),
    valor_desconto DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_desconto >= 0),
    rateio_despesa_adicional DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (rateio_despesa_adicional >= 0),
    valor_custo DECIMAL(10,2) NOT NULL CHECK (valor_custo > 0),
    valor_total DECIMAL(10,2) NOT NULL CHECK (valor_total > 0),
    quantidade DECIMAL(10,3) NOT NULL CHECK (quantidade > 0)
);

-- Seeds necessários para a entrada de estoque.
INSERT INTO tb_categoria_movimento_estoque (nome)
VALUES ('ENTRADA DE ESTOQUE')
ON CONFLICT (nome) DO NOTHING;

INSERT INTO tb_estoques (nome, descricao)
VALUES ('Estoque Geral', 'Local padrão de armazenamento e vendas')
ON CONFLICT (nome) DO NOTHING;

-- Estrutura operacional: Mercadinho da Esquina.
SET search_path TO schema_mercadinho, public;
CREATE TABLE IF NOT EXISTS tb_empresas (
    id BIGSERIAL PRIMARY KEY,
    razao_social VARCHAR(100) NOT NULL,
    nome_fantasia VARCHAR(100) NOT NULL,
    cnpj VARCHAR(20) NOT NULL UNIQUE,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    data_criacao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_config_fiscais_empresa (
    id BIGSERIAL PRIMARY KEY,
    id_empresa BIGINT NOT NULL UNIQUE REFERENCES tb_empresas(id) ON DELETE CASCADE,
    cd_regime_tributario INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tb_credenciais_empresa (
    id BIGSERIAL PRIMARY KEY,
    id_empresa BIGINT NOT NULL REFERENCES tb_empresas(id) ON DELETE CASCADE,
    tp_ambiente SMALLINT NOT NULL CHECK (tp_ambiente IN (1, 2)),
    certificado_digital BYTEA,
    senha_criptografada VARCHAR(255),
    id_csc VARCHAR(6),
    csc_nfe VARCHAR(255),
    CONSTRAINT uq_credenciais_empresa_ambiente UNIQUE (id_empresa, tp_ambiente)
);

CREATE TABLE IF NOT EXISTS tb_endereco_empresa (
    id BIGSERIAL PRIMARY KEY,
    id_empresa BIGINT NOT NULL REFERENCES tb_empresas(id) ON DELETE CASCADE,
    logradouro VARCHAR(100) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    bairro VARCHAR(100) NOT NULL,
    cep VARCHAR(9) NOT NULL,
    cd_cidade BIGINT NOT NULL,
    nome_cidade VARCHAR(100) NOT NULL,
    estado VARCHAR(2) NOT NULL,
    nome_pais VARCHAR(50) NOT NULL DEFAULT 'Brasil',
    cd_pais BIGINT NOT NULL DEFAULT 1058,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    data_criacao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_clientes (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(150) NOT NULL UNIQUE,
    tipo VARCHAR(50) NOT NULL,
    email VARCHAR(200) UNIQUE,
    telefone VARCHAR(20) UNIQUE,
    cpf VARCHAR(14) UNIQUE,
    cnpj VARCHAR(18) UNIQUE,
    contribuinte VARCHAR(50),
    is_consumidor_final BOOLEAN NOT NULL DEFAULT FALSE,
    ie VARCHAR(14) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_enderecos_clientes (
    id BIGSERIAL PRIMARY KEY,
    id_cliente BIGINT NOT NULL REFERENCES tb_clientes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    cep VARCHAR(9) NOT NULL,
    logradouro VARCHAR(255) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    bairro VARCHAR(100) NOT NULL,
    municipio VARCHAR(100) NOT NULL,
    uf VARCHAR(2) NOT NULL,
    codigo_municipio VARCHAR(7) NOT NULL,
    is_principal BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_telefones_clientes (
    id BIGSERIAL PRIMARY KEY,
    id_cliente BIGINT NOT NULL REFERENCES tb_clientes(id) ON DELETE CASCADE,
    ddd CHAR(2) NOT NULL,
    numero VARCHAR(15) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    razao_social VARCHAR(255) NOT NULL,
    cnpj VARCHAR(14) NOT NULL UNIQUE,
    inscricao_estadual VARCHAR(20),
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_enderecos_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL REFERENCES tb_fornecedores(id) ON DELETE CASCADE,
    cep VARCHAR(9) NOT NULL,
    logradouro VARCHAR(255) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    bairro VARCHAR(100) NOT NULL,
    municipio VARCHAR(100) NOT NULL,
    uf CHAR(2) NOT NULL,
    codigo_municipio VARCHAR(7) NOT NULL,
    is_principal BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_telefones_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL REFERENCES tb_fornecedores(id) ON DELETE CASCADE,
    ddd CHAR(2) NOT NULL,
    numero VARCHAR(15) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_categorias_contas_pagar (
    id BIGSERIAL PRIMARY KEY,
    descricao VARCHAR(255) NOT NULL UNIQUE,
    nome VARCHAR(255) GENERATED ALWAYS AS (descricao) STORED
);

CREATE TABLE IF NOT EXISTS tb_categorias_contas_receber (
    id BIGSERIAL PRIMARY KEY,
    descricao VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tb_formas_pagamento (
    id BIGSERIAL PRIMARY KEY,
    descricao VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tb_condicoes_pagamento (
    id BIGSERIAL PRIMARY KEY,
    descricao VARCHAR(100) NOT NULL,
    qtd_parcelas BIGINT NOT NULL CHECK (qtd_parcelas > 0),
    dias_primeiro_venc BIGINT NOT NULL CHECK (dias_primeiro_venc >= 0),
    intervalo_parcelas BIGINT NOT NULL CHECK (intervalo_parcelas >= 0)
);

CREATE TABLE IF NOT EXISTS tb_condicao_forma_pagamento (
    id_condicao BIGINT NOT NULL REFERENCES tb_condicoes_pagamento(id) ON DELETE CASCADE,
    id_forma_pagamento BIGINT NOT NULL REFERENCES tb_formas_pagamento(id) ON DELETE CASCADE,
    PRIMARY KEY (id_condicao, id_forma_pagamento)
);

CREATE TABLE IF NOT EXISTS tb_contas_pagar (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL REFERENCES tb_fornecedores(id),
    id_categoria BIGINT NOT NULL REFERENCES tb_categorias_contas_pagar(id),
    id_grupo_parcelas UUID,
    descricao VARCHAR(255),
    nr_documento VARCHAR(255),
    nr_nota_fiscal VARCHAR(255),
    valor_original DECIMAL(15,2) NOT NULL CHECK (valor_original > 0),
    saldo_restante DECIMAL(15,2) NOT NULL CHECK (saldo_restante >= 0),
    valor DECIMAL(15,2) GENERATED ALWAYS AS (valor_original) STORED,
    dt_entrada DATE NOT NULL,
    dt_vencimento DATE NOT NULL,
    nr_parcela INTEGER NOT NULL CHECK (nr_parcela > 0),
    nr_total_parcelas INTEGER NOT NULL CHECK (nr_total_parcelas > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE'
        CHECK (status IN ('PENDENTE', 'PAGO_PARCIAL', 'PAGO', 'CANCELADO')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    dt_pagamento TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_contas_pagar_vencimento
    ON tb_contas_pagar (dt_vencimento, status);

CREATE TABLE IF NOT EXISTS tb_contas_receber (
    id BIGSERIAL PRIMARY KEY,
    id_cliente BIGINT NOT NULL REFERENCES tb_clientes(id),
    id_categoria BIGINT NOT NULL REFERENCES tb_categorias_contas_receber(id),
    tipo_origem VARCHAR(50) NOT NULL,
    id_origem BIGINT,
    id_grupo_parcelas UUID,
    descricao VARCHAR(255),
    valor_original DECIMAL(15,2) NOT NULL CHECK (valor_original > 0),
    saldo_restante DECIMAL(15,2) NOT NULL CHECK (saldo_restante >= 0),
    dt_vencimento DATE NOT NULL,
    nr_parcela INTEGER NOT NULL CHECK (nr_parcela > 0),
    nr_total_parcelas INTEGER NOT NULL CHECK (nr_total_parcelas > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE'
        CHECK (status IN ('PENDENTE', 'PAGO_PARCIAL', 'PAGO', 'CANCELADO')),
    dt_emissao TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    dt_pagamento TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_movimento_financeiro (
    id BIGSERIAL PRIMARY KEY,
    tipo_movimento VARCHAR(20) NOT NULL
        CHECK (tipo_movimento IN ('CONTA_RECEBER', 'CONTA_PAGAR')),
    id_conta_pagar BIGINT REFERENCES tb_contas_pagar(id),
    id_conta_receber BIGINT REFERENCES tb_contas_receber(id),
    dt_movimento DATE NOT NULL,
    valor_movimento DECIMAL(15,2) NOT NULL CHECK (valor_movimento > 0),
    valor_acrescimo DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (valor_acrescimo >= 0),
    valor_desconto DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (valor_desconto >= 0),
    forma_pagamento VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_movimento_financeiro_origem CHECK (
        (tipo_movimento = 'CONTA_PAGAR' AND id_conta_pagar IS NOT NULL AND id_conta_receber IS NULL)
        OR (tipo_movimento = 'CONTA_RECEBER' AND id_conta_receber IS NOT NULL AND id_conta_pagar IS NULL)
    )
);

CREATE TABLE IF NOT EXISTS tb_grupos_tributarios (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    cfop_padrao VARCHAR(4) NOT NULL,
    origem_mercadoria INTEGER NOT NULL DEFAULT 0,
    csosn VARCHAR(4),
    icms_cst VARCHAR(3),
    icms_aliquota DECIMAL(5,2) NOT NULL DEFAULT 0,
    icms_mva_st DECIMAL(5,2) NOT NULL DEFAULT 0,
    icms_aliquota_st DECIMAL(5,2) NOT NULL DEFAULT 0,
    ipi_cst VARCHAR(2),
    ipi_aliquota DECIMAL(5,2) NOT NULL DEFAULT 0,
    pis_cst VARCHAR(2) NOT NULL DEFAULT '07',
    pis_aliquota DECIMAL(5,2) NOT NULL DEFAULT 0,
    cofins_cst VARCHAR(2) NOT NULL DEFAULT '07',
    cofins_aliquota DECIMAL(5,2) NOT NULL DEFAULT 0,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_produtos (
    id BIGSERIAL PRIMARY KEY,
    codigo_barras VARCHAR(14) UNIQUE,
    codigo_interno_loja VARCHAR(50),
    nome VARCHAR(150) NOT NULL,
    descricao TEXT,
    preco_custo DECIMAL(12,4) NOT NULL DEFAULT 0,
    preco_venda DECIMAL(12,2) NOT NULL DEFAULT 0,
    unidade_estoque VARCHAR(10) NOT NULL DEFAULT 'UN',
    unidade_venda VARCHAR(10) NOT NULL DEFAULT 'UN',
    peso_bruto DECIMAL(12,2) NOT NULL DEFAULT 0,
    peso_liquido DECIMAL(12,2) NOT NULL DEFAULT 0,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_produtos_codigo_interno ON tb_produtos(codigo_interno_loja);
CREATE INDEX IF NOT EXISTS idx_produtos_nome_ativos ON tb_produtos(nome) WHERE ativo = TRUE;

CREATE TABLE IF NOT EXISTS tb_produtos_fiscal (
    id_produto BIGINT PRIMARY KEY REFERENCES tb_produtos(id) ON DELETE CASCADE,
    ncm VARCHAR(8) NOT NULL,
    cest VARCHAR(7),
    id_grupo_tributario BIGINT NOT NULL REFERENCES tb_grupos_tributarios(id) ON DELETE RESTRICT,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_estoques (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    descricao VARCHAR(255),
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_produtos_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_produto BIGINT NOT NULL REFERENCES tb_produtos(id) ON DELETE CASCADE,
    id_estoque BIGINT NOT NULL REFERENCES tb_estoques(id) ON DELETE RESTRICT,
    quantidade DECIMAL(10,3) NOT NULL DEFAULT 0,
    estoque_minimo DECIMAL(10,3) NOT NULL DEFAULT 0,
    atualizado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_produto_por_estoque UNIQUE (id_produto, id_estoque)
);

CREATE INDEX IF NOT EXISTS idx_produtos_estoque_produto ON tb_produtos_estoque(id_produto);

CREATE TABLE IF NOT EXISTS tb_categoria_movimento_estoque (
    id BIGSERIAL PRIMARY KEY,
    nome VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tb_movimento_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_produto BIGINT NOT NULL REFERENCES tb_produtos(id),
    id_estoque BIGINT NOT NULL REFERENCES tb_estoques(id),
    id_usuario BIGINT NOT NULL REFERENCES public.tb_usuarios_gestao(id),
    quantidade DECIMAL(10,3) NOT NULL CHECK (quantidade > 0),
    tipo_movimento VARCHAR(10) NOT NULL CHECK (tipo_movimento IN ('ENTRADA', 'SAIDA')),
    id_categoria_movimento BIGINT NOT NULL REFERENCES tb_categoria_movimento_estoque(id),
    id_origem VARCHAR(50),
    observacao TEXT,
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_movimento_estoque_produto ON tb_movimento_estoque(id_produto, id_estoque);
CREATE INDEX IF NOT EXISTS idx_movimento_estoque_criado_em ON tb_movimento_estoque(criado_em);

CREATE TABLE IF NOT EXISTS tb_entradas_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_estoque BIGINT NOT NULL REFERENCES tb_estoques(id),
    id_fornecedor BIGINT NOT NULL REFERENCES tb_fornecedores(id),
    valor_despesa_adicional DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_despesa_adicional >= 0),
    id_usuario BIGINT NOT NULL REFERENCES public.tb_usuarios_gestao(id),
    valor_total DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_total >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'ABERTO'
        CHECK (status IN ('ABERTO', 'CONCLUIDA', 'CANCELADA')),
    criado_em TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tb_produtos_entradas_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_entrada_estoque BIGINT NOT NULL REFERENCES tb_entradas_estoque(id) ON DELETE CASCADE,
    id_produto BIGINT NOT NULL REFERENCES tb_produtos(id),
    valor_unitario DECIMAL(10,2) NOT NULL CHECK (valor_unitario > 0),
    valor_icms_st DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_icms_st >= 0),
    valor_ipi DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_ipi >= 0),
    valor_desconto DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (valor_desconto >= 0),
    rateio_despesa_adicional DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (rateio_despesa_adicional >= 0),
    valor_custo DECIMAL(10,2) NOT NULL CHECK (valor_custo > 0),
    valor_total DECIMAL(10,2) NOT NULL CHECK (valor_total > 0),
    quantidade DECIMAL(10,3) NOT NULL CHECK (quantidade > 0)
);

-- Seeds necessários para a entrada de estoque.
INSERT INTO tb_categoria_movimento_estoque (nome)
VALUES ('ENTRADA DE ESTOQUE')
ON CONFLICT (nome) DO NOTHING;

INSERT INTO tb_estoques (nome, descricao)
VALUES ('Estoque Geral', 'Local padrão de armazenamento e vendas')
ON CONFLICT (nome) DO NOTHING;

SET search_path TO public;

-- Tenants de desenvolvimento herdados da estrutura original.
INSERT INTO public.tb_empresas_gestao (nome_fantasia, email, telefone, schema, ativo)
VALUES
    ('Padaria Pão Quente', 'contato@paoquente.com.br', '11999999991', 'schema_pao_quente', TRUE),
    ('Mercadinho da Esquina', 'contato@mercadinho.com.br', '11999999992', 'schema_mercadinho', TRUE)
ON CONFLICT (email) DO UPDATE
SET nome_fantasia = EXCLUDED.nome_fantasia,
    telefone = EXCLUDED.telefone,
    schema = EXCLUDED.schema,
    ativo = EXCLUDED.ativo,
    atualizado_em = CURRENT_TIMESTAMP;

-- Credenciais de desenvolvimento. Troque-as ou remova-as antes de produção.
-- Senha: 123456
INSERT INTO public.tb_usuarios_gestao (id_empresa, nome, cpf, telefone, email, senha, ativo)
SELECT e.id, 'João da Padaria', '11111111111', '11988888881', 'joao@paoquente.com.br',
       '$2a$10$Wdisu5NesYrQ4eMAQdt9SekdU5QWTx7LWj/N2j8H/qabj/PaX7d.W', TRUE
FROM public.tb_empresas_gestao e
WHERE e.schema = 'schema_pao_quente'
ON CONFLICT (email) DO NOTHING;

INSERT INTO public.tb_usuarios_gestao (id_empresa, nome, cpf, telefone, email, senha, ativo)
SELECT e.id, 'Maria do Mercado', '22222222222', '11988888882', 'maria@mercadinho.com.br',
       '$2a$10$Wdisu5NesYrQ4eMAQdt9SekdU5QWTx7LWj/N2j8H/qabj/PaX7d.W', TRUE
FROM public.tb_empresas_gestao e
WHERE e.schema = 'schema_mercadinho'
ON CONFLICT (email) DO NOTHING;

COMMIT;
