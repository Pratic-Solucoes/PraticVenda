package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type ProdutoRepository struct {
	db *sql.DB
}

func (r *ProdutoRepository) CriarProduto(ctx context.Context, tx *sql.Tx, p *model.Produto, f *model.ProdutoFiscal) (*model.Produto, error) {
	queryProduct := `
		INSERT INTO tb_produtos (
			codigo_barras, codigo_interno_loja, nome, descricao, 
			preco_custo, preco_venda, unidade_estoque, unidade_venda, 
			peso_bruto, peso_liquido, ativo, criado_em, atualizado_em
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, criado_em, atualizado_em
	`
	err := tx.QueryRowContext(ctx, queryProduct,
		p.CodigoBarras, p.CodigoInternoLoja, p.Nome, p.Descricao,
		p.PrecoCusto, p.PrecoVenda, p.UnidadeEstoque, p.UnidadeVenda,
		p.PesoBruto, p.PesoLiquido, p.Ativo,
	).Scan(&p.ID, &p.CriadoEm, &p.AtualizadoEm)

	if err != nil {
		return nil, err
	}

	queryFiscal := `
		INSERT INTO tb_produtos_fiscal (id_produto, ncm, cest, id_grupo_tributario, atualizado_em)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	`
	_, err = tx.ExecContext(ctx, queryFiscal, p.ID, f.Ncm, f.Cest, f.IDGrupoTributario)
	if err != nil {
		return nil, err
	}

	f.IDProduto = p.ID
	return p, nil
}

func (r *ProdutoRepository) AtualizarProduto(ctx context.Context, tx *sql.Tx, id int64, p *model.Produto, f *model.ProdutoFiscal) error {
	queryProduct := `
		UPDATE tb_produtos 
		SET codigo_barras = $1, codigo_interno_loja = $2, nome = $3, descricao = $4, 
			preco_custo = $5, preco_venda = $6, unidade_estoque = $7, unidade_venda = $8, 
			peso_bruto = $9, peso_liquido = $10, ativo = $11, atualizado_em = CURRENT_TIMESTAMP
		WHERE id = $12
	`
	_, err := tx.ExecContext(ctx, queryProduct,
		p.CodigoBarras, p.CodigoInternoLoja, p.Nome, p.Descricao,
		p.PrecoCusto, p.PrecoVenda, p.UnidadeEstoque, p.UnidadeVenda,
		p.PesoBruto, p.PesoLiquido, p.Ativo, id,
	)
	if err != nil {
		return err
	}

	queryFiscal := `
		INSERT INTO tb_produtos_fiscal (id_produto, ncm, cest, id_grupo_tributario, atualizado_em)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (id_produto) DO UPDATE 
		SET ncm = EXCLUDED.ncm, cest = EXCLUDED.cest, id_grupo_tributario = EXCLUDED.id_grupo_tributario, atualizado_em = CURRENT_TIMESTAMP
	`
	_, err = tx.ExecContext(ctx, queryFiscal, id, f.Ncm, f.Cest, f.IDGrupoTributario)
	return err
}

func (r *ProdutoRepository) ListarProdutos(ctx context.Context, tx *sql.Tx, busca string) ([]*model.ProdutoCompleto, error) {
	query := `
		SELECT 
			p.id, p.codigo_barras, p.codigo_interno_loja, p.nome, p.descricao, p.preco_custo, p.preco_venda, 
			p.unidade_estoque, p.unidade_venda, p.peso_bruto, p.peso_liquido, p.ativo, p.criado_em, p.atualizado_em,
			pf.ncm, pf.cest, pf.id_grupo_tributario, pf.atualizado_em
		FROM tb_produtos p
		LEFT JOIN tb_produtos_fiscal pf ON p.id = pf.id_produto
		WHERE ($1 = '' OR p.nome ILIKE $2 OR p.codigo_barras = $1 OR p.codigo_interno_loja = $1)
		ORDER BY p.nome
	`
	rows, err := tx.QueryContext(ctx, query, busca, "%"+busca+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []*model.ProdutoCompleto
	for rows.Next() {
		var pc model.ProdutoCompleto
		var pf model.ProdutoFiscal
		var hasFiscal sql.NullInt64

		err := rows.Scan(
			&pc.Produto.ID, &pc.Produto.CodigoBarras, &pc.Produto.CodigoInternoLoja, &pc.Produto.Nome, &pc.Produto.Descricao,
			&pc.Produto.PrecoCusto, &pc.Produto.PrecoVenda, &pc.Produto.UnidadeEstoque, &pc.Produto.UnidadeVenda,
			&pc.Produto.PesoBruto, &pc.Produto.PesoLiquido, &pc.Produto.Ativo, &pc.Produto.CriadoEm, &pc.Produto.AtualizadoEm,
			&pf.Ncm, &pf.Cest, &hasFiscal, &pf.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}

		if hasFiscal.Valid {
			pf.IDGrupoTributario = hasFiscal.Int64
			pf.IDProduto = pc.Produto.ID
			pc.Fiscal = &pf
		}

		produtos = append(produtos, &pc)
	}

	return produtos, nil
}

func (r *ProdutoRepository) ObterProdutoPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.ProdutoCompleto, error) {
	query := `
		SELECT 
			p.id, p.codigo_barras, p.codigo_interno_loja, p.nome, p.descricao, p.preco_custo, p.preco_venda, 
			p.unidade_estoque, p.unidade_venda, p.peso_bruto, p.peso_liquido, p.ativo, p.criado_em, p.atualizado_em,
			pf.ncm, pf.cest, pf.id_grupo_tributario, pf.atualizado_em
		FROM tb_produtos p
		LEFT JOIN tb_produtos_fiscal pf ON p.id = pf.id_produto
		WHERE p.id = $1
	`
	var pc model.ProdutoCompleto
	var pf model.ProdutoFiscal
	var hasFiscal sql.NullInt64

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&pc.Produto.ID, &pc.Produto.CodigoBarras, &pc.Produto.CodigoInternoLoja, &pc.Produto.Nome, &pc.Produto.Descricao,
		&pc.Produto.PrecoCusto, &pc.Produto.PrecoVenda, &pc.Produto.UnidadeEstoque, &pc.Produto.UnidadeVenda,
		&pc.Produto.PesoBruto, &pc.Produto.PesoLiquido, &pc.Produto.Ativo, &pc.Produto.CriadoEm, &pc.Produto.AtualizadoEm,
		&pf.Ncm, &pf.Cest, &hasFiscal, &pf.AtualizadoEm,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if hasFiscal.Valid {
		pf.IDGrupoTributario = hasFiscal.Int64
		pf.IDProduto = pc.Produto.ID
		pc.Fiscal = &pf
	}

	return &pc, nil
}

func (r *ProdutoRepository) ExcluirProduto(ctx context.Context, tx *sql.Tx, id int64) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM tb_produtos WHERE id = $1", id)
	return err
}

func (r *ProdutoRepository) InativarProduto(ctx context.Context, tx *sql.Tx, id int64) error {
	_, err := tx.ExecContext(ctx, "UPDATE tb_produtos SET ativo = FALSE, atualizado_em = CURRENT_TIMESTAMP WHERE id = $1", id)
	return err
}

func (r *ProdutoRepository) TemMovimentacaoEstoque(ctx context.Context, tx *sql.Tx, idProduto int64) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM tb_movimento_estoque WHERE id_produto = $1"
	err := tx.QueryRowContext(ctx, query, idProduto).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProdutoRepository) TemMovimentacaoNoEstoqueEspecifico(ctx context.Context, tx *sql.Tx, idProduto, idEstoque int64) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM tb_movimento_estoque WHERE id_produto = $1 AND id_estoque = $2"
	err := tx.QueryRowContext(ctx, query, idProduto, idEstoque).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProdutoRepository) VincularAoEstoque(ctx context.Context, tx *sql.Tx, idProduto, idEstoque int64, quantidadeMinima, quantidade float64) error {
	query := `
		INSERT INTO tb_produtos_estoque (id_produto, id_estoque, quantidade, estoque_minimo, atualizado_em)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (id_produto, id_estoque) DO UPDATE 
		SET estoque_minimo = EXCLUDED.estoque_minimo, atualizado_em = CURRENT_TIMESTAMP
	`
	_, err := tx.ExecContext(ctx, query, idProduto, idEstoque, quantidade, quantidadeMinima)
	return err
}

func (r *ProdutoRepository) DesvincularDoEstoque(ctx context.Context, tx *sql.Tx, idProduto, idEstoque int64) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM tb_produtos_estoque WHERE id_produto = $1 AND id_estoque = $2", idProduto, idEstoque)
	return err
}

func (r *ProdutoRepository) BuscarEstoqueVinculos(ctx context.Context, tx *sql.Tx, idProduto int64) ([]model.ProdutoEstoqueInfo, error) {
	query := `
		SELECT pe.id_estoque, e.nome, pe.quantidade, pe.estoque_minimo
		FROM tb_produtos_estoque pe
		INNER JOIN tb_estoques e ON pe.id_estoque = e.id
		WHERE pe.id_produto = $1
		ORDER BY e.nome
	`
	rows, err := tx.QueryContext(ctx, query, idProduto)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vinculos []model.ProdutoEstoqueInfo
	for rows.Next() {
		var info model.ProdutoEstoqueInfo
		err := rows.Scan(&info.IDEstoque, &info.NomeEstoque, &info.Quantidade, &info.EstoqueMinimo)
		if err != nil {
			return nil, err
		}
		vinculos = append(vinculos, info)
	}
	return vinculos, nil
}

func (r *ProdutoRepository) ListarGruposTributarios(ctx context.Context, tx *sql.Tx) ([]*model.GrupoTributario, error) {
	query := `
		SELECT 
			id, nome, cfop_padrao, origem_mercadoria, csosn, icms_cst, icms_aliquota, icms_mva_st, 
			icms_aliquota_st, ipi_cst, ipi_aliquota, pis_cst, pis_aliquota, cofins_cst, cofins_aliquota, 
			criado_em, atualizado_em
		FROM tb_grupos_tributarios
		ORDER BY nome
	`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grupos []*model.GrupoTributario
	for rows.Next() {
		var gt model.GrupoTributario
		err := rows.Scan(
			&gt.ID, &gt.Nome, &gt.CfopPadrao, &gt.OrigemMercadoria, &gt.Csosn, &gt.IcmsCst, &gt.IcmsAliquota, &gt.IcmsMvaSt,
			&gt.IcmsAliquotaSt, &gt.IpiCst, &gt.IpiAliquota, &gt.PisCst, &gt.PisAliquota, &gt.CofinsCst, &gt.CofinsAliquota,
			&gt.CriadoEm, &gt.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}
		grupos = append(grupos, &gt)
	}

	return grupos, nil
}
