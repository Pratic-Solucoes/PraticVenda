package repository

import (
	"context"
	"database/sql"
)

type LoginRepository struct {
	db *sql.DB
}

func (r *LoginRepository) Login(ctx context.Context, username string) (uint64, string, string, error) {
	var senhaDB string
	var nome string
	var id uint64

	err := r.db.QueryRowContext(ctx, `
		select id, nome, senha
		from tb_usuarios_gestao
		where username = $1 and ativo = true;
	`, username).Scan(&id, &nome, &senhaDB)

	if err != nil {
		return 0, "", "", err
	}

	return id, nome, senhaDB, nil
}
