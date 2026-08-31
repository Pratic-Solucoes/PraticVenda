package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type LoginService struct {
	db         *sql.DB
	repository *repository.Repository
}

func (s *LoginService) Login(ctx context.Context, usuario *model.UsuarioLogin) (uint64, string, error) {

	id, nome, senhaDB, err := s.repository.Login.Login(ctx, usuario.Username)

	if err != nil {
		return 0, "", err
	}

	if err := usuario.ValidarSenha(senhaDB); err != nil {
		return 0, "", err
	}

	return id, nome, nil
}
