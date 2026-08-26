package model

type Organizacao struct {
	ID            uint64 `json:"id,omitempty"`
	NomeFantasia  string `json:"nome_fantasia,omitempty"`
	Email         string `json:"email,omitempty"`
	Telefone      string `json:"id,omitempty"`
	Celular       string `json:"id,omitempty"`
	Schema        string `json:"id,omitempty"`
	Ativo         bool   `json:"id,omitempty"`
	Criado_Em     string `json:"id,omitempty"`
	Atualizado_Em string `json:"id,omitempty"`
}
