package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
)

type FornecedorController struct {
	service *service.Service
}

func (c *FornecedorController) ListarFornecedores(w http.ResponseWriter, r *http.Request) {
	busca := r.URL.Query().Get("busca")
	
	fornecedores, err := c.service.Fornecedores.ListarFornecedores(r.Context(), busca)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao buscar fornecedores: " + err.Error()})
		return
	}

	// Retorna array vazio em vez de null se não houver fornecedores
	if fornecedores == nil {
		fornecedores = []*model.Fornecedor{}
	}

	resposta.Padrao(w, http.StatusOK, fornecedores)
}

func (c *FornecedorController) CriarFornecedor(w http.ResponseWriter, r *http.Request) {

	var fornecedor model.Fornecedor
	if err := requisicao.ProcessarRequisicao(w, r, &fornecedor); err != nil {
		// ProcessarRequisicao já envia o 400 se o JSON for inválido
		return
	}

	fornecedorCriado, err := c.service.Fornecedores.CriarFornecedor(r.Context(), &fornecedor)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, fornecedorCriado)
}

func (c *FornecedorController) ObterFornecedor(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id inválido"})
		return
	}

	fornecedor, err := c.service.Fornecedores.ObterFornecedorPorID(r.Context(), id)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": err.Error()})
		return
	}

	if fornecedor == nil {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "fornecedor não encontrado"})
		return
	}

	resposta.Padrao(w, http.StatusOK, fornecedor)
}

func (c *FornecedorController) AtualizarFornecedor(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id inválido"})
		return
	}

	var fornecedor model.Fornecedor
	if err := requisicao.ProcessarRequisicao(w, r, &fornecedor); err != nil {
		return
	}

	err = c.service.Fornecedores.AtualizarFornecedor(r.Context(), id, &fornecedor)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "fornecedor atualizado com sucesso"})
}
