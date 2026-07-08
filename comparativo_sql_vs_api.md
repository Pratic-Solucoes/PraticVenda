# Comparativo: SQL do Projeto vs. API de Nota Fiscal

Fontes analisadas: [notas.txt](file:///c:/Users/leona/OneDrive/Área%20de%20Trabalho/PraticVenda/notas.txt) × SQL files em [/scripts](file:///c:/Users/leona/OneDrive/Área%20de%20Trabalho/PraticVenda/scripts/)

---

## 1. Empresa (`empresa.sql` × Seção "Criar Empresa")

### ✅ Campos presentes

| Campo API | Coluna SQL | Observação |
|---|---|---|
| `legalName` (razão social) | `razao_social` | ✅ |
| `name` (nome fantasia) | `nome_fantasia` | ✅ |
| `federalTaxNumber` (CNPJ) | `cnpj` | ✅ |
| `taxRegime` (regime tributário) | `cd_regime_tributario` em `tb_config_fiscais_empresa` | ✅ |
| `address.street` | `logradouro` | ✅ |
| `address.number` | `numero` | ✅ |
| `address.district` | `bairro` | ✅ |
| `address.postalCode` | `cep` | ✅ |
| `address.city.name` | `nome_cidade` | ✅ |
| `address.city.code` | `cd_cidade` | ✅ |
| `address.city.state` | `estado` | ✅ |
| `address.country.name` | `nome_pais` | ✅ |
| `address.country.code` | `cd_pais` | ✅ |

### ❌ Campos ausentes

| Campo API | Descrição | Urgência |
|---|---|---|
| `stateTaxNumber` | Inscrição Estadual (IE) | 🔴 Alta — exigida para NF-e |
| `cityTaxNumber` | Inscrição Municipal (IM) | 🟡 Média — exigida para NFS-e |
| `email` | E-mail da empresa | 🔴 Alta — para envio de NF ao cliente |
| `phone` | Telefone fixo | 🟡 Média |
| `mobilePhone` | Celular | 🟡 Média |
| `address.additionalInformation` | Complemento do endereço | 🟡 Média |
| `specialTaxRegime` | Regime especial tributário (municipal) | 🟡 Média |
| `simplesNacionalTaxRegime` | Sub-regime Simples Nacional | 🟡 Média |
| `economicActivities[].code` | Código CNAE | 🔴 Alta — obrigatório para NFS-e |
| `economicActivities[].type` | Tipo de atividade (principal/secundária) | 🟡 Média |
| `id_empresa` em `tb_endereco_empresa` | FK ligando endereço à empresa | 🔴 Alta — **a FK não existe na tabela!** |

> [!CAUTION]
> A tabela `tb_endereco_empresa` **não possui** coluna `id_empresa` nem `FOREIGN KEY` para `tb_empresas`. O endereço está completamente desvinculado da empresa.

---

## 2. Cliente (`cliente.sql` × Seção "Criar Venda > customer")

### ✅ Campos presentes

| Campo API | Coluna SQL | Observação |
|---|---|---|
| `name` | `nome` | ✅ |
| `federalTaxNumber` (CPF/CNPJ) | `cpf` + `cnpj` | ✅ campos separados |
| `stateTaxNumber` (IE) | `ie` | ✅ |
| `email` | `email` | ✅ |
| `phone` / `mobilePhone` | `telefone` em `tb_clientes` + `tb_telefones_clientes` | ✅ parcial |
| `isFinalCustomer` | `is_consumidor_final` | ✅ |
| `address.street` | `logradouro` | ✅ |
| `address.number` | `numero` | ✅ |
| `address.district` | `bairro` | ✅ |
| `address.postalCode` | `cep` | ✅ |
| `address.city.name` | `municipio` | ✅ |
| `address.city.code` | `codigo_municipio` | ✅ |
| `address.city.state` | `uf` | ✅ |

### ❌ Campos ausentes

| Campo API | Descrição | Urgência |
|---|---|---|
| `legalName` | Razão Social (para PJ) | 🔴 Alta — cliente PJ precisa de razão social separada do nome fantasia |
| `cityTaxNumber` | Inscrição Municipal | 🟡 Média |
| `address.additionalInformation` | Complemento do endereço | 🟡 Média |
| `address.country.name` / `code` | País | 🟢 Baixa — relevante só para exportação |
| `address.cityName` | Nome da cidade (campo alternativo da API) | 🟢 Baixa |

> [!NOTE]
> O `tipo` (PF/PJ) na tabela já cobre parcialmente o `legalName`, mas não armazena o valor da razão social separado do nome. Considere adicionar a coluna `razao_social` para clientes PJ.

---

## 3. Produto (`produtos.sql` × Seção "Criar Produto")

### ✅ Campos presentes

| Campo API | Coluna SQL | Observação |
|---|---|---|
| `code` (código interno) | `codigo_interno_loja` | ✅ |
| `name` | `nome` | ✅ |
| `price` (preço de venda) | `preco_venda` | ✅ |
| `productInvoiceSettings.gtinCode` | `codigo_barras` | ✅ |
| `productInvoiceSettings.ncm` | `ncm` em `tb_produtos_fiscal` | ✅ |
| `productInvoiceSettings.cest` | `cest` em `tb_produtos_fiscal` | ✅ |
| `productInvoiceSettings.cfopIn` | `cfop_padrao` em `tb_grupos_tributarios` | ✅ parcial |
| `productInvoiceSettings.unitOfMeasurement` | `unidade_venda` | ✅ |
| `productInvoiceSettings.taxUnitOfMeasurement` | `unidade_estoque` | ✅ parcial |
| ICMS (cst, alíquota, MVA, ST) | campos em `tb_grupos_tributarios` | ✅ |
| IPI (cst, alíquota) | `ipi_cst`, `ipi_aliquota` | ✅ |
| PIS / COFINS (cst, alíquota) | `pis_*`, `cofins_*` | ✅ |

### ❌ Campos ausentes

| Campo API | Descrição | Urgência |
|---|---|---|
| `invoiceModel` | Tipo de NF (NF-e de produto vs. NFS-e) | 🔴 Alta — define qual nota será emitida |
| `warrantyDays` | Prazo de garantia em dias | 🟢 Baixa |
| `serviceInvoiceSettings.*` | Configurações para NFS-e (código de serviço, CNAE, etc.) | 🟡 Média — necessário se vender serviços |
| `productInvoiceSettings.cfopInterstate` | CFOP para operação interestadual | 🔴 Alta — CFOP muda por tipo de operação |
| `productInvoiceSettings.cfopTaxpayerInterstate` | CFOP interestadual para contribuinte | 🔴 Alta |
| `productInvoiceSettings.cfopInternational` | CFOP para exportação | 🟢 Baixa |
| `productInvoiceSettings.exportationTaxUnitOfMeasurement` | Unidade tributável para exportação | 🟢 Baixa |
| `productInvoiceSettings.additionalInformation` | Informações adicionais do produto na NF | 🟡 Média |
| `productInvoiceSettings.fiscalDescription` | Descrição fiscal do produto | 🟡 Média |
| `productInvoiceSettings.series` / `nextNumber` | Série e próximo número da NF-e | 🟡 Média — normalmente gerenciado pela empresa, não pelo produto |
| `origem_mercadoria` para todos CFOPs | Origem está só no grupo, não por produto | 🟡 Média |

> [!IMPORTANT]
> Os CFOPs `cfopInterstate` e `cfopTaxpayerInterstate` são **obrigatórios** para emissão de NF-e. O `cfop_padrao` no grupo tributário provavelmente cobre só operações internas. Você precisará de colunas adicionais ou uma tabela de CFOP por tipo de operação.

---

## 4. Fornecedor (`fornecedor.sql` × Seção "Criar Nota > transport.carrier")

> O `notas.txt` não tem uma seção "Criar Fornecedor" explícita. O fornecedor aparece na NF como o **emitente** (quando é a empresa) ou como **transportadora** no campo `transport.carrier`. Usei os campos do carrier como referência para os dados esperados de um fornecedor/parceiro.

### ✅ Campos presentes

| Campo API | Coluna SQL | Observação |
|---|---|---|
| `name` / `legalName` | `razao_social` | ✅ |
| `federalTaxNumber` (CNPJ) | `cnpj` | ✅ |
| `stateTaxNumber` (IE) | `inscricao_estadual` | ✅ |
| `email` | `email` | ✅ |
| `phone` | `tb_telefones_fornecedores` | ✅ |
| `address.street` | `logradouro` | ✅ |
| `address.number` | `numero` | ✅ |
| `address.district` | `bairro` | ✅ |
| `address.postalCode` | `cep` | ✅ |
| `address.city.name` | `municipio` | ✅ |
| `address.city.code` | `codigo_municipio` | ✅ |
| `address.city.state` | `uf` | ✅ |

### ❌ Campos ausentes

| Campo API | Descrição | Urgência |
|---|---|---|
| `name` (nome fantasia) | Apenas razão social está mapeada; sem nome fantasia | 🟡 Média |
| `cityTaxNumber` | Inscrição Municipal | 🟡 Média |
| `mobilePhone` | Celular do fornecedor | 🟢 Baixa |
| `address.additionalInformation` | Complemento do endereço | 🟡 Média |
| `address.country.name` / `code` | País | 🟢 Baixa |

---

## Resumo Executivo

| Entidade | Status Geral | Principal Gap |
|---|---|---|
| **Empresa** | 🟡 Parcial | IE, e-mail, telefone, CNAE e **FK do endereço ausente** |
| **Cliente** | 🟢 Bom | `razao_social` para PJ e `complemento` do endereço |
| **Produto** | 🟡 Parcial | `invoiceModel` e múltiplos CFOPs por tipo de operação |
| **Fornecedor** | 🟢 Bom | `nome_fantasia`, `complemento` e `cityTaxNumber` |

### Prioridade de Correção

1. 🔴 **`tb_endereco_empresa`** — adicionar `id_empresa` (FK) imediatamente
2. 🔴 **`tb_empresas`** — adicionar `inscricao_estadual`, `email`, `telefone`
3. 🔴 **`tb_grupos_tributarios`** — adicionar `cfop_interestadual`, `cfop_interestadual_contribuinte`
4. 🔴 **`tb_produtos`** — adicionar `modelo_nota` (NF-e / NFS-e)
5. 🟡 **`tb_empresas`** — adicionar `inscricao_municipal`, CNAE (nova tabela)
6. 🟡 **`tb_clientes`** — adicionar `razao_social`, `complemento` no endereço
