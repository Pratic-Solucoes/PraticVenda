-- 1. Tabela de Locais de Estoque
CREATE TABLE tb_estoques (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    descricao VARCHAR(255),
    ativo BOOLEAN DEFAULT TRUE NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Comentário para organização: Inserindo o estoque padrão/geral
INSERT INTO tb_estoques (nome, descricao) VALUES ('Estoque Geral', 'Local padrão de armazenamento e vendas');


-- 2. Tabela de Saldo de Produtos por Estoque (Relacional)
CREATE TABLE tb_produtos_estoque (
    id BIGSERIAL PRIMARY KEY,
    id_produto INT NOT NULL, -- Supondo FK para sua tabela de produtos principal
    id_estoque INT NOT NULL,
    quantidade DECIMAL(10,3) DEFAULT 0.000 NOT NULL,
    estoque_minimo DECIMAL(10,3) DEFAULT 0.000 NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (id_estoque) REFERENCES tb_estoques(id) ON DELETE RESTRICT,
    FOREIGN KEY (id_produto) REFERENCES tb_produtos(id) ON DELETE CASCADE,
    -- CONSTRAINT fk_prod_est_produto FOREIGN KEY (id_produto) REFERENCES tb_produtos(id_produto) ON DELETE RESTRICT,
    -- Garante que não haverá duplicidade de produto no mesmo estoque
    CONSTRAINT uq_produto_por_estoque UNIQUE (id_produto, id_estoque)
);

-- Índice para acelerar a busca de saldo de um produto específico nos estoques
CREATE INDEX idx_prod_est_produto ON tb_produtos_estoque(id_produto);

create table tb_categoria_movimento_estoque(
    id BIGSERIAL PRIMARY KEY,
    nome varchar(50) unique
);


-- 3. Tabela de Histórico de Movimentações (Rastreabilidade)
CREATE TABLE tb_movimento_estoque (
    id BIGSERIAL PRIMARY KEY, -- BIGSERIAL para aguentar milhões de movimentações ao longo dos anos
    id_produto INT NOT NULL,
    id_estoque INT NOT NULL,
    id_usuario BIGINT NOT NULL, -- Quem realizou a operação (Corrigido para BIGINT para bater com tb_usuarios_gestao)
    quantidade DECIMAL(10,3) NOT NULL, -- Sempre gravada em valor absoluto (positivo)
    tipo_movimento VARCHAR(10) NOT NULL, -- 'ENTRADA' ou 'SAIDA'
    id_categoria_movimento BIGINT NOT NULL, -- Deve ser BIGINT para bater com o ID da tabela de categoria
    id_origem VARCHAR(50), -- ID do documento que gerou isso (ex: ID da venda, número da NF-e)
    observacao TEXT,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (id_produto) REFERENCES tb_produtos(id),
    FOREIGN KEY (id_estoque) REFERENCES tb_estoques(id),
    FOREIGN KEY (id_usuario) REFERENCES tb_usuarios_gestao(id),
    FOREIGN KEY (id_categoria_movimento) REFERENCES tb_categoria_movimento_estoque(id),
    -- Proteção extra via Banco de Dados para evitar lixo nos ENUMs de texto
    CONSTRAINT chk_tipo_movimento CHECK (tipo_movimento IN ('ENTRADA', 'SAIDA'))
);
-- Índices essenciais para relatórios de estoque e auditoria rápidos
CREATE INDEX idx_mov_produto_estoque ON tb_movimento_estoque(id_produto, id_estoque);
CREATE INDEX idx_mov_criado_em ON tb_movimento_estoque(criado_em);