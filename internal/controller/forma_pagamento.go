package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"net/http"
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
