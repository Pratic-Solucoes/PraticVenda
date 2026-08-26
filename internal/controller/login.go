package controller

import (
	"fmt"
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
		resposta.Padrao(w, http.StatusBadRequest, "erro ao ler json")
		return
	}

	id, nome, schema, err := c.service.Login.Login(r.Context(), &usuarioRequest)
	if err != nil {
		fmt.Printf("ERRO REAL DURANTE O LOGIN: %v\n", err) // Print para debug
		resposta.Padrao(w, http.StatusUnauthorized, "dados login inválidos")
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
