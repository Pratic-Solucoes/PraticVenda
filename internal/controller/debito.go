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

type DebitoController struct {
	service *service.Service
}

func (c *DebitoController) CriarDebitoAvulso(w http.ResponseWriter, r *http.Request) {
	var debito model.DebitoAvulsoCriar
	if err := requisicao.ProcessarRequisicao(w, r, &debito); err != nil {
		return
	}

	if err := c.service.Debitos.LancarDebito(r.Context(), &debito); err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao processar débito: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, nil)
}

func (c *DebitoController) ListarDebitos(w http.ResponseWriter, r *http.Request) {
	busca := r.URL.Query().Get("busca")
	vencimento := r.URL.Query().Get("vencimento")
	status := r.URL.Query().Get("status")

	debitos, err := c.service.Debitos.ListarDebitos(r.Context(), busca, vencimento, status)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao buscar débitos: " + err.Error()})
		return
	}

	if debitos == nil {
		debitos = []*model.Debito{}
	}

	resposta.Padrao(w, http.StatusOK, debitos)
}

func (c *DebitoController) PagarDebito(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "ID inválido"})
		return
	}

	if err := c.service.Debitos.PagarDebito(r.Context(), id); err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao pagar débito: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, nil)
}

func (c *DebitoController) EditarDebito(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "ID inválido"})
		return
	}

	var debito model.DebitoAvulsoCriar
	if err := requisicao.ProcessarRequisicao(w, r, &debito); err != nil {
		return
	}

	if err := c.service.Debitos.EditarDebito(r.Context(), id, &debito); err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao editar débito: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, nil)
}
