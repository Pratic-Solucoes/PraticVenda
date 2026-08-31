package controller

import (
	"database/sql"
	"errors"
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CondicaoPagamentoController struct {
	service *service.Service
}

func (c *CondicaoPagamentoController) Criar(w http.ResponseWriter, r *http.Request) {
	var condicao model.CondicaoPagamento
	if err := requisicao.ProcessarRequisicao(w, r, &condicao); err != nil {
		return
	}

	if err := c.service.CondicoesPagamento.Criar(r.Context(), &condicao); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, condicao)
}

func (c *CondicaoPagamentoController) Listar(w http.ResponseWriter, r *http.Request) {
	condicoes, err := c.service.CondicoesPagamento.Listar(r.Context())
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao listar condições de pagamento"})
		return
	}
	resposta.Padrao(w, http.StatusOK, condicoes)
}

func (c *CondicaoPagamentoController) BuscarPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de condição de pagamento inválido"})
		return
	}

	condicao, err := c.service.CondicoesPagamento.BuscarPorID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "condição de pagamento não encontrada"})
		return
	}
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao buscar condição de pagamento"})
		return
	}
	resposta.Padrao(w, http.StatusOK, condicao)
}

func (c *CondicaoPagamentoController) Atualizar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de condição de pagamento inválido"})
		return
	}

	var condicao model.CondicaoPagamento
	if err := requisicao.ProcessarRequisicao(w, r, &condicao); err != nil {
		return
	}
	condicao.ID = id

	if err := c.service.CondicoesPagamento.Atualizar(r.Context(), &condicao); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "condição de pagamento não encontrada"})
			return
		}
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusOK, condicao)
}
