package model

import (
	"errors"
	"time"
)

// Fornecedor representa os dados básicos do fornecedor
type Fornecedor struct {
	ID                int64     `json:"id" db:"id"`
	RazaoSocial       string    `json:"razao_social" db:"razao_social"`
	CNPJ              string    `json:"cnpj" db:"cnpj"`
	InscricaoEstadual *string   `json:"inscricao_estadual,omitempty" db:"inscricao_estadual"`
	Email             *string   `json:"email,omitempty" db:"email"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`

	// Relacionamentos - Opcionais no momento do cadastro inicial
	Enderecos []EnderecoFornecedor `json:"enderecos,omitempty"`
	Telefones []TelefoneFornecedor `json:"telefones,omitempty"`
}

// EnderecoFornecedor representa o endereço vinculado a um fornecedor
type EnderecoFornecedor struct {
	ID              int64  `json:"id" db:"id"`
	IDFornecedor    int64  `json:"id_fornecedor" db:"id_fornecedor"`
	CEP             string `json:"cep" db:"cep"`
	Logradouro      string `json:"logradouro" db:"logradouro"`
	Numero          string `json:"numero" db:"numero"`
	Bairro          string `json:"bairro" db:"bairro"`
	Municipio       string `json:"municipio" db:"municipio"`
	UF              string `json:"uf" db:"uf"`
	CodigoMunicipio string `json:"codigo_municipio" db:"codigo_municipio"`
	IsPrincipal     bool   `json:"is_principal" db:"is_principal"`
	CreatedAt       string `json:"created_at" db:"created_at"`
}

// TelefoneFornecedor representa os telefones de contato de um fornecedor
type TelefoneFornecedor struct {
	ID           int64     `json:"id" db:"id"`
	IDFornecedor int64     `json:"id_fornecedor" db:"id_fornecedor"`
	DDD          string    `json:"ddd" db:"ddd"`
	Numero       string    `json:"numero" db:"numero"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Validar verifica se os dados obrigatórios do fornecedor estão presentes
func (f *Fornecedor) Validar() error {
	var erros []error
	if f.RazaoSocial == "" {
		erros = append(erros, errors.New("a razão social não foi informada"))
	}
	if f.CNPJ == "" {
		erros = append(erros, errors.New("o CNPJ não foi informado"))
	}
	if f.Email == nil || *f.Email == "" {
		erros = append(erros, errors.New("o email não foi informado"))
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}
	return nil
}
