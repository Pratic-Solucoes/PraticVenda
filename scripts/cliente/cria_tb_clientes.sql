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