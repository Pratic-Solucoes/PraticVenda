package controller

import (
	"fmt"
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	idParam := chi.URLParam(r, "id")
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
	idParam := chi.URLParam(r, "id")
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

func (c *ClienteController) CriarEndereco(w http.ResponseWriter, r *http.Request) {
	var idCliente int64
	var endereco model.EnderecoCliente

	idCliente, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do cliente inválido"})
		return
	}

	if err := requisicao.ProcessarRequisicao(w, r, &endereco); err != nil {
		return
	}

	enderecoCriado, err := c.service.Clientes.CriarEndereco(r.Context(), idCliente, &endereco)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, enderecoCriado)
}

func (c *ClienteController) BuscarEnderecoByID(w http.ResponseWriter, r *http.Request) {

	var idCliente, idEndereco int64

	idCliente, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do cliente inválido"})
		return
	}

	idEndereco, err = strconv.ParseInt(chi.URLParam(r, "id_endereco"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do endereço inválido"})
		return
	}

	endereco, err := c.service.Clientes.BuscarEnderecoByID(r.Context(), idCliente, idEndereco)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": err.Error()})
		return
	}

	if endereco == nil {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "endereço não encontrado"})
		return
	}

	resposta.Padrao(w, http.StatusOK, endereco)
}

func (c *ClienteController) AtualizarEndereco(w http.ResponseWriter, r *http.Request) {

	var idCliente, idEndereco int64
	var endereco model.EnderecoCliente

	idCliente, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do cliente inválido"})
		return
	}

	idEndereco, err = strconv.ParseInt(chi.URLParam(r, "id_endereco"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do endereço inválido"})
		return
	}

	if err := requisicao.ProcessarRequisicao(w, r, &endereco); err != nil {
		return
	}

	err = c.service.Clientes.EditarEndereco(r.Context(), idCliente, idEndereco, &endereco)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "endereço atualizado com sucesso"})
}

func (c *ClienteController) CriarTelefone(w http.ResponseWriter, r *http.Request) {

	var idCliente int64
	var telefone model.TelefoneCliente

	idCliente, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do cliente inválido"})
		return
	}

	if err := requisicao.ProcessarRequisicao(w, r, &telefone); err != nil {
		return
	}

	telefoneCriado, err := c.service.Clientes.CriarTelefone(r.Context(), idCliente, &telefone)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, telefoneCriado)
}

func (c *ClienteController) BuscarTelefoneByID(w http.ResponseWriter, r *http.Request) {
	var idCliente, idTelefone int64

	idCliente, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do cliente inválido"})
		return
	}

	idTelefone, err = strconv.ParseInt(chi.URLParam(r, "id_telefone"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do telefone inválido"})
		return
	}

	telefone, err := c.service.Clientes.BuscarTelefoneByID(r.Context(), idCliente, idTelefone)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": err.Error()})
		return
	}

	if telefone == nil {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "telefone não encontrado"})
		return
	}

	resposta.Padrao(w, http.StatusOK, telefone)
}

func (c *ClienteController) AtualizarTelefone(w http.ResponseWriter, r *http.Request) {
	var idCliente, idTelefone int64
	var telefone model.TelefoneCliente

	idCliente, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do cliente inválido"})
		return
	}

	idTelefone, err = strconv.ParseInt(chi.URLParam(r, "id_telefone"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id do telefone inválido"})
		return
	}

	if err := requisicao.ProcessarRequisicao(w, r, &telefone); err != nil {
		return
	}

	err = c.service.Clientes.EditarTelefone(r.Context(), idCliente, idTelefone, &telefone)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "telefone atualizado com sucesso"})
}
