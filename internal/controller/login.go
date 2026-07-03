package controller

import (
	"fmt"
	"gestao/internal/model"
	"gestao/internal/service"
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"gestao/pkg/token"
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

	tokenString, err := token.GerarTokenJWT(int(id), nome, schema)
	if err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, "erro ao gerar token")
		return
	}

	token := map[string]string{
		"token": tokenString,
	}

	resposta.Padrao(w, http.StatusOK, token)

}
