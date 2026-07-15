package repository

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
)

type EntradaEstoqueRepository struct {
	db *sql.DB
}

// RegistrarEntrada registra uma nova entrada de estoque.
// Parâmetros:
// - ctx: Contexto da requisição.
// - tx: Transação de banco de dados.
// - entrada: Dados da entrada de estoque.
// Retorno:
// - nil se a operação for bem-sucedida.
// - error se ocorrer algum erro.
func (r *EntradaEstoqueRepository) RegistrarEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error {

	stmt := `
		insert into tb_entradas_estoque(
			id_estoque,
			id_fornecedor,
			valor_despesa_adicional,
			id_usuario,
			valor_total,
			status
		)values(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6	
		)
	`

	var idEntrada int64
	err := tx.QueryRowContext(
		ctx,
		stmt+" RETURNING id",
		entrada.IDEstoque,
		entrada.IDFornecedor,
		entrada.ValorDespesaAdicional,
		entrada.IDUsuario,
		entrada.ValorTotal,
		entrada.Status,
	).Scan(&idEntrada)
	if err != nil {
		return err
	}

	if idEntrada == 0 {
		return errors.New("Nenhuma entrada de estoque criada.")
	}

	entrada.ID = uint64(idEntrada)
	return nil
}

// RegistrarProdutosEntrada registra os produtos da entrada de estoque.
// Parâmetros:
// - ctx: Contexto da requisição.
// - tx: Transação de banco de dados.
// - entrada: Dados da entrada de estoque.
// Retorno:
// - nil se a operação for bem-sucedida.
// - error se ocorrer algum erro.
func (r *EntradaEstoqueRepository) RegistrarProdutosEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error {

	// Registra os produtos na tabela de produtos da entrada de estoque
	query := `
		insert into	tb_produtos_entradas_estoque(
			id_entrada_estoque,
			id_produto,
			valor_unitario,
			valor_icms_st,
			valor_ipi,
			valor_desconto,
			rateio_despesa_adicional,
			valor_custo,
			valor_total,
			quantidade
		)
		values(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)
	`

	stmtRegistraProduto, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return errors.New("Erro ao preparar statement para inserção de produtos na tabela: tb_produtos_entrada_estoque. Erro: " + err.Error())
	}
	defer stmtRegistraProduto.Close()

	for _, produto := range entrada.Produtos {
		result, err := stmtRegistraProduto.ExecContext(
			ctx,
			entrada.ID,
			produto.IDProduto,
			produto.ValorUnitario,
			produto.ValorIcmsST,
			produto.ValorIPI,
			produto.ValorDesconto,
			produto.RateioDespesaAdicional,
			produto.ValorCusto,
			produto.ValorTotal,
			produto.Quantidade,
		)
		if err != nil {
			return errors.New("Erro ao inserir produto na tabela: tb_produtos_entrada_estoque. Erro: " + err.Error())
		}

		rowsAfetados, err := result.RowsAffected()
		if err != nil {
			return errors.New("Erro ao verificar inserção do produto na entrada de estoque. Erro: " + err.Error())
		}

		if rowsAfetados == 0 {
			return errors.New("Nenhum produto da entrada de estoque criado.")
		}
	}

	// Registra a movimentação do produto no estoque
	query = `
		insert into tb_movimento_estoque(
			id_produto,
			id_estoque,
			id_usuario,
			quantidade,
			tipo_movimento,
			id_categoria_movimento,
			id_origem
		)values(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
	`

	stmtRegistraMovimentoEstoque, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return errors.New("Erro ao preparar statement para inserção de movimentação do produto na tabela: tb_movimento_estoque. Erro: " + err.Error())
	}
	defer stmtRegistraMovimentoEstoque.Close()

	for _, produto := range entrada.Produtos {
		result, err := stmtRegistraMovimentoEstoque.ExecContext(
			ctx,
			produto.IDProduto,
			entrada.IDEstoque,
			entrada.IDUsuario,
			produto.Quantidade,
			"ENTRADA",
			1, // id_categoria_movimento, default "ENTRADA ESTOQUE" - DEVE SER O ID 1 DA TABELA tb_categoria_movimento_estoque
			entrada.ID,
		)
		if err != nil {
			return errors.New("Erro ao inserir movimentação do produto na tabela: tb_movimento_estoque. Erro: " + err.Error())
		}

		rowsAfetados, err := result.RowsAffected()
		if err != nil {
			return errors.New("Erro ao verificar inserção da movimentação do produto. Erro: " + err.Error())
		}

		if rowsAfetados == 0 {
			return errors.New("Nenhuma movimentação do produto criada.")
		}
	}

	// Atualiza o saldo do produto no estoque
	query = `
		update tb_produtos_estoque
		set quantidade = quantidade + $1,
			atualizado_em = now()
		where id_produto = $2
		and id_estoque = $3
	`
	stmtAtualizaSaldoProduto, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return errors.New("Erro ao preparar statement para atualização do saldo do produto na tabela: tb_produtos_estoque. Erro: " + err.Error())
	}
	defer stmtAtualizaSaldoProduto.Close()

	for _, produto := range entrada.Produtos {
		result, err := stmtAtualizaSaldoProduto.ExecContext(
			ctx,
			produto.Quantidade,
			produto.IDProduto,
			entrada.IDEstoque,
		)
		if err != nil {
			return errors.New("Erro ao atualizar saldo do produto na tabela: tb_produtos_estoque. Erro: " + err.Error())
		}

		rowsAfetados, err := result.RowsAffected()
		if err != nil {
			return errors.New("Erro ao verificar atualização do saldo do produto. Erro: " + err.Error())
		}

		if rowsAfetados == 0 {
			return errors.New("Produto não encontrado no estoque para atualização de saldo.")
		}
	}

	return nil
}
