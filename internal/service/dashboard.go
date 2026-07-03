package service

import (
	"context"
	"database/sql"
	"gestao/internal/repository"
	"gestao/pkg/helpers"
	"time"
)

type DashboardService struct {
	Repo *repository.DashboardRepository
	db   *sql.DB
}

func NewDashboardService(repo *repository.DashboardRepository, db *sql.DB) *DashboardService {
	return &DashboardService{Repo: repo, db: db}
}

type ResumoDashboard struct {
	TotalVencido     float64                     `json:"total_vencido"`
	TotalSemana      float64                     `json:"total_semana"`
	VendaDia         float64                     `json:"venda_dia"`
	VendaMes         float64                     `json:"venda_mes"`
	DespesasCategoria []repository.CategoriaGasto `json:"despesas_categoria"`
}

// ObterResumo consulta e agrega todos os dados do dashboard em uma única chamada usando goroutines para performance
func (s *DashboardService) ObterResumo(ctx context.Context) (*ResumoDashboard, error) {
	hoje := time.Now()

	// Lógica para encontrar a Segunda-feira (início da semana) e Domingo (fim da semana)
	// Considerando que a semana começa na Segunda-feira
	offsetSegunda := int(time.Monday - hoje.Weekday())
	if offsetSegunda > 0 {
		offsetSegunda = -6
	}
	inicioSemana := hoje.AddDate(0, 0, offsetSegunda)
	fimSemana := inicioSemana.AddDate(0, 0, 6)

	// Lógica para início e fim do mês atual
	inicioMes := time.Date(hoje.Year(), hoje.Month(), 1, 0, 0, 0, 0, hoje.Location())
	fimMes := inicioMes.AddDate(0, 1, -1)

	resumo := &ResumoDashboard{}
	
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := helpers.SetSchema(ctx, tx); err != nil {
		return nil, err
	}

	totalVencido, err := s.Repo.GetTotalDebitosAtrasados(ctx, tx)
	if err != nil {
		return nil, err
	}
	resumo.TotalVencido = totalVencido

	totalSemana, err := s.Repo.GetTotalDebitosSemana(ctx, tx, inicioSemana, fimSemana)
	if err != nil {
		return nil, err
	}
	resumo.TotalSemana = totalSemana

	cats, err := s.Repo.GetDespesasPorCategoria(ctx, tx, inicioMes, fimMes)
	if err != nil {
		return nil, err
	}
	resumo.DespesasCategoria = cats

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// TODO: Quando o módulo de Vendas for criado, substituir os mocks abaixo por chamadas ao repository
	resumo.VendaDia = 0.0 // Mock provisório
	resumo.VendaMes = 0.0 // Mock provisório

	return resumo, nil
}
