package model

import "strings"

// CategoriaContaReceber classifica créditos e receitas.
type CategoriaContaReceber struct {
	ID   int64  `json:"id" db:"id"`
	Nome string `json:"nome" db:"nome"`
}

func (c *CategoriaContaReceber) Validar() bool {
	c.Nome = strings.TrimSpace(c.Nome)
	return c.Nome != ""
}
