package controller

import (
	"gestao/internal/auth"
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/utils/requisicao"
	"gestao/utils/resposta"
	"net/http"
)

type LoginController struct {
	service *service.Service
}

func (c *LoginController) Login(w http.ResponseWriter, r *http.Request) {

	var usuarioRequest model.UsuarioLogin

	if err := requisicao.ProcessarRequisicao(w, r, &usuarioRequest); err != nil {
		return
	}

	if err := usuarioRequest.Validar(); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	id, nome, schema, err := c.service.Login.Login(r.Context(), &usuarioRequest)
	if err != nil {
		resposta.Padrao(w, http.StatusUnauthorized, map[string]string{"erro": "E-mail ou senha inválidos"})
		return
	}

	tokenString, err := auth.GerarTokenJWT(int(id), nome, schema)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, "erro ao gerar token")
		return
	}

	token := map[string]string{
		"token": tokenString,
	}

	resposta.Padrao(w, http.StatusOK, token)

}

func (c *LoginController) LoginAdministrativo(w http.ResponseWriter, r *http.Request) {

	var usuario model.UsuarioLogin
	if err := requisicao.ProcessarRequisicao(w, r, &usuario); err != nil {
		return
	}

	if err := usuario.Validar(); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	id, nome, err := c.service.Login.LoginAdministrativo(r.Context(), &usuario)
	if err != nil {
		resposta.Padrao(w, http.StatusUnauthorized, map[string]string{"erro": "Dados de login inválidos"})
		return
	}

	tokenString, err := auth.GerarTokenJWTAdministrativo(int(id), nome)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, "erro ao gerar token")
		return
	}

	token := map[string]string{
		"token": tokenString,
	}

	resposta.Padrao(w, http.StatusOK, token)
}
