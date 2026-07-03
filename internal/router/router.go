package router

import (
	"gestao/config"
	"gestao/internal/auth"
	"gestao/internal/controller"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func CarregarRotas(c *controller.Controller) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)
	r.Use(chi_middleware.Throttle(config.GetInt("ROUTER_MAX_REQUESTS_PER_MINUTE", 10)))

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	filesDir := http.Dir(filepath.Join(workDir, "web", "static"))
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(filesDir)))

	// rotas renderizar páginas
	r.Get("/", c.View.RenderizarLoginPage)
	r.Get("/dashboard", c.View.RenderizarDashboardPage)
	r.Get("/configuracao-usuario", c.View.RenderizarConfiguracaoUsuarioPage)
	r.Get("/debitos", c.View.RenderizarDebitosPage)
	r.Get("/clientes", c.View.RenderizarClientesPage)
	r.Get("/fornecedores", c.View.RenderizarFornecedoresPage)
	r.Get("/categorias-debito", c.View.RenderizarCategoriasPage)

	// rotas funcionalidades
	r.Route("/api", func(r chi.Router) {
		r.Post("/login", c.Login.Login)
		r.Post("/usuarios", c.Usuarios.CriarUsuario)
		r.Get("/usuario", auth.Autenticar(c.Usuarios.BuscarUsuarioPorID))
		r.Put("/usuario/alterar-senha", auth.Autenticar(c.Usuarios.AlterarSenha))

		r.Get("/fornecedores", auth.Autenticar(c.Fornecedores.ListarFornecedores))
		r.Post("/fornecedores", auth.Autenticar(c.Fornecedores.CriarFornecedor))
		r.Get("/fornecedores/{id}", auth.Autenticar(c.Fornecedores.ObterFornecedor))
		r.Put("/fornecedores/{id}", auth.Autenticar(c.Fornecedores.AtualizarFornecedor))

		r.Get("/categorias", auth.Autenticar(c.Categorias.ListarCategorias))
		r.Post("/categorias", auth.Autenticar(c.Categorias.CriarCategoria))

		r.Get("/debitos", auth.Autenticar(c.Debitos.ListarDebitos))
		r.Post("/debitos", auth.Autenticar(c.Debitos.CriarDebitoAvulso))
		r.Put("/debitos/{id}", auth.Autenticar(c.Debitos.EditarDebito))
		r.Put("/debitos/{id}/pagar", auth.Autenticar(c.Debitos.PagarDebito))

		r.Get("/clientes", auth.Autenticar(c.Clientes.ListarClientes))
		r.Post("/clientes", auth.Autenticar(c.Clientes.CriarCliente))
		r.Get("/clientes/{id}", auth.Autenticar(c.Clientes.ObterCliente))
		r.Put("/clientes/{id}", auth.Autenticar(c.Clientes.AtualizarCliente))

		r.Get("/dashboard/resumo", auth.Autenticar(c.Dashboard.ResumoDashboard))
	})

	return r
}
