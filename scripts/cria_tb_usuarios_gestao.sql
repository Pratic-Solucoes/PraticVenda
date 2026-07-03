create table if not exists tb_usuarios_gestao(
    id bigserial primary key ,
    id_empresa bigint not null,
    nome varchar(255) not null,
    cpf varchar(14) unique,
    telefone varchar(20) unique,
    email varchar(255) not null unique,
    senha varchar(255) not null,
    criado_em timestamp default current_timestamp,
    atualizado_em timestamp default current_timestamp ,
    ativo boolean default false,
    foreign key (id_empresa) references tb_empresas_gestao(id) on delete cascade
);
