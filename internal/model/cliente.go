package model

import (
	"errors"
)

type TipoPessoa string
type IndContribuinte int

const (
	PessoaFisica   TipoPessoa = "PF"
	PessoaJuridica TipoPessoa = "PJ"

	ContribuinteICMS   IndContribuinte = 1
	ContribuinteIsento IndContribuinte = 2
	NaoContribuinte    IndContribuinte = 9
)

type Cliente struct {
	ID                int64             `json:"id,omitempty"`
	Nome              string            `json:"nome,omitempty"`
	Tipo              TipoPessoa        `json:"tipo,omitempty"`
	Email             string            `json:"email,omitempty"`
	Telefone          string            `json:"telefone,omitempty"`
	CPF               string            `json:"cpf,omitempty"`
	CNPJ              string            `json:"cnpj,omitempty"`
	Contribuinte      IndContribuinte   `json:"contribuinte,omitempty"`
	IsConsumidorFinal bool              `json:"is_consumidor_final,omitempty"`
	IE                string            `json:"ie,omitempty"`
	Enderecos         []EnderecoCliente `json:"enderecos"`
	CreatedAt         string            `json:"created_at,omitempty"`
	UpdatedAt         string            `json:"updated_at,omitempty"`
}

type EnderecoCliente struct {
	ID              int64  `json:"id" db:"id"`
	IDCliente       int64  `json:"id_cliente"`
	CEP             string `json:"cep" db:"cep"`
	Logradouro      string `json:"logradouro" db:"logradouro"`
	Numero          string `json:"numero" db:"numero"`
	Bairro          string `json:"bairro" db:"bairro"`
	Municipio       string `json:"municipio" db:"municipio"`
	UF              string `json:"uf" db:"uf"`
	CodigoMunicipio string `json:"codigo_municipio" db:"codigo_municipio"`
	CreatedAt       string `json:"created_at" db:"created_at"`
}

type TelefoneCliente struct {
	ID        int64  `json:"id" db:"id"`
	IDCliente int64  `json:"id_cliente" db:"id_cliente"`
	DDD       string `json:"ddd" db:"ddd"`
	Numero    string `json:"numero" db:"numero"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

func (c *Cliente) Validar() error {
	var erros []error

	if c.Nome == "" {
		erros = append(erros, errors.New("nome do cliente obrigatório"))
	}
	if c.Tipo == "" {
		erros = append(erros, errors.New("tipo do cliente obrigatório"))
	}
	if c.Tipo == PessoaFisica && c.CPF == "" {
		erros = append(erros, errors.New("CPF obrigatório"))
	}
	if c.Tipo == PessoaJuridica && c.CNPJ == "" {
		erros = append(erros, errors.New("CNPJ obrigatório"))
	}
	if len(erros) != 0 {
		return errors.Join(erros...)
	}

	return nil
}

func (e *EnderecoCliente) Validar() error {
	var erros []error

	if e.CEP == "" {
		erros = append(erros, errors.New("CEP obrigatório"))
	}
	if e.Logradouro == "" {
		erros = append(erros, errors.New("logradouro obrigatório"))
	}
	if e.Numero == "" {
		erros = append(erros, errors.New("número obrigatório"))
	}
	if e.Bairro == "" {
		erros = append(erros, errors.New("bairro obrigatório"))
	}
	if e.Municipio == "" {
		erros = append(erros, errors.New("município obrigatório"))
	}
	if e.UF == "" {
		erros = append(erros, errors.New("UF obrigatório"))
	}
	if len(erros) != 0 {
		return errors.Join(erros...)
	}

	return nil
}

func (t *TelefoneCliente) Validar() error {
	var erros []error

	if t.DDD == "" {
		erros = append(erros, errors.New("DDD obrigatório"))
	}
	if t.Numero == "" {
		erros = append(erros, errors.New("número obrigatório"))
	}
	if len(erros) != 0 {
		return errors.Join(erros...)
	}

	return nil
}
