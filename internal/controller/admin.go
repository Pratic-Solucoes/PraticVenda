package controller

import (
	"gestao/internal/service"
	"net/http"
)

type AdminController struct {
	service *service.Service
}

func (c *AdminController) CarregarOrganizacoes(w http.ResponseWriter, r *http.Request) {}

func (c *AdminController) CarregarUsuarios(w http.ResponseWriter, r *http.Request) {}

func (c *AdminController) CriarOrganizacao(w http.ResponseWriter, r *http.Request) {}

func (c *AdminController) CriarUsuario(w http.ResponseWriter, r *http.Request) {}
