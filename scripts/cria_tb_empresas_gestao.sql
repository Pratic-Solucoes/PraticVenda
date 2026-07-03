create table if not exists tb_empresas_gestao(
    id bigserial primary key,
    nome_fantasia varchar(255) not null unique,
    email varchar(255) not null unique,
    telefone varchar(20) not null unique,
    schema varchar(50) not null unique,
    ativo bool not null default true
);