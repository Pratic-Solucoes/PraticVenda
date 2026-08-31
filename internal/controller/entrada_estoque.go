package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
	"strconv"
)

type EntradaEstoqueController struct {
	service *service.Service
}

func (c *EntradaEstoqueController) RegistrarEntrada(w http.ResponseWriter, r *http.Request) {
	idUsuarioFloat, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, http.StatusUnauthorized, map[string]string{"erro": "id de usuario inválido no token"})
		return
	}
	idUsuario := int64(idUsuarioFloat)

	var entrada model.EntradaEstoque

	entrada.IDUsuario = idUsuario

	if err := requisicao.ProcessarRequisicao(w, r, &entrada); err != nil {
		return
	}
	if entrada.IDEstoque == 0 {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de estoque é obrigatório"})
		return
	}

	if err := c.service.EntradaEstoque.RegistrarEntrada(r.Context(), &entrada); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, map[string]string{"mensagem": "Entrada de estoque registrada com sucesso"})
}

func (c *EntradaEstoqueController) ListarEntradas(w http.ResponseWriter, r *http.Request) {
	filtro := model.FiltroEntradaEstoque{
		Fornecedor: r.URL.Query().Get("fornecedor"),
		Data:       r.URL.Query().Get("data"),
		Status:     r.URL.Query().Get("status"),
	}

	if id := r.URL.Query().Get("id"); id != "" {
		entradaID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de entrada inválido"})
			return
		}
		filtro.ID = entradaID
	}

	entradas, err := c.service.EntradaEstoque.ListarEntradas(r.Context(), filtro)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "erro ao listar entradas de estoque"})
		return
	}

	resposta.Padrao(w, http.StatusOK, entradas)
}

// ObterEntrada retorna a entrada e seus itens para visualização ou edição.
func (c *EntradaEstoqueController) ObterEntrada(w http.ResponseWriter, r *http.Request) {
	idEntrada, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de entrada inválido"})
		return
	}
	entrada, err := c.service.EntradaEstoque.ObterEntrada(r.Context(), idEntrada)
	if err != nil {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "entrada de estoque não encontrada"})
		return
	}
	resposta.Padrao(w, http.StatusOK, entrada)
}

// EditarEntrada altera uma entrada somente enquanto ela estiver aberta.
func (c *EntradaEstoqueController) EditarEntrada(w http.ResponseWriter, r *http.Request) {
	idEntrada, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de entrada inválido"})
		return
	}
	var entrada model.EntradaEstoque
	if err := requisicao.ProcessarRequisicao(w, r, &entrada); err != nil {
		return
	}
	entrada.ID = idEntrada
	if err := c.service.EntradaEstoque.EditarEntrada(r.Context(), &entrada); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "Entrada de estoque atualizada com sucesso"})
}

// AprovarEntrada conclui uma entrada aberta e aplica seus itens ao estoque.
func (c *EntradaEstoqueController) AprovarEntrada(w http.ResponseWriter, r *http.Request) {
	idEntrada, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": "id de entrada inválido"})
		return
	}
	usuarioID, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, http.StatusUnauthorized, map[string]string{"erro": "id de usuário inválido no token"})
		return
	}
	if err := c.service.EntradaEstoque.AprovarEntrada(r.Context(), idEntrada, int64(usuarioID)); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusOK, map[string]string{"mensagem": "Entrada de estoque aprovada com sucesso"})
}
func (c *EntradaEstoqueController) CancelarEntrada(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": "id de entrada inválido"})
		return
	}
	uid, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, 401, map[string]string{"erro": "usuário inválido"})
		return
	}
	if err = c.service.EntradaEstoque.CancelarEntrada(r.Context(), id, int64(uid)); err != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 200, map[string]string{"mensagem": "Entrada cancelada"})
}
