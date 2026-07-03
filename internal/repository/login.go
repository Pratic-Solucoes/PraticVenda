package repository

import (
	"context"
	"database/sql"
)

type LoginRepository struct {
	db *sql.DB
}

func (r *LoginRepository) Login(ctx context.Context, email string) (uint64, string, string, string, error) {
	var senhaDB string
	var nome string
	var id uint64
	var schema string

	err := r.db.QueryRowContext(ctx, `
		select u.id, u.nome, u.senha, e.schema
		from tb_usuarios_gestao u
		join tb_empresas_gestao e on e.id = u.id_empresa
		where u.email = $1 and e.ativo = true;
	`, email).Scan(&id, &nome, &senhaDB, &schema)

	if err != nil {
		return 0, "", "", "", err
	}

	return id, nome, senhaDB, schema, nil
}
