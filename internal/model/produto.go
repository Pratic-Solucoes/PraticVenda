package model

import (
	"errors"
	"strings"
	"time"
)

// Produto representa os campos básicos de tb_produtos necessários para o estoque
type Produto struct {
	ID                int64     `json:"id" db:"id"`
	Composto          bool      `json:"composto" db:"composto"`
	MateriaPrima      bool      `json:"materia_prima" db:"materia_prima"`
	IDFornecedor      int64     `json:"id_fornecedor" db:"id_fornecedor"`
	IDsFornecedores   []int64   `json:"ids_fornecedores,omitempty"`
	Fornecedor        string    `json:"fornecedor,omitempty" db:"fornecedor"`
	CodigoBarras      *string   `json:"codigo_barras,omitempty" db:"codigo_barras"`
	CodigoInternoLoja *string   `json:"codigo_interno_loja,omitempty" db:"codigo_interno_loja"`
	Nome              string    `json:"nome" db:"nome"`
	Descricao         *string   `json:"descricao,omitempty" db:"descricao"`
	PrecoCusto        float64   `json:"preco_custo" db:"preco_custo"`
	PrecoVenda        float64   `json:"preco_venda" db:"preco_venda"`
	UnidadeEstoque    string    `json:"unidade_estoque" db:"unidade_estoque"`
	UnidadeVenda      string    `json:"unidade_venda" db:"unidade_venda"`
	PesoBruto         float64   `json:"peso_bruto" db:"peso_bruto"`
	PesoLiquido       float64   `json:"peso_liquido" db:"peso_liquido"`
	Ativo             bool      `json:"ativo" db:"ativo"`
	CriadoEm          time.Time `json:"criado_em" db:"criado_em"`
	AtualizadoEm      time.Time `json:"atualizado_em" db:"atualizado_em"`
}

// GrupoTributario representa a tabela tb_grupos_tributarios
type GrupoTributario struct {
	ID               int64     `json:"id" db:"id"`
	Nome             string    `json:"nome" db:"nome"`
	CfopPadrao       string    `json:"cfop_padrao" db:"cfop_padrao"`
	OrigemMercadoria int       `json:"origem_mercadoria" db:"origem_mercadoria"`
	Csosn            *string   `json:"csosn,omitempty" db:"csosn"`
	IcmsCst          *string   `json:"icms_cst,omitempty" db:"icms_cst"`
	IcmsAliquota     float64   `json:"icms_aliquota" db:"icms_aliquota"`
	IcmsMvaSt        float64   `json:"icms_mva_st" db:"icms_mva_st"`
	IcmsAliquotaSt   float64   `json:"icms_aliquota_st" db:"icms_aliquota_st"`
	IpiCst           *string   `json:"ipi_cst,omitempty" db:"ipi_cst"`
	IpiAliquota      float64   `json:"ipi_aliquota" db:"ipi_aliquota"`
	PisCst           *string   `json:"pis_cst,omitempty" db:"pis_cst"`
	PisAliquota      float64   `json:"pis_aliquota" db:"pis_aliquota"`
	CofinsCst        *string   `json:"cofins_cst,omitempty" db:"cofins_cst"`
	CofinsAliquota   float64   `json:"cofins_aliquota" db:"cofins_aliquota"`
	CriadoEm         time.Time `json:"criado_em" db:"criado_em"`
	AtualizadoEm     time.Time `json:"atualizado_em" db:"atualizado_em"`
}

// ProdutoFiscal representa a tabela tb_produtos_fiscal
type ProdutoFiscal struct {
	IDProduto         int64     `json:"id_produto" db:"id_produto"`
	Ncm               string    `json:"ncm" db:"ncm"`
	Cest              *string   `json:"cest,omitempty" db:"cest"`
	IDGrupoTributario int64     `json:"id_grupo_tributario" db:"id_grupo_tributario"`
	AtualizadoEm      time.Time `json:"atualizado_em" db:"atualizado_em"`
}

// ProdutoEstoqueInput define o input de associação de estoque para o produto
type ProdutoEstoqueInput struct {
	IDEstoque     int64   `json:"id_estoque"`
	EstoqueMinimo float64 `json:"estoque_minimo"`
	Quantidade    float64 `json:"quantidade"`
}

// ProdutoCompleto é a entidade agregada para leitura/detalhes
type ProdutoCompleto struct {
	Produto
	Fiscal   *ProdutoFiscal       `json:"fiscal,omitempty"`
	Estoques []ProdutoEstoqueInfo `json:"estoques,omitempty"`
}

// ProdutoEstoqueInfo simplifica os detalhes do estoque associado
type ProdutoEstoqueInfo struct {
	IDEstoque     int64   `json:"id_estoque"`
	NomeEstoque   string  `json:"nome_estoque"`
	Quantidade    float64 `json:"quantidade"`
	EstoqueMinimo float64 `json:"estoque_minimo"`
}

// ItemComposicaoProduto é uma matéria-prima consumida na venda de um produto composto.
type ItemComposicaoProduto struct {
	IDProdutoComponente int64   `json:"id_produto_componente"`
	NomeProduto         string  `json:"nome_produto,omitempty"`
	UnidadeEstoque      string  `json:"unidade_estoque,omitempty"`
	Quantidade          float64 `json:"quantidade"`
	PrecoCusto          float64 `json:"preco_custo,omitempty"`
}

func (i ItemComposicaoProduto) Validar(idProdutoComposto int64) error {
	if i.IDProdutoComponente <= 0 || i.Quantidade <= 0 {
		return errors.New("produto componente e quantidade são obrigatórios")
	}
	if i.IDProdutoComponente == idProdutoComposto {
		return errors.New("um produto não pode fazer parte da própria composição")
	}
	return nil
}

// ProdutoInput DTO para criação e edição de produtos
type ProdutoInput struct {
	Composto          bool                    `json:"composto"`
	MateriaPrima      bool                    `json:"materia_prima"`
	IDFornecedor      int64                   `json:"id_fornecedor"`
	IDsFornecedores   []int64                 `json:"ids_fornecedores"`
	CodigoBarras      *string                 `json:"codigo_barras"`
	CodigoInternoLoja *string                 `json:"codigo_interno_loja"`
	Nome              string                  `json:"nome"`
	Descricao         *string                 `json:"descricao"`
	PrecoCusto        float64                 `json:"preco_custo"`
	PrecoVenda        float64                 `json:"preco_venda"`
	UnidadeEstoque    string                  `json:"unidade_estoque"`
	UnidadeVenda      string                  `json:"unidade_venda"`
	PesoBruto         float64                 `json:"peso_bruto"`
	PesoLiquido       float64                 `json:"peso_liquido"`
	Ncm               string                  `json:"ncm"`
	Cest              *string                 `json:"cest"`
	IDGrupoTributario int64                   `json:"id_grupo_tributario"`
	Estoques          []ProdutoEstoqueInput   `json:"estoques"`
	Composicao        []ItemComposicaoProduto `json:"composicao"`
}

func (p *ProdutoInput) Validar() error {
	p.Nome = strings.TrimSpace(p.Nome)
	p.UnidadeEstoque = strings.TrimSpace(p.UnidadeEstoque)
	p.UnidadeVenda = strings.TrimSpace(p.UnidadeVenda)

	if p.Nome == "" {
		return errors.New("o nome do produto é obrigatório")
	}
	if p.Composto && p.MateriaPrima {
		return errors.New("um produto não pode ser composto e matéria-prima ao mesmo tempo")
	}
	if !p.Composto && len(p.IDsFornecedores) == 0 && p.IDFornecedor <= 0 {
		return errors.New("vincule ao menos um fornecedor ao produto")
	}
	if p.UnidadeEstoque == "" {
		return errors.New("a unidade de estoque é obrigatória")
	}
	if p.UnidadeVenda == "" {
		return errors.New("a unidade de venda é obrigatória")
	}
	if !p.MateriaPrima && p.PrecoVenda <= 0 {
		return errors.New("o preço de venda deve ser maior que zero")
	}
	if p.PrecoCusto < 0 {
		return errors.New("o preço de custo não pode ser negativo")
	}
	if !p.Composto && len(p.Estoques) == 0 {
		return errors.New("vincule o produto a pelo menos um estoque")
	}
	if p.Composto && len(p.Composicao) == 0 {
		return errors.New("adicione ao menos uma matéria-prima ao produto composto")
	}
	return nil
}
