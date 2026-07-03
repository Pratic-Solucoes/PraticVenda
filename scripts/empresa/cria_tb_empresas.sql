create table if not exists tb_empresas(
    id BIGSERIAL PRIMARY KEY,
    razao_social varchar(100) not null,
    nome_fantasia varchar(100) not null,
    cnpj varchar(20) not null,
    ativo boolean default true,
    data_criacao timestamp default current_timestamp,
    data_atualizacao timestamp default current_timestamp 
);