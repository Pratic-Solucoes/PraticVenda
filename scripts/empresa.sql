create table if not exists tb_empresas(
    id BIGSERIAL PRIMARY KEY,
    razao_social varchar(100) not null,
    nome_fantasia varchar(100) not null,
    cnpj varchar(20) not null,
    ativo boolean default true,
    data_criacao timestamp default current_timestamp,
    data_atualizacao timestamp default current_timestamp 
);

create table if not exists tb_endereco_empresa(
    id bigserial primary KEY,
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

CREATE TABLE IF NOT EXISTS tb_credenciais_empresa (
    id bigserial PRIMARY KEY ,
    tp_ambiente smallint NOT NULL, -- 1 para Produção, 2 para Homologação
    certificado_digital bytea NOT NULL, -- BYTEA garante espaço seguro para o binário no Postgres
    senha_criptografada varchar(255) NOT NULL,
    id_csc varchar(6) NOT NULL, -- O identificador sequencial fornecido pela SEFAZ
    csc_nfe varchar(255) NOT NULL -- O token secreto
);

CREATE TABLE IF NOT EXISTS tb_config_fiscais_empresa (
    id bigserial primary key ,
    cd_regime_tributario int not null
);
