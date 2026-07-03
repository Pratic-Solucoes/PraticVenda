package controller

import (
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"net/http"
)

type UsuarioController struct {
	service *service.Service
}

func (c *UsuarioController) CriarUsuario(w http.ResponseWriter, r *http.Request) {

	var usuarioCriar model.UsuarioCriar
	if err := requisicao.ProcessarRequisicao(w, r, &usuarioCriar); err != nil {
		// ProcessarRequisicao já enviou um HTTP 400 caso tenha falhado
		return
	}

	usuarioCriado, err := c.service.Usuarios.CriarUsuario(r.Context(), &usuarioCriar)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, usuarioCriado)
}

func (c *UsuarioController) BuscarUsuarioPorID(w http.ResponseWriter, r *http.Request) {

	// O ID do usuário vem do token JWT, que é validado pelo middleware de autenticação.
	// O pacote JWT decodifica números JSON como float64, então tratamos esse tipo.
	usuarioIDClaim := r.Context().Value("usuario_id")
	if usuarioIDClaim == nil {
		resposta.Padrao(w, http.StatusUnauthorized, map[string]string{"erro": "ID do usuário não encontrado no token"})
		return
	}

	// Convertemos o ID de float64 para int de forma segura.
	usuarioID, ok := usuarioIDClaim.(float64)
	if !ok {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Formato de ID do usuário inválido no token"})
		return
	}

	usuario, err := c.service.Usuarios.BuscarUsuarioPorID(r.Context(), int(usuarioID))
	if err != nil {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "erro ao buscar usuário: " + err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusOK, usuario)
}

func (c *UsuarioController) EditarUsuario(w http.ResponseWriter, r *http.Request) {}

func (c *UsuarioController) AlterarSenha(w http.ResponseWriter, r *http.Request) {

	type payload struct {
		SenhaAtual       string `json:"senha_atual"`
		NovaSenha        string `json:"nova_senha"`
		SenhaConfirmacao string `json:"senha_confirmacao"`
	}

	var p payload

	if err := requisicao.ProcessarRequisicao(w, r, &p); err != nil {
		return
	}

	usuarioIDClaim := r.Context().Value("usuario_id")
	if usuarioIDClaim == nil {
		resposta.Padrao(w, http.StatusUnauthorized, map[string]string{"erro": "ID do usuário não encontrado no token"})
		return
	}

	if err := c.service.Usuarios.AlterarSenha(r.Context(), int64(usuarioIDClaim.(float64)), p.SenhaAtual, p.NovaSenha, p.SenhaConfirmacao); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusOK, nil)
}
