create table if not exists tb_clientes(
    id bigserial primary key,
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
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp ,
    constraint fk_enderecos_clientes_clientes foreign key (id_cliente) references tb_clientes(id) on delete cascade on update cascade
);

CREATE TABLE if not exists tb_telefones_clientes (
    id BIGSERIAL PRIMARY KEY,
    id_cliente BIGINT NOT NULL,
    ddd CHAR(2) NOT NULL,
    numero VARCHAR(15) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_telefone_cliente FOREIGN KEY (id_cliente) REFERENCES tb_clientes(id) ON DELETE CASCADE
);
