CREATE TABLE IF NOT EXISTS tb_debitos (
    id BIGSERIAL PRIMARY KEY,
    id_fornecedor BIGINT NOT NULL,
    id_categoria BIGINT,
    descricao VARCHAR(255) NOT NULL,
    nr_documento VARCHAR(255),
    nr_nota_fiscal VARCHAR(255),
    valor DECIMAL(15,2) NOT NULL,
    dt_entrada DATE NOT NULL,
    dt_vencimento DATE NOT NULL,
    nr_parcela INT NOT NULL,
    nr_total_parcelas INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDENTE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    FOREIGN KEY (id_fornecedor) REFERENCES tb_fornecedores(id),
    FOREIGN KEY (id_categoria) REFERENCES tb_categorias_debito(id)
);
