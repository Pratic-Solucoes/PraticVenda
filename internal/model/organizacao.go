package model

import "errors"

type Organizacao struct {
	ID            uint64 `json:"id,omitempty"`
	NomeFantasia  string `json:"nome_fantasia,omitempty"`
	Email         string `json:"email,omitempty"`
	Telefone      string `json:"telefone,omitempty"`
	Celular       string `json:"celular,omitempty"`
	Schema        string `json:"schema,omitempty"`
	Ativo         bool   `json:"ativo,omitempty"`
	Criado_Em     string `json:"criado_em,omitempty"`
	Atualizado_Em string `json:"atualizado_em,omitempty"`
}

func (o *Organizacao) Validar() error {
	var erros []error
	if o.NomeFantasia == "" {
		erros = append(erros, errors.New("Nome fantasia da organização é obrigatório."))
	}
	if o.Email == "" {
		erros = append(erros, errors.New("Email da organização é obrigatório."))
	}
	if o.Celular == "" {
		erros = append(erros, errors.New("Celular da organização é obrigatório."))
	}
	if o.Schema == "" {
		erros = append(erros, errors.New("Schema da organização é obrigatório."))
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}
	return nil
}
