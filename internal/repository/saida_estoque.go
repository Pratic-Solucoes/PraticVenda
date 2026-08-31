package repository

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
)

type SaidaEstoqueRepository struct{ db *sql.DB }

// RegistrarSaida cria o cabeçalho aberto; o saldo só muda na aprovação.
func (r *SaidaEstoqueRepository) RegistrarSaida(c context.Context, tx *sql.Tx, s *model.SaidaEstoque) error {
	return tx.QueryRowContext(c, `INSERT INTO tb_saidas_estoque (id_estoque,id_usuario,valor_total,status) VALUES ($1,$2,$3,$4) RETURNING id`, s.IDEstoque, s.IDUsuario, s.ValorTotal, s.Status).Scan(&s.ID)
}
func (r *SaidaEstoqueRepository) SalvarItens(c context.Context, tx *sql.Tx, s *model.SaidaEstoque) error {
	for _, p := range s.Produtos {
		if _, e := tx.ExecContext(c, `INSERT INTO tb_produtos_saidas_estoque (id_saida_estoque,id_produto,valor_unitario,valor_custo,valor_total,quantidade) VALUES ($1,$2,$3,$4,$5,$6)`, s.ID, p.IDProduto, p.ValorUnitario, p.ValorCusto, p.ValorTotal, p.Quantidade); e != nil {
			return e
		}
	}
	return nil
}
func (r *SaidaEstoqueRepository) Listar(c context.Context, tx *sql.Tx, f model.FiltroSaidaEstoque) ([]model.SaidaEstoque, error) {
	rows, e := tx.QueryContext(c, `SELECT s.id,e.nome,s.valor_total,s.status,s.criado_em,u.nome FROM tb_saidas_estoque s JOIN tb_estoques e ON e.id=s.id_estoque JOIN public.tb_usuarios_gestao u ON u.id=s.id_usuario WHERE ($1=0 OR s.id=$1) AND ($2='' OR s.status=$2) AND ($3='' OR s.criado_em::date=$3::date) ORDER BY s.id DESC`, f.ID, f.Status, f.Data)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []model.SaidaEstoque
	for rows.Next() {
		var s model.SaidaEstoque
		if e = rows.Scan(&s.ID, &s.Estoque, &s.ValorTotal, &s.Status, &s.CriadoEm, &s.Usuario); e != nil {
			return nil, e
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
func (r *SaidaEstoqueRepository) Obter(c context.Context, tx *sql.Tx, id uint64) (*model.SaidaEstoque, error) {
	s := &model.SaidaEstoque{}
	e := tx.QueryRowContext(c, `SELECT s.id,s.id_estoque,e.nome,s.id_usuario,s.valor_total,s.status,s.criado_em,u.nome FROM tb_saidas_estoque s JOIN tb_estoques e ON e.id=s.id_estoque JOIN public.tb_usuarios_gestao u ON u.id=s.id_usuario WHERE s.id=$1`, id).Scan(&s.ID, &s.IDEstoque, &s.Estoque, &s.IDUsuario, &s.ValorTotal, &s.Status, &s.CriadoEm, &s.Usuario)
	if e != nil {
		return nil, e
	}
	rows, e := tx.QueryContext(c, `SELECT i.id,i.id_produto,p.nome,i.valor_unitario,i.valor_custo,i.valor_total,i.quantidade FROM tb_produtos_saidas_estoque i JOIN tb_produtos p ON p.id=i.id_produto WHERE i.id_saida_estoque=$1`, id)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	for rows.Next() {
		var p model.ProdutoSaidaEstoque
		if e = rows.Scan(&p.ID, &p.IDProduto, &p.NomeProduto, &p.ValorUnitario, &p.ValorCusto, &p.ValorTotal, &p.Quantidade); e != nil {
			return nil, e
		}
		s.Produtos = append(s.Produtos, p)
	}
	return s, rows.Err()
}

// Atualizar substitui os itens apenas enquanto a saída ainda estiver aberta.
func (r *SaidaEstoqueRepository) Atualizar(c context.Context, tx *sql.Tx, s *model.SaidaEstoque) error {
	result, e := tx.ExecContext(c, `UPDATE tb_saidas_estoque SET id_estoque=$1,valor_total=$2 WHERE id=$3 AND status='ABERTO'`, s.IDEstoque, s.ValorTotal, s.ID)
	if e != nil {
		return e
	}
	n, e := result.RowsAffected()
	if e != nil {
		return e
	}
	if n == 0 {
		return errors.New("somente saídas em aberto podem ser editadas")
	}
	if _, e = tx.ExecContext(c, `DELETE FROM tb_produtos_saidas_estoque WHERE id_saida_estoque=$1`, s.ID); e != nil {
		return e
	}
	return r.SalvarItens(c, tx, s)
}

// Aprovar baixa apenas itens com saldo suficiente, dentro da transação atual.
func (r *SaidaEstoqueRepository) Aprovar(c context.Context, tx *sql.Tx, id uint64, uid int64) error {
	// O lock do cabeçalho serializa aprovações concorrentes da mesma saída.
	// Assim, uma segunda requisição só observa o status depois do commit da
	// primeira e não consegue baixar o saldo novamente.
	var idEstoque uint64
	var status string
	if e := tx.QueryRowContext(c, `SELECT id_estoque,status FROM tb_saidas_estoque WHERE id=$1 FOR UPDATE`, id).Scan(&idEstoque, &status); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return errors.New("saída não encontrada")
		}
		return e
	}
	if status != "ABERTO" {
		return errors.New("somente saídas em aberto podem ser aprovadas")
	}

	rows, e := tx.QueryContext(c, `SELECT id_produto,SUM(quantidade) FROM tb_produtos_saidas_estoque WHERE id_saida_estoque=$1 GROUP BY id_produto ORDER BY id_produto`, id)
	if e != nil {
		return e
	}
	defer rows.Close()
	var itens []model.ProdutoSaidaEstoque
	for rows.Next() {
		var p model.ProdutoSaidaEstoque
		if e = rows.Scan(&p.IDProduto, &p.Quantidade); e != nil {
			return e
		}
		itens = append(itens, p)
	}
	if e = rows.Err(); e != nil {
		return e
	}
	if len(itens) == 0 {
		return errors.New("a saída não possui itens")
	}

	for _, p := range itens {
		var quantidadeDisponivel float64
		if e = tx.QueryRowContext(c, `SELECT quantidade FROM tb_produtos_estoque WHERE id_produto=$1 AND id_estoque=$2 FOR UPDATE`, p.IDProduto, idEstoque).Scan(&quantidadeDisponivel); e != nil {
			return errors.New("produto não está vinculado ao estoque")
		}
		if quantidadeDisponivel < p.Quantidade {
			return errors.New("saldo insuficiente para aprovar a saída")
		}
	}
	if _, e = tx.ExecContext(c, `UPDATE tb_saidas_estoque SET status='CONCLUIDA' WHERE id=$1 AND status='ABERTO'`, id); e != nil {
		return e
	}
	for _, p := range itens {
		if _, e = tx.ExecContext(c, `UPDATE tb_produtos_estoque SET quantidade=quantidade-$1,atualizado_em=NOW() WHERE id_produto=$2 AND id_estoque=$3`, p.Quantidade, p.IDProduto, idEstoque); e != nil {
			return e
		}
		if _, e = tx.ExecContext(c, `INSERT INTO tb_movimento_estoque (id_produto,id_estoque,id_usuario,quantidade,tipo_movimento,id_categoria_movimento,id_origem) SELECT $1,$2,$3,$4,'SAIDA',id,$5 FROM tb_categoria_movimento_estoque WHERE nome='SAIDA DE ESTOQUE'`, p.IDProduto, idEstoque, uid, p.Quantidade, id); e != nil {
			return e
		}
	}
	return nil
}

func (r *SaidaEstoqueRepository) Cancelar(c context.Context, tx *sql.Tx, id uint64) error {
	return r.CancelarComUsuario(c, tx, id, 0)
}

func (r *SaidaEstoqueRepository) CancelarComUsuario(c context.Context, tx *sql.Tx, id uint64, usuario int64) error {
	var estoque uint64
	var status string
	if e := tx.QueryRowContext(c, `SELECT id_estoque,status FROM tb_saidas_estoque WHERE id=$1 FOR UPDATE`, id).Scan(&estoque, &status); e != nil {
		return e
	}
	if status == "CANCELADA" {
		return errors.New("saída já está cancelada")
	}
	if status == "ABERTO" {
		_, e := tx.ExecContext(c, `UPDATE tb_saidas_estoque SET status='CANCELADA' WHERE id=$1`, id)
		return e
	}
	rows, e := tx.QueryContext(c, `SELECT id_produto,SUM(quantidade) FROM tb_produtos_saidas_estoque WHERE id_saida_estoque=$1 GROUP BY id_produto`, id)
	if e != nil {
		return e
	}
	type itemCancelamento struct {
		produto uint64
		qtd     float64
	}
	itens := make([]itemCancelamento, 0)
	for rows.Next() {
		var produto uint64
		var qtd float64
		if e = rows.Scan(&produto, &qtd); e != nil {
			rows.Close()
			return e
		}
		itens = append(itens, itemCancelamento{produto: produto, qtd: qtd})
	}
	if e = rows.Err(); e != nil {
		rows.Close()
		return e
	}
	if e = rows.Close(); e != nil {
		return e
	}
	for _, item := range itens {
		produto, qtd := item.produto, item.qtd
		if _, e = tx.ExecContext(c, `UPDATE tb_produtos_estoque SET quantidade=quantidade+$1,atualizado_em=NOW() WHERE id_produto=$2 AND id_estoque=$3`, qtd, produto, estoque); e != nil {
			return e
		}
		if _, e = tx.ExecContext(c, `INSERT INTO tb_movimento_estoque(id_produto,id_estoque,id_usuario,quantidade,tipo_movimento,id_categoria_movimento,id_origem) SELECT $1,$2,$3,$4,'ENTRADA',id,$5 FROM tb_categoria_movimento_estoque WHERE nome='CANCELAMENTO DE SAIDA'`, produto, estoque, usuario, qtd, id); e != nil {
			return e
		}
	}
	_, e = tx.ExecContext(c, `UPDATE tb_saidas_estoque SET status='CANCELADA' WHERE id=$1`, id)
	return e
}
