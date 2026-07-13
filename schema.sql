create table if not exists tb_empresas(
    id BIGSERIAL PRIMARY KEY,
    razao_social varchar(100) not null,
    nome_fantasia varchar(100) not null,
    cnpj varchar(20) not null,
    ativo boolean default true,
    data_criacao timestamp default current_timestamp,
    data_atualizacao timestamp default current_timestamp 
);

create table if not exists tb_categorias_contas_pagar(
    id BIGSERIAL PRIMARY KEY,
    nome varchar(255) not null
)

create table if not exists tb_clientes(
    id BIGSERIAL PRIMARY KEY,
    nome varchar(150) not null unique,
    tipo VARCHAR(50) not null,
    email varchar(200) unique default null,
    telefone varchar(20) unique default null,
    cpf varchar(14) unique default null,
    cnpj varchar(18) unique default null,
    contribuinte VARCHAR(50) default null,
    is_consumidor_final BOOLEAN DEFAULT FALSE,
    ie varchar(14) unique default null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp 
);

CREATE TABLE IF NOT EXISTS tb_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    razao_social VARCHAR(255) NOT NULL,
    cnpj VARCHAR(14) NOT NULL,
    inscricao_estadual VARCHAR(20),
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    UNIQUE (cnpj)
);

create table if not exists tb_usuarios_gestao(
    id bigint generated always as identity primary key ,
    id_empresa bigint not null,
    nome varchar(255) not null,
    cpf varchar(14) unique,
    telefone varchar(20) unique,
    email varchar(255) not null unique,
    senha varchar(255) not null,
    criado_em timestamp default current_timestamp,
    atualizado_em timestamp default current_timestamp ,
    ativo boolean default false,
    foreign key (id_empresa) references tb_empresas(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS tb_credenciais_empresa (
    id BIGINT PRIMARY KEY ,
    tp_ambiente SMALLINT NOT NULL, -- 1 para Produção, 2 para Homologação
    certificado_digital BYTEA NOT NULL, -- BYTEA garante espaço seguro para o binário no Postgres
    senha_criptografada VARCHAR(255) NOT NULL,
    id_csc VARCHAR(6) NOT NULL, -- O identificador sequencial fornecido pela SEFAZ
    csc_nfe VARCHAR(255) NOT NULL, -- O token secreto
    CONSTRAINT uq_empresa_ambiente UNIQUE (tp_ambiente) -- Garante uma config por ambiente
);

CREATE TABLE IF NOT EXISTS tb_config_fiscais_empresa (
    id BIGINT PRIMARY KEY ,
    cd_regime_tributario INT NOT NULL
);

create table if not exists tb_endereco_empresa(
    id BIGSERIAL PRIMARY KEY,
    logradouro varchar(100) not null,
    numero varchar(20) not null,
    bairro varchar(100) not null,
    cep varchar(8) not null,
    cd_cidade bigint not null,
    nome_cidade varchar(100) not null,
    estado varchar(2) not null,
    nome_pais varchar(50) not null default 'Brasil',
    cd_pais bigint not null default 1058,
    ativo boolean default true,
    data_criacao timestamp default current_timestamp,
    data_atualizacao timestamp default current_timestamp 
);

create table if not exists tb_enderecos_clientes (
    id BIGSERIAL PRIMARY KEY,
    id_cliente bigint not null,
    cep varchar(9) not null,
    logradouro varchar(255) not null,
    numero varchar(20) not null,
    bairro varchar(100) not null,
    municipio varchar(100) not null,
    uf varchar(2) not null,
    codigo_municipio varchar(7) not null,
    is_principal BOOLEAN DEFAULT FALSE,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp ,
    constraint fk_enderecos_clientes_clientes foreign key (id_cliente) references tb_clientes(id) on delete cascade on update cascade
);

CREATE TABLE IF NOT EXISTS tb_enderecos_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL,
    cep VARCHAR(8) NOT NULL,
    logradouro VARCHAR(255) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    bairro VARCHAR(100) NOT NULL,
    municipio VARCHAR(100) NOT NULL,
    uf CHAR(2) NOT NULL,
    codigo_municipio VARCHAR(7) NOT NULL,
    is_principal BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_endereco_fornecedor FOREIGN KEY (id_fornecedor) REFERENCES tb_fornecedores(id) ON DELETE CASCADE
);

CREATE TABLE if not exists tb_telefones_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL,
    ddd CHAR(2) NOT NULL,
    numero VARCHAR(15) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_telefone_fornecedor FOREIGN KEY (id_fornecedor) REFERENCES tb_fornecedores(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tb_contas_pagar (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL,
    id_categoria BIGINT,
    descricao VARCHAR(255) NOT NULL,
    nr_documento VARCHAR(255),
    nr_nota_fiscal VARCHAR(255),
    valor DECIMAL(15,2) NOT NULL,
    dt_entrada DATE NOT NULL,
    dt_vencimento DATE NOT NULL,
    nr_parcela INT NOT NULL,
    nr_total_parcelas INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDENTE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    FOREIGN KEY (id_fornecedor) REFERENCES tb_fornecedores(id),
    FOREIGN KEY (id_categoria) REFERENCES tb_categorias_contas_pagar(id)
);

-- 1. Tabela Principal: Dados Comerciais do Produto
CREATE TABLE IF NOT EXISTS tb_produtos (
    id SERIAL PRIMARY KEY,
    codigo_barras VARCHAR(14) UNIQUE, -- Suporta EAN-13, DUN-14, etc. Pode ser NULL se o produto não tiver.
    codigo_interno_loja VARCHAR(50),  -- Código customizado que o lojista cria (ex: "AUT-001")
    nome VARCHAR(150) NOT NULL,
    descricao TEXT,
    preco_custo DECIMAL(12,4) DEFAULT 0.0000 NOT NULL, -- 4 casas decimais para precisão de centavos em atacado
    preco_venda DECIMAL(12,2) DEFAULT 0.00 NOT NULL,
    unidade_estoque VARCHAR(10) DEFAULT 'UN' NOT NULL, -- Ex: 'UN', 'KG', 'CX', 'L'
    unidade_venda VARCHAR(10) DEFAULT 'UN' NOT NULL,   -- Geralmente igual à unidade de estoque no seu modelo simples
    peso_bruto DECIMAL(12,2) DEFAULT 0.00 NOT NULL,
    peso_liquido DECIMAL(12,2) DEFAULT 0.00 NOT NULL,
    ativo BOOLEAN DEFAULT TRUE NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Índices para buscas rápidas no PDV e na Gestão
CREATE INDEX IF NOT EXISTS idx_produtos_codigo_interno ON tb_produtos(codigo_interno_loja);
CREATE INDEX IF NOT EXISTS idx_produtos_nome ON tb_produtos(nome) WHERE ativo = TRUE;


-- 2. Tabela de Perfil Tributário (Regras Fiscais Reutilizáveis)
CREATE TABLE IF NOT EXISTS tb_grupos_tributarios (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL, -- Ex: "Revenda Padrão", "Produto ST", "Isento"
    
    -- Dados Padrão do Grupo
    cfop_padrao VARCHAR(4) NOT NULL,
    origem_mercadoria INT DEFAULT 0 NOT NULL,
    
    -- Simples Nacional
    csosn VARCHAR(4),
    
    -- Regime Normal e ST
    icms_cst VARCHAR(3),
    icms_aliquota DECIMAL(5,2) DEFAULT 0.00,
    icms_mva_st DECIMAL(5,2) DEFAULT 0.00,
    icms_aliquota_st DECIMAL(5,2) DEFAULT 0.00,
    
    -- IPI
    ipi_cst VARCHAR(2),
    ipi_aliquota DECIMAL(5,2) DEFAULT 0.00,
    
    -- PIS/COFINS
    pis_cst VARCHAR(2) DEFAULT '07',
    pis_aliquota DECIMAL(5,2) DEFAULT 0.00,
    cofins_cst VARCHAR(2) DEFAULT '07',
    cofins_aliquota DECIMAL(5,2) DEFAULT 0.00,
    
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


-- 3. Tabela Auxiliar: Dados Fiscais Específicos do Produto
CREATE TABLE IF NOT EXISTS tb_produtos_fiscal (
    id_produto INT PRIMARY KEY, -- Chave Primária e Estrangeira ao mesmo tempo (1:1 com produto)
    
    -- Natureza Física da Mercadoria
    ncm VARCHAR(8) NOT NULL,    -- Nomenclatura Comum do Mercosul (Obrigatório)
    cest VARCHAR(7),            -- Código Especificador da Substituição Tributária
    
    -- Vínculo com a Regra Fiscal (O Pulo do Gato)
    id_grupo_tributario INT NOT NULL,
    
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Relacionamentos
    FOREIGN KEY (id_produto) REFERENCES tb_produtos(id) ON DELETE CASCADE,
    FOREIGN KEY (id_grupo_tributario) REFERENCES tb_grupos_tributarios(id) ON DELETE RESTRICT
);

-- 1. Tabela de Locais de Estoque
CREATE TABLE IF NOT EXISTS tb_estoques (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    descricao VARCHAR(255),
    ativo BOOLEAN DEFAULT TRUE NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 2. Tabela de Saldo de Produtos por Estoque (Relacional)
CREATE TABLE IF NOT EXISTS tb_produtos_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_produto INT NOT NULL,
    id_estoque INT NOT NULL,
    quantidade DECIMAL(10,3) DEFAULT 0.000 NOT NULL,
    estoque_minimo DECIMAL(10,3) DEFAULT 0.000 NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (id_estoque) REFERENCES tb_estoques(id) ON DELETE RESTRICT,
    FOREIGN KEY (id_produto) REFERENCES tb_produtos(id) ON DELETE CASCADE,
    CONSTRAINT uq_produto_por_estoque UNIQUE (id_produto, id_estoque)
);

-- Índice para acelerar a busca de saldo de um produto específico nos estoques
CREATE INDEX IF NOT EXISTS idx_prod_est_produto ON tb_produtos_estoque(id_produto);

CREATE TABLE IF NOT EXISTS tb_categoria_movimento_estoque(
    id BIGSERIAL PRIMARY KEY,
    nome varchar(50) unique
);


-- 3. Tabela de Histórico de Movimentações (Rastreabilidade)
CREATE TABLE IF NOT EXISTS tb_movimento_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_produto INT NOT NULL,
    id_estoque INT NOT NULL,
    id_usuario BIGINT NOT NULL,
    quantidade DECIMAL(10,3) NOT NULL,
    tipo_movimento VARCHAR(10) NOT NULL,
    id_categoria_movimento BIGINT NOT NULL,
    id_origem VARCHAR(50),
    observacao TEXT,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (id_produto) REFERENCES tb_produtos(id),
    FOREIGN KEY (id_estoque) REFERENCES tb_estoques(id),
    FOREIGN KEY (id_usuario) REFERENCES tb_usuarios_gestao(id),
    FOREIGN KEY (id_categoria_movimento) REFERENCES tb_categoria_movimento_estoque(id),
    CONSTRAINT chk_tipo_movimento CHECK (tipo_movimento IN ('ENTRADA', 'SAIDA'))
);

-- Índices essenciais para relatórios de estoque e auditoria rápidos
CREATE INDEX IF NOT EXISTS idx_mov_produto_estoque ON tb_movimento_estoque(id_produto, id_estoque);
CREATE INDEX IF NOT EXISTS idx_mov_criado_em ON tb_movimento_estoque(criado_em);
