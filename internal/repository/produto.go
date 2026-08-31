package repository

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"strings"
)

type ProdutoRepository struct {
	db *sql.DB
}

func (r *ProdutoRepository) CriarProduto(ctx context.Context, tx *sql.Tx, p *model.Produto, f *model.ProdutoFiscal) (*model.Produto, error) {
	queryProduct := `
		INSERT INTO tb_produtos (
			id_fornecedor, composto, materia_prima, codigo_barras, codigo_interno_loja, nome, descricao,
			preco_custo, preco_venda, unidade_estoque, unidade_venda, 
			peso_bruto, peso_liquido, ativo, criado_em, atualizado_em
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, criado_em, atualizado_em
	`
	err := tx.QueryRowContext(ctx, queryProduct,
		nullableFornecedor(p.IDFornecedor), p.Composto, p.MateriaPrima, p.CodigoBarras, p.CodigoInternoLoja, p.Nome, p.Descricao,
		p.PrecoCusto, p.PrecoVenda, p.UnidadeEstoque, p.UnidadeVenda,
		p.PesoBruto, p.PesoLiquido, p.Ativo,
	).Scan(&p.ID, &p.CriadoEm, &p.AtualizadoEm)

	if err != nil {
		return nil, err
	}

	if possuiDadosFiscais(f) {
		queryFiscal := `
			INSERT INTO tb_produtos_fiscal (id_produto, ncm, cest, id_grupo_tributario, atualizado_em)
			VALUES ($1, NULLIF($2, ''), $3, $4, CURRENT_TIMESTAMP)
		`
		if _, err = tx.ExecContext(ctx, queryFiscal, p.ID, f.Ncm, f.Cest, idGrupoTributarioOuNulo(f)); err != nil {
			return nil, err
		}
	}

	f.IDProduto = p.ID
	return p, nil
}

func (r *ProdutoRepository) AtualizarProduto(ctx context.Context, tx *sql.Tx, id int64, p *model.Produto, f *model.ProdutoFiscal) error {
	queryProduct := `
		UPDATE tb_produtos 
		SET id_fornecedor=$1, composto=$2, materia_prima=$3, codigo_barras=$4, codigo_interno_loja=$5, nome=$6, descricao=$7,
			preco_custo=$8, preco_venda=$9, unidade_estoque=$10, unidade_venda=$11,
			peso_bruto=$12, peso_liquido=$13, ativo=$14, atualizado_em=CURRENT_TIMESTAMP WHERE id=$15
	`
	_, err := tx.ExecContext(ctx, queryProduct,
		nullableFornecedor(p.IDFornecedor), p.Composto, p.MateriaPrima, p.CodigoBarras, p.CodigoInternoLoja, p.Nome, p.Descricao,
		p.PrecoCusto, p.PrecoVenda, p.UnidadeEstoque, p.UnidadeVenda,
		p.PesoBruto, p.PesoLiquido, p.Ativo, id,
	)
	if err != nil {
		return err
	}

	if !possuiDadosFiscais(f) {
		_, err = tx.ExecContext(ctx, "DELETE FROM tb_produtos_fiscal WHERE id_produto = $1", id)
		return err
	}

	queryFiscal := `
		INSERT INTO tb_produtos_fiscal (id_produto, ncm, cest, id_grupo_tributario, atualizado_em)
		VALUES ($1, NULLIF($2, ''), $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (id_produto) DO UPDATE
		SET ncm = EXCLUDED.ncm, cest = EXCLUDED.cest, id_grupo_tributario = EXCLUDED.id_grupo_tributario, atualizado_em = CURRENT_TIMESTAMP
	`
	_, err = tx.ExecContext(ctx, queryFiscal, id, f.Ncm, f.Cest, idGrupoTributarioOuNulo(f))
	return err
}

func possuiDadosFiscais(f *model.ProdutoFiscal) bool {
	return strings.TrimSpace(f.Ncm) != "" || (f.Cest != nil && strings.TrimSpace(*f.Cest) != "") || f.IDGrupoTributario > 0
}

func idGrupoTributarioOuNulo(f *model.ProdutoFiscal) any {
	if f.IDGrupoTributario <= 0 {
		return nil
	}
	return f.IDGrupoTributario
}

func nullableFornecedor(id int64) any {
	if id <= 0 {
		return nil
	}
	return id
}

func (r *ProdutoRepository) ListarProdutos(ctx context.Context, tx *sql.Tx, busca string, idFornecedor int64) ([]*model.ProdutoCompleto, error) {
	query := `
		SELECT 
			p.id, p.composto, p.materia_prima, COALESCE(p.id_fornecedor,0), COALESCE(f.razao_social,''), p.codigo_barras, p.codigo_interno_loja, p.nome, p.descricao, p.preco_custo, p.preco_venda,
			p.unidade_estoque, p.unidade_venda, p.peso_bruto, p.peso_liquido, p.ativo, p.criado_em, p.atualizado_em,
			pf.ncm, pf.cest, pf.id_grupo_tributario, pf.atualizado_em
		FROM tb_produtos p
		LEFT JOIN tb_produtos_fiscal pf ON p.id = pf.id_produto LEFT JOIN tb_fornecedores f ON f.id = p.id_fornecedor
		WHERE ($1 = '' OR p.nome ILIKE $2 OR p.codigo_barras = $1 OR p.codigo_interno_loja = $1)
		AND ($3 = 0 OR EXISTS(SELECT 1 FROM tb_produtos_fornecedores pfv WHERE pfv.id_produto=p.id AND pfv.id_fornecedor=$3))
		ORDER BY p.nome
	`
	rows, err := tx.QueryContext(ctx, query, busca, "%"+busca+"%", idFornecedor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []*model.ProdutoCompleto
	for rows.Next() {
		var pc model.ProdutoCompleto
		var pf model.ProdutoFiscal
		var hasFiscal sql.NullInt64
		var ncmNull sql.NullString
		var fiscalAtualizadoEm sql.NullTime

		err := rows.Scan(
			&pc.Produto.ID, &pc.Produto.Composto, &pc.Produto.MateriaPrima, &pc.Produto.IDFornecedor, &pc.Produto.Fornecedor, &pc.Produto.CodigoBarras, &pc.Produto.CodigoInternoLoja, &pc.Produto.Nome, &pc.Produto.Descricao,
			&pc.Produto.PrecoCusto, &pc.Produto.PrecoVenda, &pc.Produto.UnidadeEstoque, &pc.Produto.UnidadeVenda,
			&pc.Produto.PesoBruto, &pc.Produto.PesoLiquido, &pc.Produto.Ativo, &pc.Produto.CriadoEm, &pc.Produto.AtualizadoEm,
			&ncmNull, &pf.Cest, &hasFiscal, &fiscalAtualizadoEm,
		)
		if err != nil {
			return nil, err
		}

		if hasFiscal.Valid {
			pf.IDGrupoTributario = hasFiscal.Int64
			pf.IDProduto = pc.Produto.ID
			pf.Ncm = ncmNull.String
			pf.AtualizadoEm = fiscalAtualizadoEm.Time
			pc.Fiscal = &pf
		}

		produtos = append(produtos, &pc)
	}

	return produtos, nil
}

func (r *ProdutoRepository) SincronizarFornecedores(ctx context.Context, tx *sql.Tx, idProduto int64, ids []int64) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM tb_produtos_fornecedores WHERE id_produto=$1`, idProduto); err != nil {
		return err
	}
	for _, id := range ids {
		if _, err := tx.ExecContext(ctx, `INSERT INTO tb_produtos_fornecedores(id_produto,id_fornecedor) VALUES($1,$2)`, idProduto, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProdutoRepository) ListarFornecedoresProduto(ctx context.Context, tx *sql.Tx, idProduto int64) ([]int64, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id_fornecedor FROM tb_produtos_fornecedores WHERE id_produto=$1 ORDER BY id_fornecedor`, idProduto)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *ProdutoRepository) ObterProdutoPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.ProdutoCompleto, error) {
	query := `
		SELECT 
			p.id, p.composto, p.materia_prima, COALESCE(p.id_fornecedor,0), COALESCE(f.razao_social,''), p.codigo_barras, p.codigo_interno_loja, p.nome, p.descricao, p.preco_custo, p.preco_venda,
			p.unidade_estoque, p.unidade_venda, p.peso_bruto, p.peso_liquido, p.ativo, p.criado_em, p.atualizado_em,
			pf.ncm, pf.cest, pf.id_grupo_tributario, pf.atualizado_em
		FROM tb_produtos p
		LEFT JOIN tb_produtos_fiscal pf ON p.id = pf.id_produto LEFT JOIN tb_fornecedores f ON f.id = p.id_fornecedor
		WHERE p.id = $1
	`
	var pc model.ProdutoCompleto
	var pf model.ProdutoFiscal
	var hasFiscal sql.NullInt64
	var ncmNull sql.NullString
	var fiscalAtualizadoEm sql.NullTime

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&pc.Produto.ID, &pc.Produto.Composto, &pc.Produto.MateriaPrima, &pc.Produto.IDFornecedor, &pc.Produto.Fornecedor, &pc.Produto.CodigoBarras, &pc.Produto.CodigoInternoLoja, &pc.Produto.Nome, &pc.Produto.Descricao,
		&pc.Produto.PrecoCusto, &pc.Produto.PrecoVenda, &pc.Produto.UnidadeEstoque, &pc.Produto.UnidadeVenda,
		&pc.Produto.PesoBruto, &pc.Produto.PesoLiquido, &pc.Produto.Ativo, &pc.Produto.CriadoEm, &pc.Produto.AtualizadoEm,
		&ncmNull, &pf.Cest, &hasFiscal, &fiscalAtualizadoEm,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if hasFiscal.Valid {
		pf.IDGrupoTributario = hasFiscal.Int64
		pf.IDProduto = pc.Produto.ID
		pf.Ncm = ncmNull.String
		pf.AtualizadoEm = fiscalAtualizadoEm.Time
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

func (r *ProdutoRepository) ListarComposicao(ctx context.Context, tx *sql.Tx, idProduto int64) ([]model.ItemComposicaoProduto, error) {
	rows, err := tx.QueryContext(ctx, `SELECT c.id_produto_componente,p.nome,p.unidade_estoque,c.quantidade,p.preco_custo FROM tb_composicoes_produtos c JOIN tb_produtos p ON p.id=c.id_produto_componente WHERE c.id_produto_composto=$1 ORDER BY p.nome`, idProduto)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var itens []model.ItemComposicaoProduto
	for rows.Next() {
		var item model.ItemComposicaoProduto
		if err = rows.Scan(&item.IDProdutoComponente, &item.NomeProduto, &item.UnidadeEstoque, &item.Quantidade, &item.PrecoCusto); err != nil {
			return nil, err
		}
		itens = append(itens, item)
	}
	return itens, rows.Err()
}

func (r *ProdutoRepository) SalvarComposicao(ctx context.Context, tx *sql.Tx, idProduto int64, itens []model.ItemComposicaoProduto) error {
	var produtoComposto bool
	if err := tx.QueryRowContext(ctx, `SELECT composto FROM tb_produtos WHERE id=$1`, idProduto).Scan(&produtoComposto); err != nil {
		return errors.New("produto composto não encontrado")
	}
	if !produtoComposto {
		return errors.New("apenas produtos compostos podem ter ficha técnica")
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM tb_composicoes_produtos WHERE id_produto_composto=$1`, idProduto); err != nil {
		return err
	}
	for _, item := range itens {
		var ativo, materiaPrima bool
		if err := tx.QueryRowContext(ctx, `SELECT ativo,materia_prima FROM tb_produtos WHERE id=$1`, item.IDProdutoComponente).Scan(&ativo, &materiaPrima); err != nil {
			return errors.New("produto componente não encontrado")
		}
		if !ativo {
			return errors.New("produto componente está inativo")
		}
		if !materiaPrima {
			return errors.New("os componentes da ficha técnica devem ser matérias-primas")
		}
		var composto bool
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tb_composicoes_produtos WHERE id_produto_composto=$1)`, item.IDProdutoComponente).Scan(&composto); err != nil {
			return err
		}
		if composto {
			return errors.New("um componente não pode ser outro produto composto")
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO tb_composicoes_produtos (id_produto_composto,id_produto_componente,quantidade) VALUES ($1,$2,$3)`, idProduto, item.IDProdutoComponente, item.Quantidade); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `UPDATE tb_produtos p SET preco_custo=COALESCE((SELECT SUM(c.quantidade * componente.preco_custo) FROM tb_composicoes_produtos c JOIN tb_produtos componente ON componente.id=c.id_produto_componente WHERE c.id_produto_composto=$1),0), atualizado_em=NOW() WHERE p.id=$1`, idProduto); err != nil {
		return err
	}
	return nil
}
