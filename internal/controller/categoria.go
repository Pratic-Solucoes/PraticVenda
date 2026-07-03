package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"net/http"
)

type CategoriaController struct {
	service *service.Service
}

func (c *CategoriaController) CriarCategoria(w http.ResponseWriter, r *http.Request) {
	var categoria model.CategoriaDebito
	if err := requisicao.ProcessarRequisicao(w, r, &categoria); err != nil {
		return
	}

	categoriaCriada, err := c.service.Categorias.CriarCategoria(r.Context(), &categoria)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, categoriaCriada)
}

func (c *CategoriaController) ListarCategorias(w http.ResponseWriter, r *http.Request) {
	categorias, err := c.service.Categorias.ListarCategorias(r.Context())
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao buscar categorias: " + err.Error()})
		return
	}

	if categorias == nil {
		categorias = []*model.CategoriaDebito{}
	}

	resposta.Padrao(w, http.StatusOK, categorias)
}
