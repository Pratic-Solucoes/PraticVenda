package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"net/http"
	"strconv"
)

type EntradaEstoqueController struct {
	service *service.Service
}

func (c *EstoqueController) EntradaEstoque(w http.ResponseWriter, r *http.Request) {
	idEstoque, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de estoque inválido"})
		return
	}

	idUsuario, ok := r.Context().Value("usuario_id").(int64)
	if !ok {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de usuario inválido"})
		return
	}

	var entrada model.EntradaEstoque

	entrada.IDEstoque = uint64(idEstoque)
	entrada.IDUsuario = idUsuario

	if err := requisicao.ProcessarRequisicao(w, r, &entrada); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	if err := c.service.EntradaEstoque.RegistrarEntrada(r.Context(), &entrada); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, map[string]string{"mensagem": "Entrada de estoque registrada com sucesso"})
}
