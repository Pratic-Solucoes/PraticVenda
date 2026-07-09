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
			id_fornecedor, id_categoria, id_grupo_parcelas, descricao,
			nr_documento, nr_nota_fiscal, valor_original, saldo_restante, dt_entrada, dt_vencimento,
			nr_parcela, nr_total_parcelas, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, 'PENDENTE')
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// REFATORAÇÃO: O 'saldo_restante' inicializa com o mesmo valor do 'valor_original'
	_, err = stmt.ExecContext(ctx,
		contaPagar.IDFornecedor, contaPagar.IDCategoria, contaPagar.IDGrupoParcelas, contaPagar.Descricao,
		contaPagar.NrDocumento, contaPagar.NrNotaFiscal, contaPagar.ValorOriginal, contaPagar.ValorOriginal,
		contaPagar.DtEntrada, contaPagar.DtVencimento, contaPagar.NrParcela, contaPagar.NrTotalParcelas,
	)
	return err
}

func (r *ContaPagarRepository) ListarContasPagar(ctx context.Context, tx *sql.Tx, busca, vencimento, status string) ([]*model.ContaPagar, error) {
	query := `
		SELECT d.id, d.id_fornecedor, d.id_categoria, d.id_grupo_parcelas, d.descricao, d.nr_documento, d.nr_nota_fiscal,
		       d.valor_original, d.saldo_restante, d.dt_entrada, d.dt_vencimento, d.nr_parcela, d.nr_total_parcelas, d.status, d.created_at, d.dt_pagamento, d.updated_at,
		       f.id, f.razao_social, f.cnpj,
		       c.id, c.nome
		FROM tb_contas_pagar d
		JOIN tb_fornecedores f ON d.id_fornecedor = f.id
		JOIN tb_categorias_contas_pagar c ON d.id_categoria = c.id
		WHERE 1=1
	`
	var args []interface{}

	if busca != "" {
		query += fmt.Sprintf(" AND (f.razao_social LIKE $%d OR d.id::text = $%d)", len(args)+1, len(args)+2)
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
		c := &model.CategoriaContaPagar{}

		err := rows.Scan(
			&d.ID, &d.IDFornecedor, &d.IDCategoria, &d.IDGrupoParcelas, &d.Descricao, &d.NrDocumento, &d.NrNotaFiscal,
			&d.ValorOriginal, &d.SaldoRestante, &d.DtEntrada, &d.DtVencimento, &d.NrParcela, &d.NrTotalParcelas, &d.Status, &d.CreatedAt, &d.DtPagamento, &d.UpdatedAt,
			&f.ID, &f.RazaoSocial, &f.CNPJ,
			&c.ID, &c.Nome,
		)
		if err != nil {
			return nil, err
		}

		d.Fornecedor = f
		d.Categoria = c
		contasPagar = append(contasPagar, d)
	}

	return contasPagar, nil
}

// REFATORAÇÃO: O método PagarContaPagar original pagava a conta inteira sem registrar transação.
// Agora, ele vai realizar o pagamento integral baixando o saldo_restante para 0 e criando uma entrada na tabela tb_movimento_financeiro.
// Obs: No futuro, você pode criar uma função 'PagarContaParcial(ctx, tx, id, valor_pagamento)' 
//      e reutilizar essa estrutura atualizando saldo_restante -= valor_pagamento.
func (r *ContaPagarRepository) PagarContaPagar(ctx context.Context, tx *sql.Tx, id int64) error {

	// 1. Pega os dados atuais da conta
	queryConta := `SELECT status, saldo_restante FROM tb_contas_pagar WHERE id = $1`
	var status string
	var saldoRestante float64
	
	err := tx.QueryRowContext(ctx, queryConta, id).Scan(&status, &saldoRestante)
	if err != nil {
		if err == sql.ErrNoRows {
			return CONTA_PAGAR_NAO_ENCONTRADA
		}
		return err
	}

	if status == "PAGO" {
		return CONTA_PAGAR_QUITADA
	}

	// 2. Insere na tabela de fluxo de caixa (tb_movimento_financeiro)
	// Vamos supor que ele está pagando o valor integral do saldo restante
	queryMovimento := `
		INSERT INTO tb_movimento_financeiro (
			tipo_movimento, id_conta_pagar, dt_movimento, valor_movimento
		) VALUES (
			'SAIDA', $1, CURRENT_DATE, $2
		)
	`
	_, err = tx.ExecContext(ctx, queryMovimento, id, saldoRestante)
	if err != nil {
		return fmt.Errorf("erro ao registrar movimento financeiro: %w", err)
	}

	// 3. Atualiza a conta a pagar
	queryUpdateConta := `
		UPDATE tb_contas_pagar 
		SET status = 'PAGO', saldo_restante = 0, dt_pagamento = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, queryUpdateConta, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar conta a pagar: %w", err)
	}

	return nil
}

func (r *ContaPagarRepository) EditarContaPagar(ctx context.Context, tx *sql.Tx, id int64, contaPagar *model.ContaPagarCriar) error {

	query := `SELECT status from tb_contas_pagar WHERE id = $1`
	var status string
	if err := tx.QueryRowContext(ctx, query, id).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return CONTA_PAGAR_NAO_ENCONTRADA
		}
		return err
	}

	// Permite editar se a conta estiver pendente ou atrasada. Se estiver paga, não faz sentido editar os valores originais livremente.
	if status == "PAGO" {
		return CONTA_PAGAR_QUITADA
	}

	// REFATORAÇÃO: Ao editar o 'valor_original', precisamos decidir o que acontece com o 'saldo_restante'.
	// Para simplificar, assumimos aqui que se está editando uma conta pendente, o saldo_restante iguala ao novo valor original.
	queryUpdate := `
		UPDATE tb_contas_pagar SET
			id_fornecedor = $1, id_categoria = $2, id_grupo_parcelas = $3, descricao = $4,
			nr_documento = $5, nr_nota_fiscal = $6, valor_original = $7, saldo_restante = $8, dt_entrada = $9, dt_vencimento = $10,
			nr_parcela = $11, nr_total_parcelas = $12
		WHERE id = $13
	`
	_, err := tx.ExecContext(ctx, queryUpdate,
		contaPagar.IDFornecedor, contaPagar.IDCategoria, contaPagar.IDGrupoParcelas, contaPagar.Descricao,
		contaPagar.NrDocumento, contaPagar.NrNotaFiscal, contaPagar.ValorOriginal, contaPagar.ValorOriginal, contaPagar.DtEntrada, contaPagar.DtVencimento,
		contaPagar.NrParcela, contaPagar.NrTotalParcelas, id)
	if err != nil {
		return err
	}

	return nil
}
