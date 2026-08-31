package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type CaixaService struct{ db *sql.DB }

func (s *CaixaService) Criar(ctx context.Context, usuario int64, nome string) (*model.Caixa, error) {
	if nome == "" {
		nome = "Caixa principal"
	}
	var c model.Caixa
	err := s.db.QueryRowContext(ctx, `INSERT INTO tb_caixas(id_usuario,nome) VALUES($1,$2) RETURNING id,nome`, usuario, nome).Scan(&c.ID, &c.Nome)
	return &c, err
}

func (s *CaixaService) ListarUsuarios(ctx context.Context) ([]model.UsuarioBasico, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,nome,username,COALESCE(telefone,''),email FROM tb_usuarios_gestao WHERE ativo=true ORDER BY nome`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var usuarios []model.UsuarioBasico
	for rows.Next() {
		var u model.UsuarioBasico
		if err := rows.Scan(&u.ID, &u.Nome, &u.Username, &u.Telefone, &u.Email); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, rows.Err()
}
func (s *CaixaService) Listar(ctx context.Context, usuario int64) ([]model.Caixa, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,nome FROM tb_caixas WHERE id_usuario=$1 AND ativo=true ORDER BY nome`, usuario)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var caixas []model.Caixa
	for rows.Next() {
		var c model.Caixa
		if err := rows.Scan(&c.ID, &c.Nome); err != nil {
			return nil, err
		}
		caixas = append(caixas, c)
	}
	return caixas, rows.Err()
}
func (s *CaixaService) Atual(ctx context.Context, usuario int64) (*model.ControleCaixa, error) {
	c := &model.ControleCaixa{}
	err := s.db.QueryRowContext(ctx, `SELECT id,id_caixa,status,valor_abertura FROM tb_controle_caixa WHERE id_usuario=$1 AND status='ABERTO' ORDER BY aberto_em DESC LIMIT 1`, usuario).Scan(&c.ID, &c.IDCaixa, &c.Status, &c.ValorAbertura)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}
func (s *CaixaService) Abrir(ctx context.Context, usuario, caixa int64, valor float64, senha string) (*model.ControleCaixa, error) {
	if valor < 0 {
		return nil, errors.New("o valor de abertura não pode ser negativo")
	}
	if err := s.validarSenha(ctx, usuario, senha); err != nil {
		return nil, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var pertence bool
	if err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tb_caixas WHERE id=$1 AND id_usuario=$2 AND ativo=true)`, caixa, usuario).Scan(&pertence); err != nil {
		return nil, err
	}
	if !pertence {
		return nil, errors.New("caixa inválido para o usuário")
	}
	c := &model.ControleCaixa{}
	err = tx.QueryRowContext(ctx, `INSERT INTO tb_controle_caixa(id_caixa,id_usuario,valor_abertura) VALUES($1,$2,$3) RETURNING id,id_caixa,status,valor_abertura`, caixa, usuario, valor).Scan(&c.ID, &c.IDCaixa, &c.Status, &c.ValorAbertura)
	if err != nil {
		return nil, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO tb_movimento_caixa(id_controle_caixa,tipo_movimento,tipo_credito,valor_credito) VALUES($1,'ABERTURA','DINHEIRO',$2)`, c.ID, valor); err != nil {
		return nil, err
	}
	return c, tx.Commit()
}
func (s *CaixaService) Fechar(ctx context.Context, usuario int64, f model.FechamentoCaixa) error {
	if err := f.Validar(); err != nil {
		return err
	}
	if err := s.validarSenha(ctx, usuario, f.SenhaAcesso); err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var id int64
	var abertura, dinheiro, cartaoPix float64
	err = tx.QueryRowContext(ctx, `SELECT c.id,c.valor_abertura,COALESCE(SUM(CASE WHEN m.tipo_credito='DINHEIRO' THEN CASE WHEN m.tipo_movimento='ESTORNO' THEN -m.valor_credito ELSE m.valor_credito END ELSE 0 END),0),COALESCE(SUM(CASE WHEN m.tipo_credito IN ('CARTAO','PIX') THEN CASE WHEN m.tipo_movimento='ESTORNO' THEN -m.valor_credito ELSE m.valor_credito END ELSE 0 END),0) FROM tb_controle_caixa c LEFT JOIN tb_movimento_caixa m ON m.id_controle_caixa=c.id AND m.tipo_movimento IN ('VENDA','ESTORNO') WHERE c.id_usuario=$1 AND c.status='ABERTO' GROUP BY c.id,c.valor_abertura`, usuario).Scan(&id, &abertura, &dinheiro, &cartaoPix)
	if err == sql.ErrNoRows {
		return errors.New("não há caixa aberto")
	}
	if err != nil {
		return err
	}
	dinheiro += abertura
	diferencaDinheiro, diferencaCartao := f.ValorDinheiro-dinheiro, f.ValorCartao-cartaoPix
	_, err = tx.ExecContext(ctx, `UPDATE tb_controle_caixa SET status='FECHADO',fechado_em=NOW(),valor_dinheiro_informado=$1,valor_cartao_informado=$2,valor_dinheiro_esperado=$3,valor_cartao_esperado=$4,diferenca_dinheiro=$5,diferenca_cartao=$6 WHERE id=$7`, f.ValorDinheiro, f.ValorCartao, dinheiro, cartaoPix, diferencaDinheiro, diferencaCartao, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *CaixaService) validarSenha(ctx context.Context, usuario int64, senha string) error {
	if senha == "" {
		return errors.New("a senha de acesso é obrigatória")
	}
	var hash string
	if err := s.db.QueryRowContext(ctx, `SELECT senha FROM tb_usuarios_gestao WHERE id=$1 AND ativo=true`, usuario).Scan(&hash); err != nil {
		return errors.New("usuário inválido")
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha)) != nil {
		return errors.New("senha de acesso inválida")
	}
	return nil
}
