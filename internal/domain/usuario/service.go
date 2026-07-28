package usuario

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type ServiceInterface interface {
	CriarUsuario(ctx context.Context, usuario *Usuario) (*Usuario, error)
	BuscarUsuarioPorID(ctx context.Context, usuarioID uint64) (*Usuario, error)
	AlterarSenha(ctx context.Context, usuarioID uint64, senhaAtual, senhaNova, senhaNovaConfirmacao string) error
}

type Service struct {
	repository RepositoryInterface
	db         *sql.DB
}

func NewService(r RepositoryInterface, db *sql.DB) ServiceInterface {
	return &Service{repository: r, db: db}
}

// Cria um usuário
func (s *Service) CriarUsuario(ctx context.Context, usuario *Usuario) (*Usuario, error) {

	if err := usuario.Validar(); err != nil {
		return nil, err
	}

	if err := usuario.HashSenha(); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	usuarioCriado, err := s.repository.CriarUsuario(ctx, tx, usuario)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return usuarioCriado, nil
}

// Buscar por ID
func (s *Service) BuscarUsuarioPorID(ctx context.Context, usuarioID uint64) (*Usuario, error) {

	return s.repository.BuscarUsuarioPorID(ctx, uint64(usuarioID))
}

// Altera a senha de um usuário
func (s *Service) AlterarSenha(ctx context.Context, usuarioID uint64, senhaAtual, senhaNova, senhaNovaConfirmacao string) error {

	senhaArmazenada, err := s.repository.BuscarSenhaUsuario(ctx, uint64(usuarioID))
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*senhaArmazenada), []byte(senhaAtual)); err != nil {
		return errors.New("Senha atual recebida: " + senhaAtual + " não corresponde à senha armazenada")
	}

	if senhaNova != senhaNovaConfirmacao {
		return errors.New("a nova senha e a confirmação não coincidem")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(senhaNova), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	senhaNova = string(hash)

	if err := s.repository.AtualizarSenhaUsuario(ctx, tx, usuarioID, senhaNova); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
