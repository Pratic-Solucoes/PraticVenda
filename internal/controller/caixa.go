package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
)

type CaixaController struct{ service *service.CaixaService }

func usuarioCaixa(r *http.Request) (int64, bool) {
	id, ok := r.Context().Value("usuario_id").(float64)
	return int64(id), ok
}
func (c *CaixaController) Criar(w http.ResponseWriter, r *http.Request) {
	u, ok := usuarioCaixa(r)
	if !ok {
		resposta.Padrao(w, 401, nil)
		return
	}
	var req struct {
		Nome      string `json:"nome"`
		IDUsuario int64  `json:"id_usuario"`
	}
	if requisicao.ProcessarRequisicao(w, r, &req) != nil {
		return
	}
	if req.IDUsuario > 0 {
		u = req.IDUsuario
	}
	caixa, err := c.service.Criar(r.Context(), u, req.Nome)
	if err != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 201, caixa)
}
func (c *CaixaController) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := c.service.ListarUsuarios(r.Context())
	if err != nil {
		resposta.Padrao(w, 500, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 200, usuarios)
}
func (c *CaixaController) Listar(w http.ResponseWriter, r *http.Request) {
	u, ok := usuarioCaixa(r)
	if !ok {
		resposta.Padrao(w, 401, nil)
		return
	}
	caixas, err := c.service.Listar(r.Context(), u)
	if err != nil {
		resposta.Padrao(w, 500, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 200, caixas)
}
func (c *CaixaController) Atual(w http.ResponseWriter, r *http.Request) {
	u, ok := usuarioCaixa(r)
	if !ok {
		resposta.Padrao(w, 401, nil)
		return
	}
	controle, err := c.service.Atual(r.Context(), u)
	if err != nil {
		resposta.Padrao(w, 500, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 200, map[string]any{"aberto": controle != nil, "controle": controle})
}
func (c *CaixaController) Abrir(w http.ResponseWriter, r *http.Request) {
	u, ok := usuarioCaixa(r)
	if !ok {
		resposta.Padrao(w, 401, nil)
		return
	}
	var req struct {
		IDCaixa       int64   `json:"id_caixa"`
		ValorAbertura float64 `json:"valor_abertura"`
		SenhaAcesso   string  `json:"senha_acesso"`
	}
	if requisicao.ProcessarRequisicao(w, r, &req) != nil {
		return
	}
	controle, err := c.service.Abrir(r.Context(), u, req.IDCaixa, req.ValorAbertura, req.SenhaAcesso)
	if err != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 201, controle)
}
func (c *CaixaController) Fechar(w http.ResponseWriter, r *http.Request) {
	u, ok := usuarioCaixa(r)
	if !ok {
		resposta.Padrao(w, 401, nil)
		return
	}
	var req model.FechamentoCaixa
	if requisicao.ProcessarRequisicao(w, r, &req) != nil {
		return
	}
	if err := c.service.Fechar(r.Context(), u, req); err != nil {
		resposta.Padrao(w, 400, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, 200, nil)
}
