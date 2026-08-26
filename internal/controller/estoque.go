package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type EstoqueController struct {
	service *service.Service
}

func (c *EstoqueController) CriarEstoque(w http.ResponseWriter, r *http.Request) {
	var input model.EstoqueCriar

	if err := requisicao.ProcessarRequisicao(w, r, &input); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "dados inválidos: " + err.Error()})
		return
	}

	estoque, err := c.service.Estoques.CriarEstoque(r.Context(), &input)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao criar estoque: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, estoque)
}

func (c *EstoqueController) ListarEstoques(w http.ResponseWriter, r *http.Request) {
	estoques, err := c.service.Estoques.ListarEstoques(r.Context())
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao listar estoques: " + err.Error()})
		return
	}

	if estoques == nil {
		estoques = []*model.Estoque{}
	}

	resposta.Padrao(w, http.StatusOK, estoques)
}

func (c *EstoqueController) ListarProdutosDoEstoque(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	idEstoque, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de estoque inválido"})
		return
	}

	produtos, err := c.service.Estoques.ListarProdutosDoEstoque(r.Context(), idEstoque)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao listar produtos do estoque: " + err.Error()})
		return
	}

	if produtos == nil {
		produtos = []*model.ProdutoEstoque{}
	}

	resposta.Padrao(w, http.StatusOK, produtos)
}
