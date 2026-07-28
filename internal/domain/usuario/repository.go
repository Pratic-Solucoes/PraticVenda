package usuario

import (
	"context"
	"database/sql"
)

type RepositoryInterface interface {
	CriarUsuario(ctx context.Context, tx *sql.Tx, usuario *Usuario) (*Usuario, error)
	BuscarUsuarioPorID(ctx context.Context, usuarioID uint64) (*Usuario, error)
	BuscarSenhaUsuario(ctx context.Context, usuarioID uint64) (*string, error)
	AtualizarSenhaUsuario(ctx context.Context, tx *sql.Tx, usuarioID uint64, novaSenha string) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RepositoryInterface {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CriarUsuario(ctx context.Context, tx *sql.Tx, usuario *Usuario) (*Usuario, error) {
	var id int64
	err := tx.QueryRowContext(ctx, `
		INSERT INTO tb_usuarios_gestao (nome, username, email, senha)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`, usuario.Nome, usuario.Username, usuario.Email, usuario.Senha).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &Usuario{
		ID:       id,
		Nome:     usuario.Nome,
		Username: usuario.Username,
		Email:    usuario.Email,
	}, nil
}

func (r *Repository) BuscarUsuarioPorID(ctx context.Context, usuarioID uint64) (*Usuario, error) {
	var usuario Usuario
	err := r.db.QueryRowContext(ctx, `select id, nome, username, email, celular from tb_usuarios_gestao where id = $1`, usuarioID).Scan(
		&usuario.ID, &usuario.Nome, &usuario.Username, &usuario.Email, &usuario.Celular)
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *Repository) BuscarSenhaUsuario(ctx context.Context, usuarioID uint64) (*string, error) {

	query := `select senha from tb_usuarios_gestao where id = $1`

	var senha string
	err := r.db.QueryRowContext(ctx, query, usuarioID).Scan(&senha)
	if err != nil {
		return nil, err
	}

	return &senha, nil
}

func (r *Repository) AtualizarSenhaUsuario(ctx context.Context, tx *sql.Tx, usuarioID uint64, novaSenha string) error {
	query := `UPDATE tb_usuarios_gestao SET senha = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, novaSenha, usuarioID)
	return err
}
