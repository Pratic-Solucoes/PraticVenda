package organizacao

import "errors"

type Organizacao struct {
	ID              uint64 `json:"id,omitempty"`
	IDDono          uint64 `json:"id_dono,omitempty"`
	NomeOrganizacao string `json:"nome_organizacao,omitempty"`
	Ativo           bool   `json:"ativo,omitempty"`
	CriadoEm        string `json:"criado_em,omitempty"`
	AtualizadoEm    string `json:"atualizado_em,omitempty"`
}

func (o *Organizacao) Validar() error {
	var erros []error
	if o.IDDono == 0 {
		erros = append(erros, errors.New("o id_dono não foi informado"))
	}
	if o.NomeOrganizacao == "" {
		erros = append(erros, errors.New("o nome_organizacao não foi informado"))
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}
	return nil
}
