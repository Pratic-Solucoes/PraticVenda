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

type ContaPagarController struct {
	service *service.Service
}

func (c *ContaPagarController) CriarContaPagar(w http.ResponseWriter, r *http.Request) {
	var contaPagar model.ContaPagarCriar
	if err := requisicao.ProcessarRequisicao(w, r, &contaPagar); err != nil {
		return
	}

	if err := c.service.ContasPagar.CriarContaPagar(r.Context(), &contaPagar); err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao processar conta a pagar: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, nil)
}

func (c *ContaPagarController) ListarContasPagar(w http.ResponseWriter, r *http.Request) {
	busca := r.URL.Query().Get("busca")
	vencimento := r.URL.Query().Get("vencimento")
	status := r.URL.Query().Get("status")

	contasPagar, err := c.service.ContasPagar.ListarContasPagar(r.Context(), busca, vencimento, status)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao buscar contas a pagar: " + err.Error()})
		return
	}

	if contasPagar == nil {
		contasPagar = []*model.ContaPagar{}
	}

	resposta.Padrao(w, http.StatusOK, contasPagar)
}

func (c *ContaPagarController) PagarContaPagar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "ID inválido"})
		return
	}

	var req struct {
		ValorPagamento float64 `json:"valor_pagamento"`
	}
	if err := requisicao.ProcessarRequisicao(w, r, &req); err != nil {
		return
	}

	if err := c.service.ContasPagar.PagarContaPagar(r.Context(), id, req.ValorPagamento); err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao pagar conta a pagar: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, nil)
}

func (c *ContaPagarController) EditarContaPagar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "ID inválido"})
		return
	}

	var contaPagar model.ContaPagarCriar
	if err := requisicao.ProcessarRequisicao(w, r, &contaPagar); err != nil {
		return
	}

	if err := c.service.ContasPagar.EditarContaPagar(r.Context(), id, &contaPagar); err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao editar conta a pagar: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, nil)
}
