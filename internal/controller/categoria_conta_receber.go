package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
)

type CategoriaContaReceberController struct{ service *service.Service }

func (c *CategoriaContaReceberController) CriarCategoria(w http.ResponseWriter, r *http.Request) {
	var categoria model.CategoriaContaReceber
	if err := requisicao.ProcessarRequisicao(w, r, &categoria); err != nil {
		return
	}
	resultado, err := c.service.CategoriasContasReceber.CriarCategoria(r.Context(), &categoria)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusCreated, resultado)
}
func (c *CategoriaContaReceberController) ListarCategorias(w http.ResponseWriter, r *http.Request) {
	resultado, err := c.service.CategoriasContasReceber.ListarCategorias(r.Context())
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": err.Error()})
		return
	}
	if resultado == nil {
		resultado = []model.CategoriaContaReceber{}
	}
	resposta.Padrao(w, http.StatusOK, resultado)
}
