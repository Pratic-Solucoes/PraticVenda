-- Criação de duas empresas de gestão (Tenants / Accounts)
INSERT INTO tb_empresas_gestao (nome_fantasia, email, telefone, schema, ativo)
VALUES 
    ('Padaria Pão Quente', 'contato@paoquente.com.br', '11999999991', 'schema_pao_quente', true),
    ('Mercadinho da Esquina', 'contato@mercadinho.com.br', '11999999992', 'schema_mercadinho', true)
ON CONFLICT DO NOTHING;

-- Criação de um usuário para cada empresa 
-- A senha para ambos é: 123456
-- O hash bcrypt para '123456' é: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy

-- Inserindo usuário para a Padaria Pão Quente
INSERT INTO tb_usuarios_gestao (id_empresa, nome, cpf, telefone, email, senha, ativo)
VALUES (
    (SELECT id FROM tb_empresas_gestao WHERE email = 'contato@paoquente.com.br'),
    'João da Padaria',
    '11111111111',
    '11988888881',
    'joao@paoquente.com.br',
    '$2a$10$Wdisu5NesYrQ4eMAQdt9SekdU5QWTx7LWj/N2j8H/qabj/PaX7d.W',
    true
)
ON CONFLICT DO NOTHING;

-- Inserindo usuário para o Mercadinho
INSERT INTO tb_usuarios_gestao (id_empresa, nome, cpf, telefone, email, senha, ativo)
VALUES (
    (SELECT id FROM tb_empresas_gestao WHERE email = 'contato@mercadinho.com.br'),
    'Maria do Mercado',
    '22222222222',
    '11988888882',
    'maria@mercadinho.com.br',
    '$2a$10$Wdisu5NesYrQ4eMAQdt9SekdU5QWTx7LWj/N2j8H/qabj/PaX7d.W',
    true
)
ON CONFLICT DO NOTHING;

-- Criação dos schemas de teste para simular o ambiente
CREATE SCHEMA IF NOT EXISTS schema_pao_quente;
CREATE SCHEMA IF NOT EXISTS schema_mercadinho;
