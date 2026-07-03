create table if not exists tb_empresas(
    id BIGSERIAL PRIMARY KEY,
    razao_social varchar(100) not null,
    nome_fantasia varchar(100) not null,
    cnpj varchar(20) not null,
    ativo boolean default true,
    data_criacao timestamp default current_timestamp,
    data_atualizacao timestamp default current_timestamp 
);

create table if not exists tb_categorias_debito(
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

CREATE TABLE IF NOT EXISTS tb_debitos (
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
    FOREIGN KEY (id_categoria) REFERENCES tb_categorias_debito(id)
);


