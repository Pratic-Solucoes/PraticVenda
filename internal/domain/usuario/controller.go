package usuario

import (
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"net/http"
)

type ControllerInterface interface {
	CriarUsuario(w http.ResponseWriter, r *http.Request)
	BuscarUsuarioPorID(w http.ResponseWriter, r *http.Request)
	EditarUsuario(w http.ResponseWriter, r *http.Request)
	AlterarSenha(w http.ResponseWriter, r *http.Request)
}

type Controller struct {
	service ServiceInterface
}

func NewController(service ServiceInterface) ControllerInterface {
	return &Controller{
		service: service,
	}
}

func (c *Controller) CriarUsuario(w http.ResponseWriter, r *http.Request) {

	var usuarioCriar Usuario
	if err := requisicao.ProcessarRequisicao(w, r, &usuarioCriar); err != nil {
		// ProcessarRequisicao já enviou um HTTP 400 caso tenha falhado
		return
	}

	usuarioCriado, err := c.service.CriarUsuario(r.Context(), &usuarioCriar)
	if err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	resposta.Padrao(w, http.StatusCreated, usuarioCriado)
}

func (c *Controller) BuscarUsuarioPorID(w http.ResponseWriter, r *http.Request) {

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

	usuario, err := c.service.BuscarUsuarioPorID(r.Context(), uint64(usuarioID))
	if err != nil {
		resposta.Padrao(w, http.StatusNotFound, map[string]string{"erro": "erro ao buscar usuário: " + err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusOK, usuario)
}

func (c *Controller) EditarUsuario(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) AlterarSenha(w http.ResponseWriter, r *http.Request) {

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

	usuarioID, ok := usuarioIDClaim.(float64)
	if !ok {
		resposta.Padrao(w, http.StatusInternalServerError, map[string]string{"erro": "Formato de ID do usuário inválido no token"})
		return
	}

	if err := c.service.AlterarSenha(r.Context(), uint64(usuarioID), p.SenhaAtual, p.NovaSenha, p.SenhaConfirmacao); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}
	resposta.Padrao(w, http.StatusOK, nil)
}
