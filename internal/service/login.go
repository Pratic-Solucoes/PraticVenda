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
	if err := usuario.Validar(); err != nil {
		return 0, "", "", err
	}

	id, nome, senhaDB, schema, err := s.repository.Login.Login(ctx, usuario.Email)

	if err != nil {
		return 0, "", "", err
	}

	if err := usuario.ValidarSenha(senhaDB); err != nil {
		return 0, "", "", errors.New("dados login inválidos")
	}

	return id, nome, schema, nil
}
