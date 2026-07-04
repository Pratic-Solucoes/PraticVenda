package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type ClienteRepository struct {
	db *sql.DB
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullIfZeroInt(i model.IndContribuinte) interface{} {
	if i == 0 {
		return nil
	}
	return i
}

func (r *ClienteRepository) CriarCliente(ctx context.Context, tx *sql.Tx, c *model.Cliente) (*model.Cliente, error) {
	var id int64

	query := `
		INSERT INTO tb_clientes (nome, tipo, email, cpf, cnpj, contribuinte, is_consumidor_final, ie)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at;
	`
	err := tx.QueryRowContext(ctx, query, c.Nome, c.Tipo, nullIfEmpty(c.Email),
		nullIfEmpty(c.CPF), nullIfEmpty(c.CNPJ), nullIfZeroInt(c.Contribuinte), c.IsConsumidorFinal, nullIfEmpty(c.IE)).Scan(
		&id, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	c.ID = id

	var endID int64
	query = `
			insert into tb_enderecos_clientes(id_cliente, cep, logradouro, numero, bairro, municipio, uf, codigo_municipio)
			values ($1, $2, $3, $4, $5, $6, $7, $8)
			returning id, created_at;
		`
	if len(c.Enderecos) > 0 {
		for index, _ := range c.Enderecos {
			err = tx.QueryRowContext(
				ctx, query, c.ID, c.Enderecos[index].CEP, c.Enderecos[index].Logradouro,
				c.Enderecos[index].Numero, c.Enderecos[index].Bairro, c.Enderecos[index].Municipio, c.Enderecos[index].UF,
				c.Enderecos[index].CodigoMunicipio).Scan(
				&endID, &c.Enderecos[index].CreatedAt,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return c, nil
}

func (r *ClienteRepository) ListarClientes(ctx context.Context, tx *sql.Tx, busca string) ([]model.Cliente, error) {
	query := `
		SELECT id, nome, tipo, email, cpf, cnpj
		FROM tb_clientes
	`
	var rows *sql.Rows
	var err error

	if busca != "" {
		query += " WHERE nome ILIKE $1 OR cpf ILIKE $2 OR cnpj ILIKE $3"
		buscaParam := "%" + busca + "%"
		query += " ORDER BY id DESC"
		rows, err = tx.QueryContext(ctx, query, buscaParam, buscaParam, buscaParam)
	} else {
		query += " ORDER BY id DESC"
		rows, err = tx.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []model.Cliente
	for rows.Next() {
		var c model.Cliente
		var email, cpf, cnpj sql.NullString

		if err := rows.Scan(
			&c.ID, &c.Nome, &c.Tipo, &email, &cpf, &cnpj,
		); err != nil {
			return nil, err
		}

		c.Email = email.String
		c.CPF = cpf.String
		c.CNPJ = cnpj.String

		clientes = append(clientes, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clientes, nil
}

func (r *ClienteRepository) ObterClientePorID(ctx context.Context, tx *sql.Tx, id int64) (*model.Cliente, error) {
	query := `
		SELECT id, nome, tipo, email, cpf, cnpj, contribuinte, is_consumidor_final, ie, created_at, updated_at
		FROM tb_clientes
		WHERE id = $1
	`
	c := &model.Cliente{}
	var contribuinte sql.NullInt64
	var email, cpf, cnpj, ie sql.NullString

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Nome, &c.Tipo, &email, &cpf, &cnpj, &contribuinte, &c.IsConsumidorFinal, &ie, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	c.Email = email.String
	c.CPF = cpf.String
	c.CNPJ = cnpj.String
	c.IE = ie.String
	if contribuinte.Valid {
		c.Contribuinte = model.IndContribuinte(contribuinte.Int64)
	}

	// Buscar endereços
	queryEndereco := `
		SELECT id, id_cliente, cep, logradouro, numero, bairro, municipio, uf, codigo_municipio, created_at
		FROM tb_enderecos_clientes
		WHERE id_cliente = $1
	`
	rows, err := tx.QueryContext(ctx, queryEndereco, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	c.Enderecos = make([]model.EnderecoCliente, 0)
	for rows.Next() {
		var endereco model.EnderecoCliente
		err := rows.Scan(
			&endereco.ID, &endereco.IDCliente, &endereco.CEP, &endereco.Logradouro, &endereco.Numero, &endereco.Bairro,
			&endereco.Municipio, &endereco.UF, &endereco.CodigoMunicipio, &endereco.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		c.Enderecos = append(c.Enderecos, endereco)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Buscar telefones
	queryTelefone := `
		SELECT id, id_cliente, ddd, numero, created_at
		FROM tb_telefones_clientes
		WHERE id_cliente = $1
	`
	rowsTel, err := tx.QueryContext(ctx, queryTelefone, id)
	if err != nil {
		return nil, err
	}
	defer rowsTel.Close()

	c.Telefones = make([]model.TelefoneCliente, 0)
	for rowsTel.Next() {
		var telefone model.TelefoneCliente
		err := rowsTel.Scan(
			&telefone.ID, &telefone.IDCliente, &telefone.DDD, &telefone.Numero, &telefone.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		c.Telefones = append(c.Telefones, telefone)
	}
	if err = rowsTel.Err(); err != nil {
		return nil, err
	}

	return c, nil
}

func (r *ClienteRepository) AtualizarCliente(ctx context.Context, tx *sql.Tx, ID_Cliente int64, c *model.Cliente) error {
	// query para atualizar os dados do cliente
	query := `
		UPDATE tb_clientes
		SET nome = $1, tipo = $2, email = $3, cpf = $4, cnpj = $5, contribuinte = $6, is_consumidor_final = $7, ie = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
	`
	_, err := tx.ExecContext(ctx, query, c.Nome, c.Tipo, nullIfEmpty(c.Email),
		nullIfEmpty(c.CPF), nullIfEmpty(c.CNPJ), nullIfZeroInt(c.Contribuinte), c.IsConsumidorFinal, nullIfEmpty(c.IE), ID_Cliente)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClienteRepository) CriarEndereco(ctx context.Context, tx *sql.Tx, idCliente int64, e *model.EnderecoCliente) (*model.EnderecoCliente, error) {
	query := `
		INSERT INTO tb_enderecos_clientes (id_cliente, cep, logradouro, numero, bairro, municipio, uf, codigo_municipio)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at;
	`
	err := tx.QueryRowContext(ctx, query, idCliente, e.CEP, e.Logradouro, e.Numero, e.Bairro, e.Municipio, e.UF, e.CodigoMunicipio).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *ClienteRepository) EditarEndereco(ctx context.Context, tx *sql.Tx, idCliente int64, idEndereco int64, e *model.EnderecoCliente) error {
	query := `
		UPDATE tb_enderecos_clientes
		SET cep = $1, logradouro = $2, numero = $3, bairro = $4, municipio = $5, uf = $6, codigo_municipio = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id_cliente = $8 AND id = $9
	`

	_, err := tx.ExecContext(ctx, query, e.CEP, e.Logradouro, e.Numero, e.Bairro, e.Municipio, e.UF, e.CodigoMunicipio, idCliente, idEndereco)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClienteRepository) BuscarEnderecoByID(ctx context.Context, tx *sql.Tx, idCliente int64, idEndereco int64) (*model.EnderecoCliente, error) {
	query := `
		SELECT id, id_cliente, cep, logradouro, numero, bairro, municipio, uf, codigo_municipio, created_at
		FROM tb_enderecos_clientes
		WHERE id_cliente = $1 AND id = $2
	`
	e := &model.EnderecoCliente{}
	err := tx.QueryRowContext(ctx, query, idCliente, idEndereco).Scan(
		&e.ID, &e.IDCliente, &e.CEP, &e.Logradouro, &e.Numero, &e.Bairro, &e.Municipio, &e.UF, &e.CodigoMunicipio, &e.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return e, nil
}

func (r *ClienteRepository) CriarTelefone(ctx context.Context, tx *sql.Tx, idCliente int64, t *model.TelefoneCliente) (*model.TelefoneCliente, error) {
	query := `
		INSERT INTO tb_telefones_clientes (id_cliente, ddd, numero)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`
	err := tx.QueryRowContext(ctx, query, idCliente, t.DDD, t.Numero).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *ClienteRepository) EditarTelefone(ctx context.Context, tx *sql.Tx, idCliente int64, idTelefone int64, t *model.TelefoneCliente) error {
	query := `
		UPDATE tb_telefones_clientes
		SET ddd = $1, numero = $2
		WHERE id_cliente = $3 AND id = $4
	`
	_, err := tx.ExecContext(ctx, query, t.DDD, t.Numero, idCliente, idTelefone)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClienteRepository) BuscarTelefoneByID(ctx context.Context, tx *sql.Tx, idCliente int64, idTelefone int64) (*model.TelefoneCliente, error) {
	query := `
		SELECT id, id_cliente, ddd, numero, created_at
		FROM tb_telefones_clientes
		WHERE id_cliente = $1 AND id = $2
	`
	t := &model.TelefoneCliente{}
	err := tx.QueryRowContext(ctx, query, idCliente, idTelefone).Scan(
		&t.ID, &t.IDCliente, &t.DDD, &t.Numero, &t.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}
