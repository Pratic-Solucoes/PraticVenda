package model

import (
	"errors"
	"time"
)

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

// ProdutoInput DTO para criação e edição de produtos
type ProdutoInput struct {
	CodigoBarras      *string               `json:"codigo_barras"`
	CodigoInternoLoja *string               `json:"codigo_interno_loja"`
	Nome              string                `json:"nome"`
	Descricao         *string               `json:"descricao"`
	PrecoCusto        float64               `json:"preco_custo"`
	PrecoVenda        float64               `json:"preco_venda"`
	UnidadeEstoque    string                `json:"unidade_estoque"`
	UnidadeVenda      string                `json:"unidade_venda"`
	PesoBruto         float64               `json:"peso_bruto"`
	PesoLiquido       float64               `json:"peso_liquido"`
	Ncm               string                `json:"ncm"`
	Cest              *string               `json:"cest"`
	IDGrupoTributario int64                 `json:"id_grupo_tributario"`
	Estoques          []ProdutoEstoqueInput `json:"estoques"`
}

func (p *ProdutoInput) Validar() error {
	if p.Nome == "" {
		return errors.New("o nome do produto é obrigatório")
	}
	if p.PrecoVenda < 0 {
		return errors.New("o preço de venda não pode ser negativo")
	}
	if p.PrecoCusto < 0 {
		return errors.New("o preço de custo não pode ser negativo")
	}
	if p.Ncm == "" {
		return errors.New("o NCM do produto é obrigatório")
	}
	if p.IDGrupoTributario <= 0 {
		return errors.New("o grupo tributário é obrigatório")
	}
	return nil
}
