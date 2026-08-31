package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type SaidaEstoqueService struct {
	repository *repository.Repository
	db         *sql.DB
}

// RegistrarSaida guarda uma saída aberta para revisão antes da baixa definitiva.
func (s *SaidaEstoqueService) RegistrarSaida(c context.Context, x *model.SaidaEstoque) error {
	if x.IDEstoque == 0 || len(x.Produtos) == 0 {
		return errors.New("estoque e produtos são obrigatórios")
	}
	if e := validarItensSaida(x.Produtos); e != nil {
		return e
	}
	var t float64
	for _, p := range x.Produtos {
		t += p.ValorTotal
	}
	x.ValorTotal = t
	statusSolicitado := x.Status
	if statusSolicitado == "" {
		statusSolicitado = "ABERTO"
	}
	if statusSolicitado != "ABERTO" && statusSolicitado != "CONCLUIDA" {
		return errors.New("uma nova saída deve ser salva em aberto ou concluída")
	}
	// A saída sempre nasce em aberto no banco. Quando solicitada como concluída,
	// a aprovação é realizada na mesma transação após gravar os itens.
	x.Status = "ABERTO"
	tx, e := s.db.BeginTx(c, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if e = s.repository.SaidaEstoque.RegistrarSaida(c, tx, x); e != nil {
		return e
	}
	if e = s.repository.SaidaEstoque.SalvarItens(c, tx, x); e != nil {
		return e
	}
	if statusSolicitado == "CONCLUIDA" {
		if e = s.repository.SaidaEstoque.Aprovar(c, tx, x.ID, x.IDUsuario); e != nil {
			return e
		}
		x.Status = "CONCLUIDA"
	}
	return tx.Commit()
}
func (s *SaidaEstoqueService) ListarSaidas(c context.Context, f model.FiltroSaidaEstoque) ([]model.SaidaEstoque, error) {
	tx, e := s.db.BeginTx(c, nil)
	if e != nil {
		return nil, e
	}
	defer tx.Rollback()
	out, e := s.repository.SaidaEstoque.Listar(c, tx, f)
	if e != nil {
		return nil, e
	}
	return out, tx.Commit()
}
func (s *SaidaEstoqueService) ObterSaida(c context.Context, id uint64) (*model.SaidaEstoque, error) {
	tx, e := s.db.BeginTx(c, nil)
	if e != nil {
		return nil, e
	}
	defer tx.Rollback()
	out, e := s.repository.SaidaEstoque.Obter(c, tx, id)
	if e != nil {
		return nil, e
	}
	return out, tx.Commit()
}

func (s *SaidaEstoqueService) EditarSaida(c context.Context, x *model.SaidaEstoque) error {
	if x.ID == 0 || x.IDEstoque == 0 || len(x.Produtos) == 0 {
		return errors.New("estoque e produtos são obrigatórios")
	}
	if e := validarItensSaida(x.Produtos); e != nil {
		return e
	}
	var total float64
	for _, p := range x.Produtos {
		total += p.ValorTotal
	}
	x.ValorTotal = total
	tx, e := s.db.BeginTx(c, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if e = s.repository.SaidaEstoque.Atualizar(c, tx, x); e != nil {
		return e
	}
	return tx.Commit()
}

// AprovarSaida efetiva a baixa de saldo após a validação transacional.
func (s *SaidaEstoqueService) AprovarSaida(c context.Context, id uint64, uid int64) error {
	tx, e := s.db.BeginTx(c, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if e = s.repository.SaidaEstoque.Aprovar(c, tx, id, uid); e != nil {
		return e
	}
	return tx.Commit()
}

func (s *SaidaEstoqueService) CancelarSaida(c context.Context, id uint64, usuario int64) error {
	tx, e := s.db.BeginTx(c, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if e = s.repository.SaidaEstoque.CancelarComUsuario(c, tx, id, usuario); e != nil {
		return e
	}
	return tx.Commit()
}

// validarItensSaida evita que o mesmo produto seja gravado em mais de uma
// linha. Além de tornar a saída clara para o usuário, impede validações de
// saldo parciais quando a operação for aprovada.
func validarItensSaida(itens []model.ProdutoSaidaEstoque) error {
	produtos := make(map[uint64]struct{}, len(itens))
	for _, p := range itens {
		if e := p.Validar(); e != nil {
			return e
		}
		if _, existe := produtos[p.IDProduto]; existe {
			return errors.New("um produto só pode constar uma vez na saída")
		}
		produtos[p.IDProduto] = struct{}{}
	}
	return nil
}
