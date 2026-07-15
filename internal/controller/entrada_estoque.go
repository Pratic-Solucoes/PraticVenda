package controller

import (
	"encoding/json"
	"gestao/internal/model"
	"gestao/internal/service"
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

	idUsuarioFloat, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, http.StatusUnauthorized, map[string]string{"erro": "id de usuario inválido no token"})
		return
	}
	idUsuario := int64(idUsuarioFloat)

	var entrada model.EntradaEstoque

	entrada.IDEstoque = uint64(idEstoque)
	entrada.IDUsuario = idUsuario

	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "Erro ao decodificar dados: " + err.Error()})
		return
	}

	if err := c.service.EntradaEstoque.RegistrarEntrada(r.Context(), &entrada); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, map[string]string{"mensagem": "Entrada de estoque registrada com sucesso"})
}
