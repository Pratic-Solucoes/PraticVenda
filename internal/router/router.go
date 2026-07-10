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
	r.Get("/contas-pagar", c.View.RenderizarContasPagarPage)
	r.Get("/clientes", c.View.RenderizarClientesPage)
	r.Get("/fornecedores", c.View.RenderizarFornecedoresPage)
	r.Get("/categorias-contas-pagar", c.View.RenderizarCategoriasPage)
	r.Get("/pdv", c.View.RenderizarPdvPage)
	r.Get("/formas-pagamento", c.View.RenderizarFormasPagamentoPage)

	// rotas funcionalidades
	r.Route("/api", func(r chi.Router) {
		r.Post("/login", c.Login.Login)

		// Rotas Usuarios
		r.Post("/usuarios", c.Usuarios.CriarUsuario)
		r.Get("/usuario", auth.Autenticar(c.Usuarios.BuscarUsuarioPorID))
		r.Put("/usuario/alterar-senha", auth.Autenticar(c.Usuarios.AlterarSenha))

		// Rotas Fornecedores
		r.Get("/fornecedores", auth.Autenticar(c.Fornecedores.ListarFornecedores))
		r.Post("/fornecedores", auth.Autenticar(c.Fornecedores.CriarFornecedor))
		r.Get("/fornecedores/{id}", auth.Autenticar(c.Fornecedores.ObterFornecedor))
		r.Put("/fornecedores/{id}", auth.Autenticar(c.Fornecedores.AtualizarFornecedor))

		// Rotas Formas de Pagamento
		r.Post("/formas-pagamento", auth.Autenticar(c.FormasPagamento.Criar))
		r.Get("/formas-pagamento", auth.Autenticar(c.FormasPagamento.Listar))

		// Rotas Contas a Pagar
		r.Get("/contas-pagar/categorias", auth.Autenticar(c.CategoriasContasPagar.ListarCategorias))
		r.Post("/contas-pagar/categorias", auth.Autenticar(c.CategoriasContasPagar.CriarCategoria))

		r.Get("/contas-pagar", auth.Autenticar(c.ContasPagar.ListarContasPagar))
		r.Post("/contas-pagar", auth.Autenticar(c.ContasPagar.CriarContaPagar))
		r.Put("/contas-pagar/{id}", auth.Autenticar(c.ContasPagar.EditarContaPagar))
		r.Put("/contas-pagar/{id}/pagar", auth.Autenticar(c.ContasPagar.PagarContaPagar))
		r.Post("/contas-pagar/{id}/pagar-parcial", nil)

		// Rotas Clientes
		r.Get("/clientes", auth.Autenticar(c.Clientes.ListarClientes))
		r.Post("/clientes", auth.Autenticar(c.Clientes.CriarCliente))
		r.Get("/clientes/{id}", auth.Autenticar(c.Clientes.ObterCliente))
		r.Put("/clientes/{id}", auth.Autenticar(c.Clientes.AtualizarCliente))

		r.Post("/clientes/{id}/enderecos", auth.Autenticar(c.Clientes.CriarEndereco))
		r.Get("/clientes/{id}/enderecos/{id_endereco}", auth.Autenticar(c.Clientes.BuscarEnderecoByID))
		r.Put("/clientes/{id}/enderecos/{id_endereco}", auth.Autenticar(c.Clientes.AtualizarEndereco))

		r.Post("/clientes/{id}/telefones", auth.Autenticar(c.Clientes.CriarTelefone))
		r.Get("/clientes/{id}/telefones/{id_telefone}", auth.Autenticar(c.Clientes.BuscarTelefoneByID))
		r.Put("/clientes/{id}/telefones/{id_telefone}", auth.Autenticar(c.Clientes.AtualizarTelefone))

		r.Get("/dashboard/resumo", auth.Autenticar(c.Dashboard.ResumoDashboard))
	})

	return r
}
