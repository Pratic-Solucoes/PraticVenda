package service

import (
	"context"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type AdminService struct {
	repository *repository.Repository
}

func (s *AdminService) CarregarOrganizacoes(ctx context.Context) ([]model.Organizacao, error) {
	return s.repository.Admin.CarregarOrganizacoes(ctx)
}

func (s *AdminService) CarregarUsuarios(ctx context.Context) ([]model.Usuario, error) {
	return s.repository.Admin.CarregarUsuarios(ctx)
}
