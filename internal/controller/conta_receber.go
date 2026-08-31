package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type ContaReceberController struct{ service *service.Service }

func (c *ContaReceberController) CriarContaReceber(w http.ResponseWriter, r *http.Request) {
	var conta model.ContaReceberCriar
	if err := requisicao.ProcessarRequisicao(w, r, &conta); err != nil {
		return
	}
	if err := c.service.ContasReceber.CriarContaReceber(r.Context(), &conta); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusCreated, nil)
}
func (c *ContaReceberController) ListarContasReceber(w http.ResponseWriter, r *http.Request) {
	resultado, err := c.service.ContasReceber.ListarContasReceber(r.Context(), r.URL.Query().Get("busca"), r.URL.Query().Get("vencimento"), r.URL.Query().Get("status"))
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": err.Error()})
		return
	}
	if resultado == nil {
		resultado = []*model.ContaReceber{}
	}
	resposta.Padrao(w, http.StatusOK, resultado)
}
func (c *ContaReceberController) ReceberConta(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "ID inválido"})
		return
	}
	var entrada struct {
		ValorRecebimento     float64 `json:"valor_recebimento"`
		IDFormaPagamentoReal int64   `json:"id_forma_pagamento_real"`
	}
	if err := requisicao.ProcessarRequisicao(w, r, &entrada); err != nil {
		return
	}
	if err := c.service.ContasReceber.ReceberConta(r.Context(), id, entrada.ValorRecebimento, entrada.IDFormaPagamentoReal); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusOK, nil)
}
