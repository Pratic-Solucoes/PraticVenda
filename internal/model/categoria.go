package model

// CategoriaContaPagar representa a categoria da conta a pagar (baseado no SQL tb_categorias_contas_pagar)
type CategoriaContaPagar struct {
	ID   int64  `json:"id" db:"id"`
	Nome string `json:"nome" db:"nome"`
}
