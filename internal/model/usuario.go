package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Usuario representa a entidade principal que reflete a tabela de usuários.
type Usuario struct {
	ID           int64     `json:"id" db:"id"`
	Nome         string    `json:"nome" db:"nome"`
	Username     string    `json:"username" db:"username"`
	CPF          *string   `json:"cpf,omitempty" db:"cpf"`
	Telefone     *string   `json:"telefone,omitempty" db:"telefone"`
	Email        string    `json:"email" db:"email"`
	Senha        string    `json:"senha,omitempty" db:"senha"`
	Ativo        bool      `json:"ativo" db:"ativo"`
	CriadoEm     time.Time `json:"criado_em" db:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em" db:"atualizado_em"`
}

// DTOs (Data Transfer Objects) - Usados no tráfego da API

type UsuarioCriar struct {
	Nome     string `json:"nome"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Senha    string `json:"senha"`
}

type UsuarioLogin struct {
	Username string `json:"username"`
	Senha    string `json:"senha"`
}

type UsuarioBasico struct {
	ID       int64  `json:"id"`
	Nome     string `json:"nome"`
	Username string `json:"username"`
	Telefone string `json:"telefone"`
	Email    string `json:"email"`
}

// Validar verifica se os dados de criação são válidos
func (u *UsuarioCriar) Validar() error {
	var erros []error
	if u.Nome == "" {
		erros = append(erros, errors.New("o nome não foi informado"))
	}
	if u.Username == "" {
		erros = append(erros, errors.New("o username não foi informado"))
	}
	if u.Senha == "" {
		erros = append(erros, errors.New("a senha não foi informada"))
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}
	return nil
}

// HashSenha criptografa a senha do usuário
func (u *UsuarioCriar) HashSenha() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Senha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Senha = string(hash)
	return nil
}

// Validar verifica se os dados de login são válidos
func (u *UsuarioLogin) Validar() error {
	var erros []error
	if u.Username == "" {
		erros = append(erros, errors.New("o username não foi informado"))
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
func (u *UsuarioLogin) ValidarSenha(senhaDB string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(senhaDB), []byte(u.Senha)); err != nil {
		return errors.New("senha inválida")
	}
	return nil
}
