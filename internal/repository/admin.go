package repository

import "database/sql"

type AdminRepository struct {
	db *sql.DB
}

func (r *AdminRepository) CarregarOrganizacoes()
