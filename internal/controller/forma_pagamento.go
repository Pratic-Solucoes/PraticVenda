package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
	"strconv"
)

type FormaPagamentoController struct {
	service *service.Service
}

func (c *FormaPagamentoController) Criar(w http.ResponseWriter, r *http.Request) {

	var p model.FormaPagamento

	if err := requisicao.ProcessarRequisicao(w, r, &p); err != nil {
		return
	}

	fp, err := c.service.FormasPagamento.Criar(r.Context(), &p)
	if err != nil {
		http.Error(w, "erro ao criar forma de pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resposta.Padrao(w, http.StatusCreated, fp)

}

func (c *FormaPagamentoController) Listar(w http.ResponseWriter, r *http.Request) {

	fp, err := c.service.FormasPagamento.Listar(r.Context())
	if err != nil {
		http.Error(w, "erro ao listar formas de pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resposta.Padrao(w, http.StatusOK, fp)

}

func (c *FormaPagamentoController) BuscarPorID(w http.ResponseWriter, r *http.Request) {

	idFp, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "erro ao buscar forma de pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fp, err := c.service.FormasPagamento.BuscarPorID(r.Context(), int64(idFp))
	if err != nil {
		http.Error(w, "erro ao buscar forma de pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resposta.Padrao(w, http.StatusOK, fp)

}

func (c *FormaPagamentoController) Atualizar(w http.ResponseWriter, r *http.Request) {
	idFp, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "erro ao capturar ID da forma de pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var fp model.FormaPagamento
	if err := requisicao.ProcessarRequisicao(w, r, &fp); err != nil {
		return
	}

	fp.ID = uint64(idFp)

	fpAtualizado, err := c.service.FormasPagamento.Atualizar(r.Context(), &fp)
	if err != nil {
		http.Error(w, "erro ao atualizar forma de pagamento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resposta.Padrao(w, http.StatusOK, fpAtualizado)
}
