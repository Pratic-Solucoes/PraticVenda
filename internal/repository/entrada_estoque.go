package repository

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"strconv"
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

// SalvarItensEntrada grava os itens de uma entrada aberta sem movimentar saldo.
func (r *EntradaEstoqueRepository) SalvarItensEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error {
	const query = `INSERT INTO tb_produtos_entradas_estoque
		(id_entrada_estoque, id_produto, valor_unitario, valor_icms_st, valor_ipi, valor_desconto, rateio_despesa_adicional, valor_custo, valor_total, quantidade)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	for _, produto := range entrada.Produtos {
		if _, err := tx.ExecContext(ctx, query, entrada.ID, produto.IDProduto, produto.ValorUnitario, produto.ValorIcmsST,
			produto.ValorIPI, produto.ValorDesconto, produto.RateioDespesaAdicional, produto.ValorCusto, produto.ValorTotal, produto.Quantidade); err != nil {
			return err
		}
	}
	return nil
}

func (r *EntradaEstoqueRepository) ValidarProdutosFornecedor(ctx context.Context, tx *sql.Tx, idFornecedor uint64, produtos []model.ProdutoEntradaEstoque) error {
	for _, produto := range produtos {
		var encontrado bool
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tb_produtos p JOIN tb_produtos_fornecedores pf ON pf.id_produto=p.id WHERE p.id=$1 AND pf.id_fornecedor=$2 AND p.ativo=TRUE)`, produto.IDProduto, idFornecedor).Scan(&encontrado); err != nil {
			return err
		}
		if !encontrado {
			return errors.New("todos os produtos da entrada devem pertencer ao fornecedor selecionado")
		}
	}
	return nil
}

// ObterEntrada busca os dados completos, incluindo os itens, de uma entrada.
func (r *EntradaEstoqueRepository) ObterEntrada(ctx context.Context, tx *sql.Tx, id uint64) (*model.EntradaEstoque, error) {
	entrada := &model.EntradaEstoque{}
	err := tx.QueryRowContext(ctx, `SELECT ee.id, ee.id_estoque, e.nome, ee.id_fornecedor, f.razao_social,
		ee.valor_despesa_adicional, ee.id_usuario, ee.valor_total, ee.status, ee.criado_em, u.nome
		FROM tb_entradas_estoque ee
		JOIN tb_estoques e ON e.id = ee.id_estoque JOIN tb_fornecedores f ON f.id = ee.id_fornecedor
		JOIN public.tb_usuarios_gestao u ON u.id = ee.id_usuario WHERE ee.id = $1`, id).Scan(
		&entrada.ID, &entrada.IDEstoque, &entrada.Estoque, &entrada.IDFornecedor, &entrada.Fornecedor,
		&entrada.ValorDespesaAdicional, &entrada.IDUsuario, &entrada.ValorTotal, &entrada.Status, &entrada.CriadoEm, &entrada.Usuario)
	if err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, `SELECT i.id, i.id_produto, p.nome, i.valor_unitario, i.valor_icms_st, i.valor_ipi, i.valor_desconto,
		i.rateio_despesa_adicional, i.valor_custo, i.valor_total, i.quantidade FROM tb_produtos_entradas_estoque i
		JOIN tb_produtos p ON p.id = i.id_produto WHERE i.id_entrada_estoque = $1 ORDER BY i.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var produto model.ProdutoEntradaEstoque
		if err := rows.Scan(&produto.ID, &produto.IDProduto, &produto.NomeProduto, &produto.ValorUnitario, &produto.ValorIcmsST, &produto.ValorIPI,
			&produto.ValorDesconto, &produto.RateioDespesaAdicional, &produto.ValorCusto, &produto.ValorTotal, &produto.Quantidade); err != nil {
			return nil, err
		}
		entrada.Produtos = append(entrada.Produtos, produto)
	}
	return entrada, rows.Err()
}

// AtualizarEntrada substitui os dados e itens somente se a entrada estiver aberta.
func (r *EntradaEstoqueRepository) AtualizarEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error {
	result, err := tx.ExecContext(ctx, `UPDATE tb_entradas_estoque SET id_estoque=$1, id_fornecedor=$2,
		valor_despesa_adicional=$3, valor_total=$4 WHERE id=$5 AND status='ABERTO'`, entrada.IDEstoque, entrada.IDFornecedor,
		entrada.ValorDespesaAdicional, entrada.ValorTotal, entrada.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("somente entradas em aberto podem ser editadas")
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM tb_produtos_entradas_estoque WHERE id_entrada_estoque=$1`, entrada.ID); err != nil {
		return err
	}
	return r.SalvarItensEntrada(ctx, tx, entrada)
}

// AprovarEntrada muda o estado para concluída e aplica os itens ao saldo do estoque.
func (r *EntradaEstoqueRepository) AprovarEntrada(ctx context.Context, tx *sql.Tx, id uint64, usuarioID int64) error {
	result, err := tx.ExecContext(ctx, `UPDATE tb_entradas_estoque SET status='CONCLUIDA' WHERE id=$1 AND status='ABERTO'`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("somente entradas em aberto podem ser aprovadas")
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO tb_movimento_estoque (id_produto,id_estoque,id_usuario,quantidade,tipo_movimento,id_categoria_movimento,id_origem)
		SELECT i.id_produto,e.id_estoque,$2,i.quantidade,'ENTRADA',c.id,e.id
		FROM tb_produtos_entradas_estoque i JOIN tb_entradas_estoque e ON e.id=i.id_entrada_estoque
		JOIN tb_categoria_movimento_estoque c ON c.nome='ENTRADA DE ESTOQUE' WHERE e.id=$1`, id, usuarioID); err != nil {
		return err
	}
	result, err = tx.ExecContext(ctx, `UPDATE tb_produtos_estoque pe SET quantidade=pe.quantidade+itens.quantidade, atualizado_em=NOW()
		FROM (SELECT i.id_produto,e.id_estoque,SUM(i.quantidade) quantidade FROM tb_produtos_entradas_estoque i
		JOIN tb_entradas_estoque e ON e.id=i.id_entrada_estoque WHERE e.id=$1 GROUP BY i.id_produto,e.id_estoque) itens
		WHERE pe.id_produto=itens.id_produto AND pe.id_estoque=itens.id_estoque`, id)
	if err != nil {
		return err
	}
	rows, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("nenhum produto da entrada está vinculado ao estoque")
	}
	if _, err = tx.ExecContext(ctx, `UPDATE tb_produtos p SET preco_custo=itens.valor_custo, atualizado_em=NOW()
		FROM (SELECT DISTINCT ON (id_produto) id_produto,valor_custo FROM tb_produtos_entradas_estoque WHERE id_entrada_estoque=$1 ORDER BY id_produto,id DESC) itens
		WHERE p.id=itens.id_produto`, id); err != nil {
		return err
	}
	return nil
}

func (r *EntradaEstoqueRepository) CancelarEntrada(ctx context.Context, tx *sql.Tx, id uint64, usuarioID int64) error {
	var estoque uint64
	var status string
	if err := tx.QueryRowContext(ctx, `SELECT id_estoque,status FROM tb_entradas_estoque WHERE id=$1 FOR UPDATE`, id).Scan(&estoque, &status); err != nil {
		return err
	}
	if status == "CANCELADA" {
		return errors.New("entrada já está cancelada")
	}
	if status == "ABERTO" {
		_, err := tx.ExecContext(ctx, `UPDATE tb_entradas_estoque SET status='CANCELADA' WHERE id=$1`, id)
		return err
	}
	rows, err := tx.QueryContext(ctx, `SELECT id_produto,SUM(quantidade) FROM tb_produtos_entradas_estoque WHERE id_entrada_estoque=$1 GROUP BY id_produto`, id)
	if err != nil {
		return err
	}
	type itemCancelamento struct {
		produto uint64
		qtd     float64
	}
	itens := make([]itemCancelamento, 0)
	for rows.Next() {
		var produto uint64
		var qtd float64
		if err = rows.Scan(&produto, &qtd); err != nil {
			rows.Close()
			return err
		}
		itens = append(itens, itemCancelamento{produto: produto, qtd: qtd})
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err = rows.Close(); err != nil {
		return err
	}
	for _, item := range itens {
		produto, qtd := item.produto, item.qtd
		var saldo float64
		if err = tx.QueryRowContext(ctx, `SELECT quantidade FROM tb_produtos_estoque WHERE id_produto=$1 AND id_estoque=$2 FOR UPDATE`, produto, estoque).Scan(&saldo); err != nil {
			return err
		}
		if saldo < qtd {
			return errors.New("não é possível cancelar a entrada: parte do saldo já foi consumida")
		}
		if _, err = tx.ExecContext(ctx, `UPDATE tb_produtos_estoque SET quantidade=quantidade-$1,atualizado_em=NOW() WHERE id_produto=$2 AND id_estoque=$3`, qtd, produto, estoque); err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO tb_movimento_estoque(id_produto,id_estoque,id_usuario,quantidade,tipo_movimento,id_categoria_movimento,id_origem) SELECT $1,$2,$3,$4,'SAIDA',id,$5 FROM tb_categoria_movimento_estoque WHERE nome='CANCELAMENTO DE ENTRADA'`, produto, estoque, usuarioID, qtd, id); err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `UPDATE tb_entradas_estoque SET status='CANCELADA' WHERE id=$1`, id)
	return err
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
	for _, produto := range entrada.Produtos {
		if _, err := tx.ExecContext(ctx, `UPDATE tb_produtos SET preco_custo=$1, atualizado_em=NOW() WHERE id=$2`, produto.ValorCusto, produto.IDProduto); err != nil {
			return err
		}
	}

	return nil
}

func (r *EntradaEstoqueRepository) ListarEntradasEstoque(ctx context.Context, tx *sql.Tx, filtro model.FiltroEntradaEstoque) ([]model.EntradaEstoque, error) {
	query := `
		SELECT ee.id, e.nome, f.razao_social, ee.valor_total, ee.status, ee.criado_em, u.nome
		FROM tb_entradas_estoque ee
		INNER JOIN tb_estoques e ON e.id = ee.id_estoque
		INNER JOIN tb_fornecedores f ON f.id = ee.id_fornecedor
		INNER JOIN public.tb_usuarios_gestao u ON u.id = ee.id_usuario
		WHERE 1 = 1`
	args := make([]any, 0, 4)

	if filtro.ID != 0 {
		args = append(args, filtro.ID)
		query += " AND ee.id = $" + strconv.Itoa(len(args))
	}
	if filtro.Fornecedor != "" {
		args = append(args, "%"+filtro.Fornecedor+"%")
		query += " AND f.razao_social ILIKE $" + strconv.Itoa(len(args))
	}
	if filtro.Data != "" {
		args = append(args, filtro.Data)
		query += " AND ee.criado_em::date = $" + strconv.Itoa(len(args)) + "::date"
	}

	if filtro.Status != "" {
		args = append(args, filtro.Status)
		query += " AND ee.status = $" + strconv.Itoa(len(args))
	}

	query += " ORDER BY ee.criado_em DESC, ee.id DESC"

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entradas []model.EntradaEstoque

	for rows.Next() {
		var entrada model.EntradaEstoque
		if err := rows.Scan(
			&entrada.ID,
			&entrada.Estoque,
			&entrada.Fornecedor,
			&entrada.ValorTotal,
			&entrada.Status,
			&entrada.CriadoEm,
			&entrada.Usuario,
		); err != nil {
			return nil, err
		}
		entradas = append(entradas, entrada)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entradas, nil

}
