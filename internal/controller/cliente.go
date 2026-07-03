package controller

import (
	"fmt"
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"net/http"
	"strconv"
)

type ClienteController struct {
	service *service.Service
}

func (c *ClienteController) CriarCliente(w http.ResponseWriter, r *http.Request) {
	var cliente model.Cliente

	if err := requisicao.ProcessarRequisicao(w, r, &cliente); err != nil {
		fmt.Println("Erro ao processar requisição:", err)
		return
	}

	clienteCriado, err := c.service.Clientes.CriarCliente(r.Context(), &cliente)
	if err != nil {
		fmt.Println("Erro ao criar cliente:", err)
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, clienteCriado)
}

func (c *ClienteController) ListarClientes(w http.ResponseWriter, r *http.Request) {
	busca := r.URL.Query().Get("busca")

	clientes, err := c.service.Clientes.ListarClientes(r.Context(), busca)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Erro ao buscar clientes: " + err.Error()})
		return
	}

	if clientes == nil {
		clientes = []model.Cliente{}
	}

	resposta.Padrao(w, http.StatusOK, clientes)
}

func (c *ClienteController) ObterCliente(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id inválido"})
		return
	}

	cliente, err := c.service.Clientes.ObterClientePorID(r.Context(), id)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": err.Error()})
		return
	}

	if cliente == nil {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "cliente não encontrado"})
		return
	}

	resposta.Padrao(w, http.StatusOK, cliente)
}

func (c *ClienteController) AtualizarCliente(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id inválido"})
		return
	}

	var cliente model.Cliente
	if err := requisicao.ProcessarRequisicao(w, r, &cliente); err != nil {
		fmt.Println("Erro ao processar requisição:", err)
		return
	}

	err = c.service.Clientes.AtualizarCliente(r.Context(), id, &cliente)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "cliente atualizado com sucesso"})
}
