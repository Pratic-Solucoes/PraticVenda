package usuario

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Usuario representa a entidade principal que reflete a tabela tb_usuarios_gestao
type Usuario struct {
	ID              int64      `json:"id,omitempty" db:"id"`
	Nome            string     `json:"nome,omitempty" db:"nome"`
	Username        string     `json:"username,omitempty" db:"username"`
	Email           string     `json:"email,omitempty" db:"email"`
	Celular         *string    `json:"celular,omitempty" db:"celular"`
	Senha           string     `json:"-" db:"senha"` // Oculto no JSON por segurança
	TermosAceitos   bool       `json:"termos_aceitos,omitempty"`
	TermosAceitosEm *time.Time `json:"termos_aceitos_em,omitempty" db:"termos_aceitos_em"`
	Ativo           bool       `json:"ativo,omitempty" db:"ativo"`
	CriadoEm        time.Time  `json:"criado_em,omitempty" db:"criado_em"`
	AtualizadoEm    time.Time  `json:"atualizado_em,omitempty" db:"atualizado_em"`
}

// Validar verifica se os dados de criação são válidos
func (u *Usuario) Validar() error {
	var erros []error
	if u.Nome == "" {
		erros = append(erros, errors.New("o nome não foi informado"))
	}
	if u.Username == "" {
		erros = append(erros, errors.New("o username não foi informado"))
	}
	if u.Email == "" {
		erros = append(erros, errors.New("o email não foi informado"))
	}
	if u.Senha == "" {
		erros = append(erros, errors.New("a senha não foi informada"))
	}
	if !u.TermosAceitos {
		erros = append(erros, errors.New("é obrigatório aceitar os termos de uso"))
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}
	return nil
}

// HashSenha criptografa a senha do usuário
func (u *Usuario) HashSenha() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Senha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Senha = string(hash)
	return nil
}

// ValidarLogin verifica se os dados de login são válidos (por email ou username)
func (u *Usuario) ValidarLogin() error {
	var erros []error
	if u.Email == "" && u.Username == "" {
		erros = append(erros, errors.New("o email ou username não foi informado"))
	}
	if u.Senha == "" {
		erros = append(erros, errors.New("a senha não foi informada"))
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}
	return nil
}

// ValidarSenha compara a senha enviada no login com a do banco
func (u *Usuario) ValidarSenha(senhaDB string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(senhaDB), []byte(u.Senha)); err != nil {
		return errors.New("senha inválida")
	}
	return nil
}
