# 📋 PraticVenda — Análise de Gaps para o MVP

> Legenda: ✅ Implementado · ⚠️ Parcial / Incompleto · ❌ Faltando

---

## 1. 🔐 Autenticação & Usuário

| Funcionalidade | Status | Observação |
|---|---|---|
| Login | ✅ | Backend + Frontend prontos |
| Criar usuário | ✅ | Rota `POST /api/usuarios` |
| Buscar dados do usuário logado | ✅ | Rota `GET /api/usuario` |
| Alterar senha | ✅ | Rota `PUT /api/usuario/alterar-senha` |
| Logout (limpar token no frontend) | ⚠️ | Sem botão/ação de logout na sidebar |
| Recuperação de senha | ❌ | Não existe |
| Níveis de acesso / permissões | ❌ | Sem controle de perfil (admin, operador, etc.) |
| Multi-empresa (troca de empresa logada) | ❌ | Sem suporte a múltiplas empresas por usuário |

---

## 2. 🏢 Empresa

| Funcionalidade | Status | Observação |
|---|---|---|
| Visualizar dados da empresa | ⚠️ | Controller e model existem mas sem rota/página completa |
| Criar/editar dados da empresa | ❌ | Sem rota, sem página, sem JS |
| Configurar dados fiscais (CNPJ, IE, regime tributário) | ❌ | Tabela `tb_config_fiscais_empresa` existe mas sem CRUD |
| Configurar certificado digital | ❌ | Tabela `tb_credenciais_empresa` existe mas sem CRUD |
| Configurar endereço da empresa | ❌ | Tabela `tb_endereco_empresa` existe mas sem CRUD |

---

## 3. 🛍️ Produtos / Catálogo

| Funcionalidade | Status | Observação |
|---|---|---|
| Listar produtos | ❌ | **Nenhuma tabela no schema, nenhum model, nenhum controller** |
| Criar produto | ❌ | Ausente em toda a stack |
| Editar produto | ❌ | Ausente |
| Inativar produto | ❌ | Ausente |
| Categorias de produtos | ❌ | Ausente |
| Estoque (saldo em estoque) | ❌ | Sem tabela de estoque no schema |
| Unidade de medida | ❌ | Ausente |
| Código de barras / SKU | ❌ | Ausente |
| Preço de custo e preço de venda | ❌ | Ausente |

> [!CAUTION]
> Este é o **maior gap do MVP**. Sem o módulo de Produtos, o PDV não pode funcionar e o módulo de Vendas não pode ser construído.

---

## 4. 🛒 Vendas / PDV

| Funcionalidade | Status | Observação |
|---|---|---|
| Tela do PDV (layout) | ⚠️ | HTML da tela existe, mas é apenas layout estático sem JS |
| Buscar produto por código/nome no PDV | ❌ | Nenhum JS, nenhuma rota de busca |
| Adicionar item à venda | ❌ | Sem JS e sem backend |
| Remover item da venda | ❌ | Ausente |
| Alterar quantidade de item | ❌ | Ausente |
| Selecionar cliente na venda | ❌ | Campo existe no HTML mas sem bind |
| Aplicar desconto | ❌ | Ausente |
| Finalizar venda (registrar no banco) | ❌ | Sem tabela `tb_vendas`, sem rota, sem controller |
| Salvar venda em aberto (orçamento) | ❌ | Botão existe no HTML, sem funcionalidade |
| Cancelar venda | ❌ | Ausente |
| Histórico de vendas | ❌ | Sem tela, sem rota |
| Detalhe de uma venda | ❌ | Ausente |
| Módulo de pagamento (dinheiro/cartão/pix) | ❌ | Sem modal de pagamento, sem lógica |
| Geração de NF-e / NFC-e | ❌ | Estrutura prevista no `notas.txt` mas sem implementação |
| Impressão de cupom/recibo | ❌ | Ausente |

> [!CAUTION]
> O módulo de Vendas é **completamente não funcional** — depende do módulo de Produtos e de tabelas que ainda não existem no banco.

---

## 5. 👥 Clientes

| Funcionalidade | Status | Observação |
|---|---|---|
| Listar clientes | ✅ | Rota e frontend prontos |
| Criar cliente | ⚠️ | Apenas pelo modal de edição; sem modal/fluxo dedicado de criação |
| Editar cliente (dados principais) | ✅ | JS `editarCliente.js` implementado |
| Inativar / Excluir cliente | ❌ | Sem rota `DELETE` e sem botão na UI |
| Gerenciar endereços do cliente | ⚠️ | Rotas existem no backend; integração frontend incompleta |
| Gerenciar telefones do cliente | ⚠️ | Rotas existem no backend; integração frontend incompleta |
| Busca/filtro por nome, CPF, CNPJ | ❌ | Listagem sem filtro funcional |

---

## 6. 🏭 Fornecedores

| Funcionalidade | Status | Observação |
|---|---|---|
| Listar fornecedores | ✅ | Rota e frontend prontos |
| Criar fornecedor (modal) | ✅ | `criarFornecedorModal.js` implementado |
| Editar fornecedor | ✅ | `editarFornecedor.js` implementado |
| Inativar / Excluir fornecedor | ❌ | Sem rota `DELETE` e sem botão na UI |
| Gerenciar endereços do fornecedor | ⚠️ | Tabela existe no banco, mas sem integração completa no frontend |
| Gerenciar telefones do fornecedor | ⚠️ | Tabela existe no banco, mas sem integração completa no frontend |
| Busca/filtro | ❌ | Listagem sem filtro funcional |

---

## 7. 💸 Contas a Pagar

| Funcionalidade | Status | Observação |
|---|---|---|
| Listar contas a pagar | ✅ | Implementado |
| Criar conta a pagar (avulso) | ✅ | Modal implementado |
| Editar conta a pagar | ✅ | `editarContaPagar.js` implementado |
| Visualizar detalhes | ✅ | `visualizarContaPagar.js` implementado |
| Pagar conta (total) | ✅ | Rota e JS implementados |
| Pagar conta parcialmente | ⚠️ | Rota `POST /contas-pagar/{id}/pagar-parcial` tem `nil` como handler — **não implementado** |
| Cancelar conta a pagar | ❌ | Sem rota `DELETE`/cancelamento e sem botão na UI |
| Filtro por status/data/fornecedor | ❌ | Sem filtro funcional na listagem |
| Categorias de contas a pagar | ✅ | Listar e criar implementados |
| Editar / excluir categoria | ❌ | Sem rota e sem UI |

---

## 8. 💳 Formas de Pagamento

| Funcionalidade | Status | Observação |
|---|---|---|
| Listar formas de pagamento | ✅ | Implementado |
| Criar forma de pagamento | ✅ | Implementado |
| Editar forma de pagamento | ✅ | Implementado |
| Excluir forma de pagamento | ❌ | Sem rota `DELETE` e sem botão na UI |

---

## 9. 📊 Dashboard

| Funcionalidade | Status | Observação |
|---|---|---|
| KPIs resumo (contas a pagar pendentes, etc.) | ✅ | `kpis.js` e rota `/api/dashboard/resumo` |
| Gráficos | ⚠️ | `graficos.js` existe mas dados são mockados/limitados |
| KPIs de Vendas (faturamento do dia/mês) | ❌ | Impossível sem o módulo de Vendas |
| Alertas de contas vencidas | ❌ | Ausente |

---

## 10. ⚙️ Configurações Gerais

| Funcionalidade | Status | Observação |
|---|---|---|
| Configuração de usuário (alterar senha/nome) | ✅ | Tela e JS implementados |
| Configuração da empresa | ❌ | Ausente (ver módulo Empresa) |
| Configuração fiscal | ❌ | Ausente |

---

## 🗺️ Resumo de Prioridades para o MVP

```
Alta Prioridade (Bloqueiam o núcleo do negócio)
├── ❌ Módulo de Produtos (tabelas, backend, frontend completo)
├── ❌ Tabelas de Vendas no banco de dados
├── ❌ Módulo de Vendas / PDV (backend + frontend funcional)
└── ❌ Modal de Pagamento no PDV

Média Prioridade (Completar o que já está pela metade)
├── ⚠️ Pagamento parcial de conta a pagar (handler nulo)
├── ⚠️ Endereços e telefones de clientes/fornecedores (integração frontend)
├── ⚠️ Filtros nas listagens (clientes, fornecedores, contas a pagar)
├── ❌ Excluir / Cancelar registros (clientes, fornecedores, contas)
└── ❌ Editar / Excluir categorias de contas a pagar

Baixa Prioridade (Podem vir após o MVP)
├── ❌ Recuperação de senha
├── ❌ Níveis de permissão de usuário
├── ❌ Configurações da empresa e fiscais
├── ❌ Geração de NF-e / NFC-e
└── ❌ Logout explícito (botão na sidebar)
```
