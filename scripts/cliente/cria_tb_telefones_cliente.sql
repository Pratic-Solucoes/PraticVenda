CREATE TABLE if not exists tb_telefones_clientes (
    id BIGSERIAL PRIMARY KEY,
    id_cliente BIGINT NOT NULL,
    ddd CHAR(2) NOT NULL,
    numero VARCHAR(15) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_telefone_cliente FOREIGN KEY (id_cliente) REFERENCES tb_clientes(id) ON DELETE CASCADE
);