package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type ProdutoService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *ProdutoService) CriarProduto(ctx context.Context, input *model.ProdutoInput) (*model.ProdutoCompleto, error) {
	if err := input.Validar(); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	p := &model.Produto{
		Composto:          input.Composto,
		MateriaPrima:      input.MateriaPrima,
		IDFornecedor:      fornecedorPrincipal(input),
		CodigoBarras:      input.CodigoBarras,
		CodigoInternoLoja: input.CodigoInternoLoja,
		Nome:              input.Nome,
		Descricao:         input.Descricao,
		PrecoCusto:        input.PrecoCusto,
		PrecoVenda:        input.PrecoVenda,
		UnidadeEstoque:    input.UnidadeEstoque,
		UnidadeVenda:      input.UnidadeVenda,
		PesoBruto:         input.PesoBruto,
		PesoLiquido:       input.PesoLiquido,
		Ativo:             true,
	}

	f := &model.ProdutoFiscal{
		Ncm:               input.Ncm,
		Cest:              input.Cest,
		IDGrupoTributario: input.IDGrupoTributario,
	}

	pCriado, err := s.repository.Produtos.CriarProduto(ctx, tx, p, f)
	if err != nil {
		return nil, err
	}
	if err = s.repository.Produtos.SincronizarFornecedores(ctx, tx, pCriado.ID, fornecedoresInput(input)); err != nil {
		return nil, err
	}
	if input.Composto {
		if err = s.repository.Produtos.SalvarComposicao(ctx, tx, pCriado.ID, input.Composicao); err != nil {
			return nil, err
		}
	}

	// Cria vínculos iniciais de estoque
	if !input.Composto {
		for _, est := range input.Estoques {
			err = s.repository.Produtos.VincularAoEstoque(ctx, tx, pCriado.ID, est.IDEstoque, est.EstoqueMinimo, est.Quantidade)
			if err != nil {
				return nil, fmt.Errorf("erro ao vincular estoque %d: %w", est.IDEstoque, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.ObterProdutoPorID(ctx, pCriado.ID)
}

func (s *ProdutoService) ListarProdutos(ctx context.Context, busca string, idFornecedor int64) ([]*model.ProdutoCompleto, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	produtos, err := s.repository.Produtos.ListarProdutos(ctx, tx, busca, idFornecedor)
	if err != nil {
		return nil, err
	}

	// Adiciona informações de estoques vinculados a cada produto
	for _, p := range produtos {
		vinculos, err := s.repository.Produtos.BuscarEstoqueVinculos(ctx, tx, p.Produto.ID)
		if err != nil {
			return nil, err
		}
		p.Estoques = vinculos
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return produtos, nil
}

func (s *ProdutoService) ObterProdutoPorID(ctx context.Context, id int64) (*model.ProdutoCompleto, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	p, err := s.repository.Produtos.ObterProdutoPorID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	vinculos, err := s.repository.Produtos.BuscarEstoqueVinculos(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	p.Estoques = vinculos
	p.IDsFornecedores, err = s.repository.Produtos.ListarFornecedoresProduto(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProdutoService) AtualizarProduto(ctx context.Context, id int64, input *model.ProdutoInput) error {
	if err := input.Validar(); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	p, err := s.repository.Produtos.ObterProdutoPorID(ctx, tx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("produto não encontrado")
	}

	produtoUpdate := &model.Produto{
		Composto:          input.Composto,
		MateriaPrima:      input.MateriaPrima,
		IDFornecedor:      fornecedorPrincipal(input),
		CodigoBarras:      input.CodigoBarras,
		CodigoInternoLoja: input.CodigoInternoLoja,
		Nome:              input.Nome,
		Descricao:         input.Descricao,
		PrecoCusto:        input.PrecoCusto,
		PrecoVenda:        input.PrecoVenda,
		UnidadeEstoque:    input.UnidadeEstoque,
		UnidadeVenda:      input.UnidadeVenda,
		PesoBruto:         input.PesoBruto,
		PesoLiquido:       input.PesoLiquido,
		Ativo:             p.Produto.Ativo, // Mantém o status ativo original do produto
	}

	fiscalUpdate := &model.ProdutoFiscal{
		IDProduto:         id,
		Ncm:               input.Ncm,
		Cest:              input.Cest,
		IDGrupoTributario: input.IDGrupoTributario,
	}

	err = s.repository.Produtos.AtualizarProduto(ctx, tx, id, produtoUpdate, fiscalUpdate)
	if err != nil {
		return err
	}
	if err = s.repository.Produtos.SincronizarFornecedores(ctx, tx, id, fornecedoresInput(input)); err != nil {
		return err
	}
	if input.Composto {
		if _, err = tx.ExecContext(ctx, `DELETE FROM tb_produtos_estoque WHERE id_produto=$1`, id); err != nil {
			return err
		}
		if err = s.repository.Produtos.SalvarComposicao(ctx, tx, id, input.Composicao); err != nil {
			return err
		}
		return tx.Commit()
	}

	// Sincronizar estoques
	vinculosAtuais, err := s.repository.Produtos.BuscarEstoqueVinculos(ctx, tx, id)
	if err != nil {
		return err
	}

	// Mapear vínculos enviados no input
	inputEstoquesMap := make(map[int64]model.ProdutoEstoqueInput)
	for _, est := range input.Estoques {
		inputEstoquesMap[est.IDEstoque] = est
	}

	// Desvincular estoques que não estão no input
	for _, atual := range vinculosAtuais {
		if _, existe := inputEstoquesMap[atual.IDEstoque]; !existe {
			// Validar se tem movimentação antes de desvincular
			temMov, err := s.repository.Produtos.TemMovimentacaoNoEstoqueEspecifico(ctx, tx, id, atual.IDEstoque)
			if err != nil {
				return err
			}
			if temMov {
				return fmt.Errorf("não é possível desvincular o produto do estoque '%s' porque ele possui movimentações", atual.NomeEstoque)
			}
			err = s.repository.Produtos.DesvincularDoEstoque(ctx, tx, id, atual.IDEstoque)
			if err != nil {
				return err
			}
		}
	}

	// Vincular/Atualizar os estoques enviados no input
	for _, est := range input.Estoques {
		err = s.repository.Produtos.VincularAoEstoque(ctx, tx, id, est.IDEstoque, est.EstoqueMinimo, est.Quantidade)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func fornecedoresInput(input *model.ProdutoInput) []int64 {
	if input.Composto {
		return []int64{}
	}
	if len(input.IDsFornecedores) > 0 {
		return input.IDsFornecedores
	}
	return []int64{input.IDFornecedor}
}

func fornecedorPrincipal(input *model.ProdutoInput) int64 {
	fornecedores := fornecedoresInput(input)
	if len(fornecedores) == 0 {
		return 0
	}
	return fornecedores[0]
}

func (s *ProdutoService) ExcluirOuInativarProduto(ctx context.Context, id int64) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	p, err := s.repository.Produtos.ObterProdutoPorID(ctx, tx, id)
	if err != nil {
		return "", err
	}
	if p == nil {
		return "", errors.New("produto não encontrado")
	}

	temMov, err := s.repository.Produtos.TemMovimentacaoEstoque(ctx, tx, id)
	if err != nil {
		return "", err
	}

	if temMov {
		err = s.repository.Produtos.InativarProduto(ctx, tx, id)
		if err != nil {
			return "", err
		}
		if err := tx.Commit(); err != nil {
			return "", err
		}
		return "inativado", nil
	}

	err = s.repository.Produtos.ExcluirProduto(ctx, tx, id)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return "excluido", nil
}

func (s *ProdutoService) VincularProdutoEstoque(ctx context.Context, idProduto, idEstoque int64, qtdMinima float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.repository.Produtos.VincularAoEstoque(ctx, tx, idProduto, idEstoque, qtdMinima, 0.0)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ProdutoService) DesvincularProdutoEstoque(ctx context.Context, idProduto, idEstoque int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	temMov, err := s.repository.Produtos.TemMovimentacaoNoEstoqueEspecifico(ctx, tx, idProduto, idEstoque)
	if err != nil {
		return err
	}
	if temMov {
		return errors.New("o produto possui movimentações neste estoque e não pode ser desvinculado")
	}

	err = s.repository.Produtos.DesvincularDoEstoque(ctx, tx, idProduto, idEstoque)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ProdutoService) ListarGruposTributarios(ctx context.Context) ([]*model.GrupoTributario, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	grupos, err := s.repository.Produtos.ListarGruposTributarios(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return grupos, nil
}

func (s *ProdutoService) ListarComposicao(ctx context.Context, idProduto int64) ([]model.ItemComposicaoProduto, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	itens, err := s.repository.Produtos.ListarComposicao(ctx, tx, idProduto)
	if err != nil {
		return nil, err
	}
	return itens, tx.Commit()
}

func (s *ProdutoService) SalvarComposicao(ctx context.Context, idProduto int64, itens []model.ItemComposicaoProduto) error {
	if idProduto <= 0 || len(itens) == 0 {
		return errors.New("informe ao menos uma matéria-prima")
	}
	vistos := make(map[int64]struct{}, len(itens))
	for _, item := range itens {
		if err := item.Validar(idProduto); err != nil {
			return err
		}
		if _, ok := vistos[item.IDProdutoComponente]; ok {
			return errors.New("uma matéria-prima só pode constar uma vez")
		}
		vistos[item.IDProdutoComponente] = struct{}{}
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if produto, err := s.repository.Produtos.ObterProdutoPorID(ctx, tx, idProduto); err != nil {
		return err
	} else if produto == nil {
		return errors.New("produto composto não encontrado")
	}
	if err := s.repository.Produtos.SalvarComposicao(ctx, tx, idProduto, itens); err != nil {
		return err
	}
	return tx.Commit()
}
