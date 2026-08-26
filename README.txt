Um sistema de gestão empresarial (ERP) com backend em Go e frontend reativo em Javascript puro. O sistema adota uma arquitetura de SPA (Single Page Application), onde o Go serve o HTML inicial e uma API RESTful para as operações dinâmicas.
O Sistema é focado em vendas, principalmente em pequenas/médias empresas, feito para facilitar a vida do usuário, provendo uma gestão completa, sem complicação.

Tecnologias:

- Backend
Linguagem: Go (Golang) v1.25.0
Roteador HTTP: Chi Router
Banco de Dados: PostgreSQL
Segurança: Bcrypt para senhas e JWT (golang-jwt/v5) para sessões.

- Frontend
* **Javascript (Vanilla)**: Manipulação do DOM e consumo de API via `fetch`.
* **Renderização**: HTML com Go Templates para componentização.
* **Estilização**: Bootstrap 5 e CSS customizado.

---

## Banco de dados

Toda a estrutura foi consolidada em `database/001_inicializar_banco.sql`. Para criar o banco local já existente, execute:

```bash
psql "$DATABASE_URL" -f database/001_inicializar_banco.sql
```

No Docker Compose, o arquivo é executado automaticamente somente quando o volume do PostgreSQL é criado pela primeira vez.

---

## 🏛️ Arquitetura

O projeto segue um padrão em camadas, separando as responsabilidades:
- **`Controller`**: Camada de entrada que lida com as requisições HTTP.
- **`Service`**: Orquestra as regras de negócio da aplicação.
- **`Repository`**: Responsável pelo acesso e manipula
