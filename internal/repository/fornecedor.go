package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type FornecedorRepository struct {
	db *sql.DB
}

func (r *FornecedorRepository) ListarFornecedores(ctx context.Context, tx *sql.Tx, busca string) ([]*model.Fornecedor, error) {
	query := `
		SELECT id, razao_social, cnpj, inscricao_estadual, email, created_at, updated_at
		FROM tb_fornecedores
	`
	var rows *sql.Rows
	var err error

	if busca != "" {
		query += " WHERE razao_social LIKE $1 OR cnpj LIKE $2"
		buscaParam := "%" + busca + "%"
		rows, err = tx.QueryContext(ctx, query, buscaParam, buscaParam)
	} else {
		rows, err = tx.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fornecedores []*model.Fornecedor
	for rows.Next() {
		f := &model.Fornecedor{}
		err := rows.Scan(&f.ID, &f.RazaoSocial, &f.CNPJ, &f.InscricaoEstadual, &f.Email, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, err
		}
		fornecedores = append(fornecedores, f)
	}
	return fornecedores, nil
}

// CriarFornecedor insere um novo fornecedor e seus relacionamentos (se existirem) no banco de dados.
// Recebe um tx *sql.Tx pois a inserção em múltiplas tabelas (fornecedores, endereços, telefones) deve ser transacional.
func (r *FornecedorRepository) CriarFornecedor(ctx context.Context, tx *sql.Tx, f *model.Fornecedor) (*model.Fornecedor, error) {
	var id int64

	// 1. Inserir Fornecedor
	queryFornecedor := `
		INSERT INTO tb_fornecedores (razao_social, cnpj, inscricao_estadual, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at;
	`
	err := tx.QueryRowContext(ctx, queryFornecedor, f.RazaoSocial, f.CNPJ, f.InscricaoEstadual, f.Email).Scan(&id, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	f.ID = id

	// 2. Inserir Endereços (Opcional)
	for i, end := range f.Enderecos {
		var endID int64
		queryEndereco := `
			INSERT INTO tb_enderecos_fornecedores (id_fornecedor, cep, logradouro, numero, bairro, municipio, uf, codigo_municipio, is_principal)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, created_at;
		`
		err = tx.QueryRowContext(ctx, queryEndereco,
			f.ID, end.CEP, end.Logradouro, end.Numero, end.Bairro, end.Municipio, end.UF, end.CodigoMunicipio, end.IsPrincipal,
		).Scan(&endID, &end.CreatedAt)

		if err != nil {
			return nil, err
		}

		// Atualiza os dados do item no slice
		f.Enderecos[i].ID = endID
		f.Enderecos[i].IDFornecedor = f.ID
		f.Enderecos[i].CreatedAt = end.CreatedAt
	}

	// 3. Inserir Telefones (Opcional)
	for i, tel := range f.Telefones {
		var telID int64
		queryTelefone := `
			INSERT INTO tb_telefones_fornecedores (id_fornecedor, ddd, numero)
			VALUES ($1, $2, $3)
			RETURNING id, created_at;
		`
		err = tx.QueryRowContext(ctx, queryTelefone, f.ID, tel.DDD, tel.Numero).Scan(&telID, &tel.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Atualiza os dados do item no slice
		f.Telefones[i].ID = telID
		f.Telefones[i].IDFornecedor = f.ID
		f.Telefones[i].CreatedAt = tel.CreatedAt
	}

	return f, nil
}

func (r *FornecedorRepository) ObterFornecedorPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.Fornecedor, error) {
	queryFornecedor := `
		SELECT id, razao_social, cnpj, inscricao_estadual, email, created_at, updated_at
		FROM tb_fornecedores
		WHERE id = $1
	`
	f := &model.Fornecedor{}
	err := tx.QueryRowContext(ctx, queryFornecedor, id).Scan(
		&f.ID, &f.RazaoSocial, &f.CNPJ, &f.InscricaoEstadual, &f.Email, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	queryEnderecos := `
		SELECT id, id_fornecedor, cep, logradouro, numero, bairro, municipio, uf, codigo_municipio, is_principal, created_at
		FROM tb_enderecos_fornecedores
		WHERE id_fornecedor = $1
	`
	rowsEnd, err := tx.QueryContext(ctx, queryEnderecos, id)
	if err == nil {
		defer rowsEnd.Close()
		for rowsEnd.Next() {
			end := model.EnderecoFornecedor{}
			err := rowsEnd.Scan(
				&end.ID, &end.IDFornecedor, &end.CEP, &end.Logradouro, &end.Numero, &end.Bairro,
				&end.Municipio, &end.UF, &end.CodigoMunicipio, &end.IsPrincipal, &end.CreatedAt,
			)
			if err == nil {
				f.Enderecos = append(f.Enderecos, end)
			}
		}
	}

	queryTelefones := `
		SELECT id, id_fornecedor, ddd, numero, created_at
		FROM tb_telefones_fornecedores
		WHERE id_fornecedor = $1
	`
	rowsTel, err := tx.QueryContext(ctx, queryTelefones, id)
	if err == nil {
		defer rowsTel.Close()
		for rowsTel.Next() {
			tel := model.TelefoneFornecedor{}
			err := rowsTel.Scan(&tel.ID, &tel.IDFornecedor, &tel.DDD, &tel.Numero, &tel.CreatedAt)
			if err == nil {
				f.Telefones = append(f.Telefones, tel)
			}
		}
	}

	return f, nil
}

func (r *FornecedorRepository) AtualizarFornecedor(ctx context.Context, tx *sql.Tx, id int64, f *model.Fornecedor) error {
	queryFornecedor := `
		UPDATE tb_fornecedores
		SET razao_social = $1, cnpj = $2, inscricao_estadual = $3, email = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	_, err := tx.ExecContext(ctx, queryFornecedor, f.RazaoSocial, f.CNPJ, f.InscricaoEstadual, f.Email, id)
	if err != nil {
		return err
	}

	// Deleta endereços e telefones atuais
	_, err = tx.ExecContext(ctx, "DELETE FROM tb_enderecos_fornecedores WHERE id_fornecedor = $1", id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM tb_telefones_fornecedores WHERE id_fornecedor = $1", id)
	if err != nil {
		return err
	}

	// Insere novos endereços
	for _, end := range f.Enderecos {
		queryEndereco := `
			INSERT INTO tb_enderecos_fornecedores (id_fornecedor, cep, logradouro, numero, bairro, municipio, uf, codigo_municipio, is_principal)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
		_, err = tx.ExecContext(ctx, queryEndereco,
			id, end.CEP, end.Logradouro, end.Numero, end.Bairro, end.Municipio, end.UF, end.CodigoMunicipio, end.IsPrincipal,
		)
		if err != nil {
			return err
		}
	}

	// Insere novos telefones
	for _, tel := range f.Telefones {
		queryTelefone := `
			INSERT INTO tb_telefones_fornecedores (id_fornecedor, ddd, numero)
			VALUES ($1, $2, $3)
		`
		_, err = tx.ExecContext(ctx, queryTelefone, id, tel.DDD, tel.Numero)
		if err != nil {
			return err
		}
	}

	return nil
}
