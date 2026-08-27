package controller

import (
	"gestao/internal/service"
	"gestao/utils/resposta"
	"net/http"
)

type AdminController struct {
	service *service.Service
}

func (c *AdminController) CarregarOrganizacoes(w http.ResponseWriter, r *http.Request) {
	organizacoes, err := c.service.Admin.CarregarOrganizacoes(r.Context())
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao carregar organizações"})
		return
	}

	resposta.Padrao(w, http.StatusOK, organizacoes)
}

func (c *AdminController) CarregarUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := c.service.Admin.CarregarUsuarios(r.Context())
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao carregar usuários"})
		return
	}

	resposta.Padrao(w, http.StatusOK, usuarios)
}

func (c *AdminController) CriarOrganizacao(w http.ResponseWriter, r *http.Request) {}

func (c *AdminController) CriarUsuario(w http.ResponseWriter, r *http.Request) {}
