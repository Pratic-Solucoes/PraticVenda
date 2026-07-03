package model

import "errors"

type Empresa struct {
	ID           int64  `json:"id,omitempty"`
	RazaoSocial  string `json:"razao_social,omitempty"`
	NomeFantasia string `json:"nome_fantasia,omitempty"`
	CNPJ         string `json:"cnpj,omitempty"`
	Ativo        bool   `json:"ativo,omitempty"`
	CriadoEm     string `json:"criado_em,omitempty"`
	AtualizadoEm string `json:"atualizado_em,omitempty"`
}

type CredenciaisEmpresa struct {
	ID                 int64
	TipoAmbiente       int8
	CertificadoDigital string
	SenhaCertificado   string
	IDCsc              string
	CscNfe             string
}

type DadosFiscaisEmpresa struct {
	ID                     int64
	CodigoRegimeTributario int8
}

type EnderecoEmpresa struct {
	ID           int64
	Logradouro   string
	Numero       int64
	Bairro       string
	Cep          int64
	CodigoCidade int64
	NomeCidade   string
	Estado       string
	Pais         string
	CodigoPais   int64
	CriadoEm     string
	AtualizadoEm string
}

func (e *Empresa) Validar() error {
	var erros []error
	if e.RazaoSocial == "" {
		erros = append(erros, errors.New("Necessário informar razão social"))
	}
	if e.NomeFantasia == "" {
		erros = append(erros, errors.New("Necessário informar nome fantasia"))
	}
	if e.CNPJ == "" {
		erros = append(erros, errors.New("Necessário informar CNPJ"))
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}

	return nil

}
