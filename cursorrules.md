# Diretrizes Gerais de Comportamento
- Nunca gere código de backend a menos que eu peça explicitamente no chat ou composer.
- Para o frontend, você tem autonomia para sugerir e criar componentes completos, seguindo estritamente o design pattern definido na seção de Frontend deste arquivo.
- Sempre priorize o desempenho, a legibilidade e a simplicidade do código. Evite abstrações desnecessárias ou "mágicas".

# Regras de Desenvolvimento: FrontEnd
- **Stack:** HTML5, CSS3 (Bootstrap) e JavaScript Puro (Vanilla JS).
- **Estrutura de Pastas (Diretório `/web`):**
  - HTML: Localizado em `/web/template/` -> Dividido em `pages/` (telas comuns) e `componentes/` (modais personalizados).
  - CSS/JS: Localizados em `/web/static/`.
- **Arquitetura JavaScript:**
  - Divida os arquivos JS estritamente por **domínio** (Ex: `cliente.js`, `fornecedores.js`). Cada arquivo é responsável pelas funcionalidades daquele domínio.
  - **Quebre a lógica em funções de responsabilidade única**. Exemplo: Uma função exclusiva para fazer o `fetch` na API, outra função isolada para popular a tabela/modal, outra para capturar o evento do formulário, etc.
- **Reutilização (`/web/static/utils/`):**
  - Sempre verifique e utilize as funções utilitárias presentes na pasta `utils` antes de criar lógica nova (Ex: funções para capturar tokens, validar respostas do servidor, formatações, etc.).

# Regras de Desenvolvimento: BackEnd
- **Arquitetura:** MVC (Model-View-Controller) isolando explicitamente as camadas de **Service** e **Repository**.
- **Desacoplamento:** Uso obrigatório de **interfaces** em todas as camadas para a inicialização e comunicação dos componentes.
- **Injeção de Dependência:** Toda a árvore de dependências deve ser injetada de forma manual e explícita, sendo inicializada centralizadamente no arquivo `/cmd/api/main.go`.
- **Filosofia de Código:** Código sem abstrações complexas, focado em legibilidade, manutenibilidade e fluxo explícito de dados/erros.