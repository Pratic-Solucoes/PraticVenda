CREATE TABLE IF NOT EXISTS tb_enderecos_fornecedores (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL,
    cep VARCHAR(8) NOT NULL,
    logradouro VARCHAR(255) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    bairro VARCHAR(100) NOT NULL,
    municipio VARCHAR(100) NOT NULL,
    uf CHAR(2) NOT NULL,
    codigo_municipio VARCHAR(7) NOT NULL,
    is_principal BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_endereco_fornecedor FOREIGN KEY (id_fornecedor) REFERENCES tb_fornecedores(id) ON DELETE CASCADE
);