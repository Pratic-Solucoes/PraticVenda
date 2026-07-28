package usuario

import (
	"gestao/internal/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CarregarRotas(c ControllerInterface) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/usuarios", func(r chi.Router) {
		r.Post("/", http.HandlerFunc(c.CriarUsuario))
		r.Get("/", auth.Autenticar(c.BuscarUsuarioPorID))
		r.Put("/", auth.Autenticar(c.EditarUsuario))
		r.Put("/alterar-senha", auth.Autenticar(c.AlterarSenha))
	})
	return r
}
