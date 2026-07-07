-- 1. Tabela Principal: Dados Comerciais do Produto
CREATE TABLE tb_produtos (
    id SERIAL PRIMARY KEY,
    codigo_barras VARCHAR(14) UNIQUE, -- Suporta EAN-13, DUN-14, etc. Pode ser NULL se o produto não tiver.
    codigo_interno_loja VARCHAR(50),  -- Código customizado que o lojista cria (ex: "AUT-001")
    nome VARCHAR(150) NOT NULL,
    descricao TEXT,
    preco_custo DECIMAL(12,4) DEFAULT 0.0000 NOT NULL, -- 4 casas decimais para precisão de centavos em atacado
    preco_venda DECIMAL(12,2) DEFAULT 0.00 NOT NULL,
    unidade_estoque VARCHAR(10) DEFAULT 'UN' NOT NULL, -- Ex: 'UN', 'KG', 'CX', 'L'
    unidade_venda VARCHAR(10) DEFAULT 'UN' NOT NULL,   -- Geralmente igual à unidade de estoque no seu modelo simples
    ativo BOOLEAN DEFAULT TRUE NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Índices para buscas rápidas no PDV e na Gestão
CREATE INDEX idx_produtos_codigo_interno ON tb_produtos(codigo_interno_loja);
CREATE INDEX idx_produtos_nome ON tb_produtos(nome) WHERE ativo = TRUE;


-- 2. Tabela de Perfil Tributário (Regras Fiscais Reutilizáveis)
CREATE TABLE tb_grupos_tributarios (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL, -- Ex: "Revenda Padrão", "Produto ST", "Isento"
    
    -- Dados Padrão do Grupo
    cfop_padrao VARCHAR(4) NOT NULL,
    origem_mercadoria INT DEFAULT 0 NOT NULL,
    
    -- Simples Nacional
    csosn VARCHAR(4),
    
    -- Regime Normal e ST
    icms_cst VARCHAR(3),
    icms_aliquota DECIMAL(5,2) DEFAULT 0.00,
    icms_mva_st DECIMAL(5,2) DEFAULT 0.00,
    icms_aliquota_st DECIMAL(5,2) DEFAULT 0.00,
    
    -- IPI
    ipi_cst VARCHAR(2),
    ipi_aliquota DECIMAL(5,2) DEFAULT 0.00,
    
    -- PIS/COFINS
    pis_cst VARCHAR(2) DEFAULT '07',
    pis_aliquota DECIMAL(5,2) DEFAULT 0.00,
    cofins_cst VARCHAR(2) DEFAULT '07',
    cofins_aliquota DECIMAL(5,2) DEFAULT 0.00,
    
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


-- 3. Tabela Auxiliar: Dados Fiscais Específicos do Produto
CREATE TABLE tb_produtos_fiscal (
    id_produto INT PRIMARY KEY, -- Chave Primária e Estrangeira ao mesmo tempo (1:1 com produto)
    
    -- Natureza Física da Mercadoria
    ncm VARCHAR(8) NOT NULL,    -- Nomenclatura Comum do Mercosul (Obrigatório)
    cest VARCHAR(7),            -- Código Especificador da Substituição Tributária
    
    -- Vínculo com a Regra Fiscal (O Pulo do Gato)
    id_grupo_tributario INT NOT NULL,
    
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Relacionamentos
    FOREIGN KEY (id_produto) REFERENCES tb_produtos(id) ON DELETE CASCADE,
    FOREIGN KEY (id_grupo_tributario) REFERENCES tb_grupos_tributarios(id) ON DELETE RESTRICT
);

