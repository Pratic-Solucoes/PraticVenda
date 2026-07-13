package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type EstoqueRepository struct {
	db *sql.DB
}

func (r *EstoqueRepository) CriarEstoque(ctx context.Context, tx *sql.Tx, e *model.Estoque) (*model.Estoque, error) {
	query := `INSERT INTO tb_estoques (nome, descricao, ativo) VALUES ($1, $2, $3) RETURNING id, criado_em`

	err := tx.QueryRowContext(ctx, query, e.Nome, e.Descricao, e.Ativo).Scan(&e.ID, &e.CriadoEm)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *EstoqueRepository) ListarEstoques(ctx context.Context, tx *sql.Tx) ([]*model.Estoque, error) {
	query := `SELECT id, nome, descricao, ativo, criado_em FROM tb_estoques ORDER BY nome`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estoques []*model.Estoque
	for rows.Next() {
		var e model.Estoque
		err := rows.Scan(&e.ID, &e.Nome, &e.Descricao, &e.Ativo, &e.CriadoEm)
		if err != nil {
			return nil, err
		}
		estoques = append(estoques, &e)
	}

	return estoques, nil
}

func (r *EstoqueRepository) ListarProdutosDoEstoque(ctx context.Context, tx *sql.Tx, idEstoque int64) ([]*model.ProdutoEstoque, error) {
	query := `
		SELECT 
			pe.id, pe.id_produto, pe.id_estoque, pe.quantidade, pe.estoque_minimo, pe.atualizado_em,
			p.id, p.codigo_barras, p.codigo_interno_loja, p.nome, p.descricao, p.preco_custo, p.preco_venda, 
			p.unidade_estoque, p.unidade_venda, p.peso_bruto, p.peso_liquido, p.ativo, p.criado_em, p.atualizado_em
		FROM tb_produtos_estoque pe
		INNER JOIN tb_produtos p ON pe.id_produto = p.id
		WHERE pe.id_estoque = $1
		ORDER BY p.nome
	`

	rows, err := tx.QueryContext(ctx, query, idEstoque)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtosEstoque []*model.ProdutoEstoque
	for rows.Next() {
		var pe model.ProdutoEstoque
		var p model.Produto
		err := rows.Scan(
			&pe.ID, &pe.IDProduto, &pe.IDEstoque, &pe.Quantidade, &pe.EstoqueMinimo, &pe.AtualizadoEm,
			&p.ID, &p.CodigoBarras, &p.CodigoInternoLoja, &p.Nome, &p.Descricao, &p.PrecoCusto, &p.PrecoVenda,
			&p.UnidadeEstoque, &p.UnidadeVenda, &p.PesoBruto, &p.PesoLiquido, &p.Ativo, &p.CriadoEm, &p.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}
		pe.Produto = &p
		produtosEstoque = append(produtosEstoque, &pe)
	}

	return produtosEstoque, nil
}
