package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gestao/internal/model"
)

type ContaPagarRepository struct {
	db *sql.DB
}

var (
	CONTA_PAGAR_QUITADA       = errors.New("conta a pagar já quitada")
	CONTA_PAGAR_NAO_ENCONTRADA = errors.New("conta a pagar não encontrada")
)

func (r *ContaPagarRepository) CriarContaPagar(ctx context.Context, tx *sql.Tx, contaPagar *model.ContaPagarCriar) error {
	query := `
		INSERT INTO tb_contas_pagar (
			id_fornecedor, id_categoria, descricao,
			nr_documento, nr_nota_fiscal, valor, dt_entrada, dt_vencimento,
			nr_parcela, nr_total_parcelas, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'PENDENTE')
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		contaPagar.IDFornecedor, contaPagar.IDCategoria, contaPagar.Descricao,
		contaPagar.NrDocumento, contaPagar.NrNotaFiscal, contaPagar.Valor, contaPagar.DtEntrada, contaPagar.DtVencimento,
		contaPagar.NrParcela, contaPagar.NrTotalParcelas,
	)
	return err
}

func (r *ContaPagarRepository) ListarContasPagar(ctx context.Context, tx *sql.Tx, busca, vencimento, status string) ([]*model.ContaPagar, error) {
	query := `
		SELECT d.id, d.id_fornecedor, d.id_categoria, d.descricao, d.nr_documento, d.nr_nota_fiscal,
		       d.valor, d.dt_entrada, d.dt_vencimento, d.nr_parcela, d.nr_total_parcelas, d.status, d.created_at, d.updated_at,
		       f.id, f.razao_social, f.cnpj,
		       c.id, c.nome
		FROM tb_contas_pagar d
		JOIN tb_fornecedores f ON d.id_fornecedor = f.id
		LEFT JOIN tb_categorias_contas_pagar c ON d.id_categoria = c.id
		WHERE 1=1
	`
	var args []interface{}

	if busca != "" {
		query += fmt.Sprintf(" AND (f.razao_social LIKE $%d OR f.id = $%d)", len(args)+1, len(args)+2)
		buscaParam := "%" + busca + "%"
		args = append(args, buscaParam, busca)
	}

	if vencimento != "" {
		query += fmt.Sprintf(" AND d.dt_vencimento = $%d", len(args)+1)
		args = append(args, vencimento)
	}

	if status != "" {
		query += fmt.Sprintf(" AND d.status = $%d", len(args)+1)
		args = append(args, status)
	} else {
		query += " AND d.status = 'PENDENTE'"
	}

	query += " ORDER BY d.dt_vencimento ASC"

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contasPagar []*model.ContaPagar
	for rows.Next() {
		d := &model.ContaPagar{}
		f := &model.Fornecedor{}

		var cId sql.NullInt64
		var cNome sql.NullString

		err := rows.Scan(
			&d.ID, &d.IDFornecedor, &d.IDCategoria, &d.Descricao, &d.NrDocumento, &d.NrNotaFiscal,
			&d.Valor, &d.DtEntrada, &d.DtVencimento, &d.NrParcela, &d.NrTotalParcelas, &d.Status, &d.CreatedAt, &d.UpdatedAt,
			&f.ID, &f.RazaoSocial, &f.CNPJ,
			&cId, &cNome,
		)
		if err != nil {
			return nil, err
		}

		d.Fornecedor = f
		if cId.Valid {
			d.Categoria = &model.CategoriaContaPagar{
				ID:   cId.Int64,
				Nome: cNome.String,
			}
		}
		contasPagar = append(contasPagar, d)
	}

	return contasPagar, nil
}

func (r *ContaPagarRepository) PagarContaPagar(ctx context.Context, tx *sql.Tx, id int64) error {

	query := `SELECT status from tb_contas_pagar WHERE id = $1`
	stmt, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		return err
	}

	var status string
	if stmt.Next() {
		if err := stmt.Scan(&status); err != nil {
			return err
		}
	}
	stmt.Close()

	if status != "PENDENTE" {
		return CONTA_PAGAR_QUITADA
	}

	query = "UPDATE tb_contas_pagar SET status = 'PAGO' WHERE id = $1"

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return CONTA_PAGAR_NAO_ENCONTRADA
	}
	return nil
}

func (r *ContaPagarRepository) EditarContaPagar(ctx context.Context, tx *sql.Tx, id int64, contaPagar *model.ContaPagarCriar) error {

	query := `SELECT status from tb_contas_pagar WHERE id = $1`
	stmt, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		return err
	}

	var status string
	if stmt.Next() {
		if err := stmt.Scan(&status); err != nil {
			return err
		}
	}
	stmt.Close()

	if status != "PENDENTE" {
		return CONTA_PAGAR_QUITADA
	}

	query = `
		UPDATE tb_contas_pagar SET
			id_fornecedor = $1, id_categoria = $2, descricao = $3,
			nr_documento = $4, nr_nota_fiscal = $5, valor = $6, dt_entrada = $7, dt_vencimento = $8,
			nr_parcela = $9, nr_total_parcelas = $10
		WHERE id = $11
	`
	result, err := tx.ExecContext(ctx, query,
		contaPagar.IDFornecedor, contaPagar.IDCategoria, contaPagar.Descricao,
		contaPagar.NrDocumento, contaPagar.NrNotaFiscal, contaPagar.Valor, contaPagar.DtEntrada, contaPagar.DtVencimento,
		contaPagar.NrParcela, contaPagar.NrTotalParcelas, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return CONTA_PAGAR_NAO_ENCONTRADA
	}
	return nil
}
