# Entrada de Estoque — MVP

## Contexto

O sistema já possui a estrutura de estoques (locais/armazéns), produtos vinculados a estoques e a tabela `tb_movimento_estoque` no banco. A rota `POST /api/estoques/{id}/entrada` já está declarada no router, mas aponta para `nil`. O controller `EntradaEstoque` existe mas está incompleto.

O objetivo é implementar o fluxo completo de **lançamento de entrada de estoque**: o operador escolhe um estoque, seleciona produtos e informa a quantidade que entrou, e o sistema incrementa o saldo e registra a movimentação.

---

## Campos da Entrada de Estoque (MVP)

Para o MVP, mantemos o mínimo necessário e funcional:

| Campo | Obrigatório | Descrição |
|---|---|---|
| `id_produto` | ✅ | Produto que está entrando |
| `quantidade` | ✅ | Quantidade que entrou (positivo) |
| `observacao` | ❌ | Texto livre (ex: "Compra NF 001") |

> [!NOTE]
> O `id_estoque` vem da URL (`/api/estoques/{id}/entrada`) e o `id_usuario` é extraído do JWT no contexto.
> Campos como `id_categoria_movimento`, `id_origem`, custo da entrada e fornecedor ficam para uma fase posterior.

---

## Open Questions

> [!IMPORTANT]
> A tabela `tb_movimento_estoque` exige `id_categoria_movimento NOT NULL` (FK para `tb_categoria_movimento_estoque`). Para o MVP, precisamos de uma estratégia:
> - **Opção A (recomendada):** Criar uma categoria padrão "Entrada Manual" com ID fixo (ex: ID=1) via seed no banco. O backend usa esse ID hardcoded no MVP.
> - **Opção B:** Tornar a FK nullable temporariamente no schema.
> - **Opção C:** Adicionar o campo no formulário do frontend.
>
> **Por favor, confirme qual opção prefere antes de executarmos.**

---

## Proposed Changes

### Backend — Model

#### [MODIFY] [produto.go](file:///c:/Users/leona/OneDrive/Área de Trabalho/PraticVenda/internal/model/produto.go)
- O `ProdutoEstoqueEntrada` já existe com `id_produto` e `quantidade`. 
- Adicionar campo `Observacao *string` ao struct.

---

### Backend — Repository

#### [MODIFY] [estoque.go](file:///c:/Users/leona/OneDrive/Área de Trabalho/PraticVenda/internal/repository/estoque.go)
- Adicionar método `RegistrarEntrada(ctx, tx, idEstoque, idProduto, idUsuario int64, quantidade float64, observacao *string) error`
  - Incrementa `quantidade` em `tb_produtos_estoque` (UPDATE ... SET quantidade = quantidade + $1)
  - Insere linha em `tb_movimento_estoque` com `tipo_movimento = 'ENTRADA'`

#### [MODIFY] [repository.go](file:///c:/Users/leona/OneDrive/Área de Trabalho/PraticVenda/internal/repository/repository.go)
- Adicionar `RegistrarEntrada(...)` na interface `Estoques`

---

### Backend — Service

#### [MODIFY] [estoque.go](file:///c:/Users/leona/OneDrive/Área de Trabalho/PraticVenda/internal/service/estoque.go)
- Adicionar método `RegistrarEntrada(ctx, idEstoque, idUsuario int64, input *model.ProdutoEstoqueEntrada) error`
  - Valida: quantidade > 0, produto pertence ao estoque
  - Abre transação, chama repository, commita

#### [MODIFY] [service.go](file:///c:/Users/leona/OneDrive/Área de Trabalho/PraticVenda/internal/service/service.go)
- Adicionar `RegistrarEntrada(...)` na interface `Estoques`

---

### Backend — Controller

#### [MODIFY] [estoque.go (controller)](file:///c:/Users/leona/OneDrive/Área de Trabalho/PraticVenda/internal/controller/estoque.go)
- Completar `EntradaEstoque`:
  - Lê `{id}` da URL (id do estoque)
  - Decodifica body `ProdutoEstoqueEntrada`
  - Extrai `id_usuario` do JWT no contexto
  - Chama `service.Estoques.RegistrarEntrada(...)`

---

### Backend — Router

#### [MODIFY] [router.go](file:///c:/Users/leona/OneDrive/Área de Trabalho/PraticVenda/internal/router/router.go)
- Trocar `nil` por `auth.Autenticar(c.Estoques.EntradaEstoque)` na rota `POST /estoques/{id}/entrada`

---

### Frontend — Modal

#### [NEW] modalEntradaEstoque.html
- Modal simples com:
  - Select de produto (populado com os produtos do estoque selecionado)
  - Campo de quantidade (number, min=0.001)
  - Campo de observação (textarea opcional)
  - Botão "Confirmar Entrada"

---

### Frontend — JS

#### [NEW] entradaEstoque.js
- `setupEntradaEstoque()` — registra listener no formulário do modal
- `fetchEntrada(idEstoque, body)` — chama `POST /api/estoques/{id}/entrada`
- Após sucesso: fecha modal, recarrega tabela de produtos do estoque

#### [MODIFY] main.js (estoque)
- Importar e chamar `setupEntradaEstoque()`

---

### Frontend — HTML (estoques.html)

#### [MODIFY] [estoques.html](file:///c:/Users/leona/OneDrive/Área de Trabalho/PraticVenda/web/template/pages/estoques.html)
- Adicionar botão "Entrada de Estoque" ao lado do botão "Novo Estoque" (ativado apenas quando um estoque está selecionado)
- Incluir `{{template "modalEntradaEstoque" .}}`

---

## Verification Plan

### Manual
1. Acessar `/estoques`, selecionar um estoque com produtos vinculados.
2. Clicar em "Entrada de Estoque", preencher produto + quantidade e confirmar.
3. Verificar que a quantidade na tabela aumentou.
4. Verificar linha inserida em `tb_movimento_estoque` via banco.
5. Testar com quantidade inválida (0 ou negativa) — deve retornar erro.
6. Testar com produto não vinculado ao estoque — deve retornar erro.
