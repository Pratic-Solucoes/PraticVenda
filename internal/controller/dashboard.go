package controller

import (
	"gestao/internal/service"
	"gestao/pkg/resposta"
	"net/http"
)

type DashboardController struct {
	service *service.DashboardService
}

func NewDashboardController(s *service.DashboardService) *DashboardController {
	return &DashboardController{service: s}
}

func (c *DashboardController) ResumoDashboard(w http.ResponseWriter, r *http.Request) {
	resumo, err := c.service.ObterResumo(r.Context())
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao buscar resumo do dashboard: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, resumo)
}
