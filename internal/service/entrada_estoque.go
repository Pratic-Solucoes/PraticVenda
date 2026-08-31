package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"gestao/internal/repository"
	"strings"
	"time"
)

type EntradaEstoqueService struct {
	db          *sql.DB
	repositorio *repository.Repository
}

func (s *EntradaEstoqueService) RegistrarEntrada(ctx context.Context, entrada *model.EntradaEstoque) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if entrada.IDEstoque == 0 || entrada.IDFornecedor == 0 || len(entrada.Produtos) == 0 {
		return errors.New("dados da entrada são obrigatórios")
	}

	// Verifica se algum produto da entrada de estoque está com dado incorreto
	for _, produto := range entrada.Produtos {
		if err := produto.Validar(); err != nil {
			return err
		}
	}
	if err := s.repositorio.EntradaEstoque.ValidarProdutosFornecedor(ctx, tx, entrada.IDFornecedor, entrada.Produtos); err != nil {
		return err
	}
	// Calcula o valor total da entrada somando os produtos
	var total float64
	for _, p := range entrada.Produtos {
		total += p.ValorTotal
	}
	entrada.ValorTotal = total

	// Lógica para criar a entrada de estoque no BD
	if err := s.repositorio.EntradaEstoque.RegistrarEntrada(ctx, tx, entrada); err != nil {
		return err
	}

	if entrada.Status == "CONCLUIDA" {
		// Lógica para inserir os produtos, movimento de estoque e atualizar o saldo do produto no estoque
		// Só haverá movimentação de estoque caso a entrada seja concluída
		if err := s.repositorio.EntradaEstoque.RegistrarProdutosEntrada(ctx, tx, entrada); err != nil {
			return err
		}
	} else if err := s.repositorio.EntradaEstoque.SalvarItensEntrada(ctx, tx, entrada); err != nil {
		// Entradas abertas também guardam seus itens para visualização e edição posterior.
		return err
	}

	return tx.Commit()
}

// ObterEntrada carrega uma entrada com seus itens, sem alterar seu estado.
func (s *EntradaEstoqueService) ObterEntrada(ctx context.Context, id uint64) (*model.EntradaEstoque, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	entrada, err := s.repositorio.EntradaEstoque.ObterEntrada(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return entrada, nil
}

// EditarEntrada substitui os dados e itens de uma entrada aberta.
func (s *EntradaEstoqueService) EditarEntrada(ctx context.Context, entrada *model.EntradaEstoque) error {
	if entrada.ID == 0 || entrada.IDEstoque == 0 || entrada.IDFornecedor == 0 || len(entrada.Produtos) == 0 {
		return errors.New("dados da entrada são obrigatórios")
	}
	var total float64
	for _, produto := range entrada.Produtos {
		if err := produto.Validar(); err != nil {
			return err
		}
		total += produto.ValorTotal
	}
	entrada.ValorTotal = total
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.repositorio.EntradaEstoque.ValidarProdutosFornecedor(ctx, tx, entrada.IDFornecedor, entrada.Produtos); err != nil {
		return err
	}
	if err := s.repositorio.EntradaEstoque.AtualizarEntrada(ctx, tx, entrada); err != nil {
		return err
	}
	return tx.Commit()
}

// AprovarEntrada conclui uma entrada aberta e registra sua movimentação de estoque.
func (s *EntradaEstoqueService) AprovarEntrada(ctx context.Context, id uint64, usuarioID int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.repositorio.EntradaEstoque.AprovarEntrada(ctx, tx, id, usuarioID); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *EntradaEstoqueService) CancelarEntrada(ctx context.Context, id uint64, usuarioID int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = s.repositorio.EntradaEstoque.CancelarEntrada(ctx, tx, id, usuarioID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *EntradaEstoqueService) ListarEntradas(ctx context.Context, filtro model.FiltroEntradaEstoque) ([]model.EntradaEstoque, error) {
	filtro.Fornecedor = strings.TrimSpace(filtro.Fornecedor)
	filtro.Status = strings.ToUpper(strings.TrimSpace(filtro.Status))

	if filtro.Data != "" {
		if _, err := time.Parse("2006-01-02", filtro.Data); err != nil {
			return nil, errors.New("data do filtro inválida")
		}
	}
	if filtro.Status != "" && filtro.Status != "ABERTO" && filtro.Status != "CONCLUIDA" && filtro.Status != "CANCELADA" {
		return nil, errors.New("status do filtro inválido")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	entradas, err := s.repositorio.EntradaEstoque.ListarEntradasEstoque(ctx, tx, filtro)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return entradas, nil
}
