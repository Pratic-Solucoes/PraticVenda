package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
)

type ConfiguracaoPDVController struct {
	service *service.ConfiguracaoPDVService
}

func (c *ConfiguracaoPDVController) Obter(w http.ResponseWriter, r *http.Request) {
	configuracao, err := c.service.Obter(r.Context())
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": err.Error()})
		return
	}
	if configuracao == nil {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "configuração de PDV não definida"})
		return
	}
	resposta.Padrao(w, http.StatusOK, configuracao)
}

func (c *ConfiguracaoPDVController) Salvar(w http.ResponseWriter, r *http.Request) {
	var configuracao model.ConfiguracaoPDV
	if requisicao.ProcessarRequisicao(w, r, &configuracao) != nil {
		return
	}
	if err := c.service.Salvar(r.Context(), configuracao); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusOK, configuracao)
}
