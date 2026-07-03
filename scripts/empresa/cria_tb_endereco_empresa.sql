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