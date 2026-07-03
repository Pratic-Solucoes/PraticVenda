package repository

import (
	"context"
	"database/sql"
	"time"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// GetTotalDebitosAtrasados retorna o valor total de débitos vencidos e não pagos
func (r *DashboardRepository) GetTotalDebitosAtrasados(ctx context.Context, tx *sql.Tx) (float64, error) {
	query := `
		SELECT COALESCE(SUM(valor), 0)
		FROM tb_debitos
		WHERE status = 'PENDENTE' AND dt_vencimento < CURRENT_DATE
	`
	var total float64
	err := tx.QueryRowContext(ctx, query).Scan(&total)
	return total, err
}

// GetTotalDebitosSemana retorna o valor total de débitos que vencem na semana (Segunda a Domingo)
func (r *DashboardRepository) GetTotalDebitosSemana(ctx context.Context, tx *sql.Tx, inicioSemana, fimSemana time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(valor), 0)
		FROM tb_debitos
		WHERE status = 'PENDENTE'
		  AND dt_vencimento >= $1
		  AND dt_vencimento <= $2
	`
	var total float64
	err := tx.QueryRowContext(ctx, query, inicioSemana.Format("2006-01-02"), fimSemana.Format("2006-01-02")).Scan(&total)
	return total, err
}

type CategoriaGasto struct {
	Categoria string  `json:"categoria"`
	Total     float64 `json:"total"`
}

// GetDespesasPorCategoria retorna as categorias que mais geraram gastos no mês atual
func (r *DashboardRepository) GetDespesasPorCategoria(ctx context.Context, tx *sql.Tx, inicioMes, fimMes time.Time) ([]CategoriaGasto, error) {
	query := `
		SELECT
			COALESCE(c.nome, 'Sem Categoria') as categoria,
			SUM(d.valor) as total
		FROM tb_debitos d
		LEFT JOIN tb_categorias_debito c ON d.id_categoria = c.id
		WHERE d.dt_vencimento >= $1 AND d.dt_vencimento <= $2
		GROUP BY c.nome
		ORDER BY total DESC
		LIMIT 5
	`
	rows, err := tx.QueryContext(ctx, query, inicioMes.Format("2006-01-02"), fimMes.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categorias []CategoriaGasto
	for rows.Next() {
		var c CategoriaGasto
		if err := rows.Scan(&c.Categoria, &c.Total); err != nil {
			return nil, err
		}
		categorias = append(categorias, c)
	}
	return categorias, nil
}
