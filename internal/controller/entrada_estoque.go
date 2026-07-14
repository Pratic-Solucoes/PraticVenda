package controller

import (
	"gestao/internal/service"
	"gestao/pkg/resposta"
	"net/http"
	"strconv"
)

type EntradaEstoqueController struct {
	service *service.Service
}

func (c *EstoqueController) EntradaEstoque(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de estoque inválido"})
		return
	}

	// TODO: implementar lógica de entrada de estoque
	resposta.Padrao(w, http.StatusNotImplemented, map[string]string{"erro": "endpoint ainda não implementado"})
}
