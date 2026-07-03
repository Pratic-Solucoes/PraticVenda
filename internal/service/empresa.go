package service

import (
	"database/sql"
	"gestao/internal/repository"
)

type EmpresaService struct {
	db         *sql.DB
	repository *repository.Repository
}

/*
func (s *EmpresaService) CriarEmpresa(ctx context.Context, empresa *model.Empresa) error {

	if err := empresa.Validar(); err != nil {
		return err
	}
}
*/
