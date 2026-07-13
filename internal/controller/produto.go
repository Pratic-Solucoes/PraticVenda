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

type ProdutoController struct {
	service *service.Service
}

func (c *ProdutoController) CriarProduto(w http.ResponseWriter, r *http.Request) {
	var input model.ProdutoInput

	if err := requisicao.ProcessarRequisicao(w, r, &input); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "dados inválidos: " + err.Error()})
		return
	}

	produto, err := c.service.Produtos.CriarProduto(r.Context(), &input)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao criar produto: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, produto)
}

func (c *ProdutoController) ListarProdutos(w http.ResponseWriter, r *http.Request) {
	busca := r.URL.Query().Get("busca")

	produtos, err := c.service.Produtos.ListarProdutos(r.Context(), busca)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao listar produtos: " + err.Error()})
		return
	}

	if produtos == nil {
		produtos = []*model.ProdutoCompleto{}
	}

	resposta.Padrao(w, http.StatusOK, produtos)
}

func (c *ProdutoController) ObterProduto(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de produto inválido"})
		return
	}

	produto, err := c.service.Produtos.ObterProdutoPorID(r.Context(), id)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao obter produto: " + err.Error()})
		return
	}

	if produto == nil {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "produto não encontrado"})
		return
	}

	resposta.Padrao(w, http.StatusOK, produto)
}

func (c *ProdutoController) AtualizarProduto(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de produto inválido"})
		return
	}

	var input model.ProdutoInput
	if err := requisicao.ProcessarRequisicao(w, r, &input); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "dados inválidos: " + err.Error()})
		return
	}

	err = c.service.Produtos.AtualizarProduto(r.Context(), id, &input)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao atualizar produto: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "produto atualizado com sucesso"})
}

func (c *ProdutoController) ExcluirProduto(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de produto inválido"})
		return
	}

	status, err := c.service.Produtos.ExcluirOuInativarProduto(r.Context(), id)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao excluir/inativar produto: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "produto processado com sucesso", "status": status})
}

func (c *ProdutoController) VincularEstoque(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	idProduto, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de produto inválido"})
		return
	}

	var input struct {
		IDEstoque     int64   `json:"id_estoque"`
		EstoqueMinimo float64 `json:"estoque_minimo"`
	}

	if err := requisicao.ProcessarRequisicao(w, r, &input); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "dados inválidos: " + err.Error()})
		return
	}

	err = c.service.Produtos.VincularProdutoEstoque(r.Context(), idProduto, input.IDEstoque, input.EstoqueMinimo)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao vincular produto ao estoque: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "produto vinculado ao estoque com sucesso"})
}

func (c *ProdutoController) DesvincularEstoque(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	idProduto, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de produto inválido"})
		return
	}

	idEstoqueParam := chi.URLParam(r, "id_estoque")
	idEstoque, err := strconv.ParseInt(idEstoqueParam, 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de estoque inválido"})
		return
	}

	err = c.service.Produtos.DesvincularProdutoEstoque(r.Context(), idProduto, idEstoque)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao desvincular produto do estoque: " + err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "produto desvinculado do estoque com sucesso"})
}

func (c *ProdutoController) ListarGruposTributarios(w http.ResponseWriter, r *http.Request) {
	grupos, err := c.service.Produtos.ListarGruposTributarios(r.Context())
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao listar grupos tributários: " + err.Error()})
		return
	}

	if grupos == nil {
		grupos = []*model.GrupoTributario{}
	}

	resposta.Padrao(w, http.StatusOK, grupos)
}
