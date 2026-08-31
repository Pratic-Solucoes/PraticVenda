package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type VendaPDVController struct{ service *service.VendaPDVService }

func (c *VendaPDVController) Finalizar(w http.ResponseWriter, r *http.Request) {
	var v model.VendaPDV
	if e := requisicao.ProcessarRequisicao(w, r, &v); e != nil {
		return
	}
	uid, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, 401, map[string]string{"erro": "usuário inválido"})
		return
	}
	id, e := c.service.Finalizar(r.Context(), v, int64(uid))
	if e != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": e.Error()})
		return
	}
	resposta.Padrao(w, 201, map[string]any{"mensagem": "venda concluída", "id": id})
}

func (c *VendaPDVController) SalvarPreVenda(w http.ResponseWriter, r *http.Request) {
	var v model.VendaPDV
	if requisicao.ProcessarRequisicao(w, r, &v) != nil {
		return
	}
	uid, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, 401, map[string]string{"erro": "usuário inválido"})
		return
	}
	id, err := c.service.SalvarPreVenda(r.Context(), v, int64(uid))
	if err != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 201, map[string]any{"mensagem": "pré-venda salva", "id": id})
}
func (c *VendaPDVController) ListarPreVendas(w http.ResponseWriter, r *http.Request) {
	itens, err := c.service.ListarPreVendas(r.Context())
	if err != nil {
		resposta.Padrao(w, 500, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 200, itens)
}
func (c *VendaPDVController) ObterPreVenda(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		resposta.Padrao(w, 400, map[string]string{"erro": "id inválido"})
		return
	}
	v, err := c.service.ObterPreVenda(r.Context(), id)
	if err != nil {
		resposta.Padrao(w, 404, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 200, v)
}
func (c *VendaPDVController) Cancelar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		resposta.Padrao(w, 400, map[string]string{"erro": "id inválido"})
		return
	}
	uid, ok := r.Context().Value("usuario_id").(float64)
	if !ok {
		resposta.Padrao(w, 401, map[string]string{"erro": "usuário inválido"})
		return
	}
	if err = c.service.Cancelar(r.Context(), id, int64(uid)); err != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 200, map[string]string{"mensagem": "venda cancelada e estornada"})
}
