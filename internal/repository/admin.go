package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type AdminRepository struct {
	db *sql.DB
}

func (r *AdminRepository) CarregarOrganizacoes(ctx context.Context) ([]model.Organizacao, error) {
	query := `
	select id, nome_fantasia, email, telefone, ''::text as celular, schema, ativo, criado_em, atualizado_em
	from public.tb_empresas_gestao
	order by nome_fantasia
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizacoes []model.Organizacao
	for rows.Next() {
		var o model.Organizacao
		if err := rows.Scan(
			&o.ID,
			&o.NomeFantasia,
			&o.Email,
			&o.Telefone,
			&o.Celular,
			&o.Schema,
			&o.Ativo,
			&o.Criado_Em,
			&o.Atualizado_Em,
		); err != nil {
			return nil, err
		}
		organizacoes = append(organizacoes, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return organizacoes, err
}

func (r *AdminRepository) CarregarUsuarios(ctx context.Context) ([]model.Usuario, error) {
	query := `
		select id, id_empresa, nome, cpf, telefone, email, ativo, criado_em, atualizado_em
		from public.tb_usuarios_gestao
		order by nome
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []model.Usuario
	for rows.Next() {
		var u model.Usuario
		if err := rows.Scan(
			&u.ID,
			&u.IDEmpresa,
			&u.Nome,
			&u.CPF,
			&u.Telefone,
			&u.Email,
			&u.Ativo,
			&u.CriadoEm,
			&u.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}
