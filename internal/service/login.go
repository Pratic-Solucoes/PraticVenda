package service

import (
	"context"
	"database/sql"
	"errors"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type LoginService struct {
	db         *sql.DB
	repository *repository.Repository
}

func (s *LoginService) Login(ctx context.Context, usuario *model.UsuarioLogin) (uint64, string, string, error) {

	id, nome, senhaDB, schema, err := s.repository.Login.Login(ctx, usuario.Email)

	if err != nil {
		return 0, "", "", err
	}

	if err := usuario.ValidarSenha(senhaDB); err != nil {
		return 0, "", "", err
	}

	return id, nome, schema, nil
}

func (s *LoginService) LoginAdministrativo(ctx context.Context, usuario *model.UsuarioLogin) (uint64, string, error) {
	id, nome, senhaDB, ativo, err := s.repository.Login.LoginAdministrativo(ctx, usuario.Email)
	if err != nil {
		return 0, "", err
	}

	if err := usuario.ValidarSenha(senhaDB); err != nil {
		return 0, "", err
	}

	if !ativo {
		return 0, "", errors.New("usuario inativo")
	}

	return id, nome, nil

}
