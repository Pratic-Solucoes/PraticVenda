package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
	"strconv"
)

type SaidaEstoqueController struct{ service *service.Service }

func idSaida(r *http.Request) (uint64, error) { return strconv.ParseUint(r.PathValue("id"), 10, 64) }

// RegistrarSaida cria uma saída em aberto, sem reduzir o saldo ainda.
func (c *SaidaEstoqueController) RegistrarSaida(w http.ResponseWriter, r *http.Request) {
	var s model.SaidaEstoque
	if e := requisicao.ProcessarRequisicao(w, r, &s); e != nil {
		return
	}
	uid, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, 401, map[string]string{"erro": "usuário inválido"})
		return
	}
	s.IDUsuario = int64(uid)
	if e := c.service.SaidaEstoque.RegistrarSaida(r.Context(), &s); e != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": e.Error()})
		return
	}
	resposta.Padrao(w, 201, s)
}
func (c *SaidaEstoqueController) ListarSaidas(w http.ResponseWriter, r *http.Request) {
	out, e := c.service.SaidaEstoque.ListarSaidas(r.Context(), model.FiltroSaidaEstoque{Status: r.URL.Query().Get("status"), Data: r.URL.Query().Get("data")})
	if e != nil {
		resposta.Padrao(w, 500, map[string]string{"erro": e.Error()})
		return
	}
	resposta.Padrao(w, 200, out)
}
func (c *SaidaEstoqueController) ObterSaida(w http.ResponseWriter, r *http.Request) {
	id, e := idSaida(r)
	if e != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": "id inválido"})
		return
	}
	out, e := c.service.SaidaEstoque.ObterSaida(r.Context(), id)
	if e != nil {
		resposta.Padrao(w, 404, map[string]string{"erro": "saída não encontrada"})
		return
	}
	resposta.Padrao(w, 200, out)
}

// EditarSaida atualiza itens somente enquanto a saída estiver em aberto.
func (c *SaidaEstoqueController) EditarSaida(w http.ResponseWriter, r *http.Request) {
	id, e := idSaida(r)
	if e != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": "id inválido"})
		return
	}
	var s model.SaidaEstoque
	if e = requisicao.ProcessarRequisicao(w, r, &s); e != nil {
		return
	}
	s.ID = id
	if e = c.service.SaidaEstoque.EditarSaida(r.Context(), &s); e != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": e.Error()})
		return
	}
	resposta.Padrao(w, 200, map[string]string{"mensagem": "Saída atualizada com sucesso"})
}

// AprovarSaida confirma a operação e reduz o saldo dos itens.
func (c *SaidaEstoqueController) AprovarSaida(w http.ResponseWriter, r *http.Request) {
	id, e := idSaida(r)
	if e != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": "id inválido"})
		return
	}
	uid, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, 401, map[string]string{"erro": "usuário inválido"})
		return
	}
	if e = c.service.SaidaEstoque.AprovarSaida(r.Context(), id, int64(uid)); e != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": e.Error()})
		return
	}
	resposta.Padrao(w, 200, map[string]string{"mensagem": "Saída aprovada"})
}

func (c *SaidaEstoqueController) CancelarSaida(w http.ResponseWriter, r *http.Request) {
	id, e := idSaida(r)
	if e != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": "id inválido"})
		return
	}
	uid, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, 401, map[string]string{"erro": "usuário inválido"})
		return
	}
	if e = c.service.SaidaEstoque.CancelarSaida(r.Context(), id, int64(uid)); e != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": e.Error()})
		return
	}
	resposta.Padrao(w, 200, map[string]string{"mensagem": "Saída cancelada"})
}
