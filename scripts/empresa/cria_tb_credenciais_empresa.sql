CREATE TABLE IF NOT EXISTS tb_credenciais_empresa (
    id bigserial PRIMARY KEY ,
    tp_ambiente smallint NOT NULL, -- 1 para Produção, 2 para Homologação
    certificado_digital bytea NOT NULL, -- BYTEA garante espaço seguro para o binário no Postgres
    senha_criptografada varchar(255) NOT NULL,
    id_csc varchar(6) NOT NULL, -- O identificador sequencial fornecido pela SEFAZ
    csc_nfe varchar(255) NOT NULL -- O token secreto
);