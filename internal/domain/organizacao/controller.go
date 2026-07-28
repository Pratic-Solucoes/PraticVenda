package organizacao

import (
	"gestao/pkg/requisicao"
	"gestao/pkg/resposta"
	"net/http"
)

type ControllerInterface interface {
	Criar(w http.ResponseWriter, r *http.Request)
}

type OrganizacaoController struct {
	service ServiceInterface
}

func NewOrganizacaoController() ControllerInterface {
	return &OrganizacaoController{}
}

func (c *OrganizacaoController) Criar(w http.ResponseWriter, r *http.Request) {
	var organizacao Organizacao
	var idUsuario = r.Context().Value("idUsuario").(int64)

	organizacao.IDDono = uint64(idUsuario)

	if err := requisicao.ProcessarRequisicao(w, r, &organizacao); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	if err := c.service.Criar(ctx, organizacao); err != nil {
		resposta.Padrao(w, http.StatusInternalServerError, err)
		return
	}
}

func (c *OrganizacaoController) Editar(w http.ResponseWriter, r *http.Request) {

	organizacaoId := r.PathValue("id")
	if organizacaoId == "" {
		resposta.Padrao(w, http.StatusBadRequest, "ID da organização não encontrado")
		return
	}

	var organizacao Organizacao
	if err := requisicao.ProcessarRequisicao(w, r, &organizacao); err != nil {
		resposta.Padrao(w, http.StatusBadRequest, err)
		return
	}

}
