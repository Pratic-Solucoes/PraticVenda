create table if not exists tb_usuarios(
    id bigserial primary key ,
    nome varchar(255) not null,
    username varchar(50) not null unique,
    email varchar(255) not null unique,
    celular varchar(20) unique,
    senha varchar(255) not null,
    termos_aceitos bool not null default false,
    termos_aceitos_em timestamp default null,
    ativo boolean default false,
    ativo boolean default false,
    criado_em timestamp default current_timestamp,
    atualizado_em timestamp default current_timestamp ,
);