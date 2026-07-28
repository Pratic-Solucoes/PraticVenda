create table if not exists tb_organizacoes(
    id bigserial primary key,
    id_dono bigint not null,
    nome_organizacao varchar(150) not null unique,
    schema varchar(50) not null unique,
    ativo boolean not null default false,
    criado_em timestamp not null default now(),
    atualizado_em timestamp not null default now(),
    foreign key (id_dono) references tb_usuarios(id)
);