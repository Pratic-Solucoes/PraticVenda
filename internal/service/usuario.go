package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"gestao/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UsuarioService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *UsuarioService) CriarUsuario(ctx context.Context, usuario *model.UsuarioCriar) (*model.UsuarioBasico, error) {

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

	usuarioCriado, err := s.repository.Usuarios.CriarUsuario(ctx, tx, usuario)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return usuarioCriado, nil
}

func (s *UsuarioService) BuscarUsuarioPorID(ctx context.Context, usuarioID int) (*model.Usuario, error) {

	return s.repository.Usuarios.BuscarUsuarioPorID(ctx, usuarioID)
}

func (s *UsuarioService) AlterarSenha(ctx context.Context, usuarioID int64, senhaAtual, senhaNova, senhaNovaConfirmacao string) error {

	senhaArmazenada, err := s.repository.Usuarios.BuscarSenhaUsuario(ctx, usuarioID)
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

	if err := s.repository.Usuarios.AtualizarSenhaUsuario(ctx, tx, usuarioID, senhaNova); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
