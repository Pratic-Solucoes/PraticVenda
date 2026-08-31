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
	staticHandler := http.StripPrefix("/static/", http.FileServer(filesDir))
	r.Handle("/static/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		staticHandler.ServeHTTP(w, req)
	}))

	// rotas renderizar páginas
	r.Get("/", c.View.RenderizarLoginPage)
	r.Get("/dashboard", c.View.RenderizarDashboardPage)
	r.Get("/configuracao-usuario", c.View.RenderizarConfiguracaoUsuarioPage)
	r.Get("/configuracao-caixa", c.View.RenderizarConfiguracaoCaixaPage)
	r.Get("/configuracao-pdv", c.View.RenderizarConfiguracaoPDVPage)
	r.Get("/contas-pagar", c.View.RenderizarContasPagarPage)
	r.Get("/contas-receber", c.View.RenderizarContasReceberPage)
	r.Get("/clientes", c.View.RenderizarClientesPage)
	r.Get("/fornecedores", c.View.RenderizarFornecedoresPage)
	r.Get("/categorias-contas-pagar", c.View.RenderizarCategoriasPage)
	r.Get("/pdv", c.View.RenderizarPdvPage)
	r.Get("/formas-pagamento", c.View.RenderizarFormasPagamentoPage)
	r.Get("/estoques", c.View.RenderizarEstoquesPage)
	r.Get("/entrada-estoque", c.View.RenderizarEntradaEstoquePage)
	r.Get("/saida-estoque", c.View.RenderizarSaidaEstoquePage)
	r.Get("/produtos", c.View.RenderizarProdutosPage)
	r.Get("/condicoes-pagamento", c.View.RenderizarCondicoesPagamentoPage)

	// rotas funcionalidades
	r.Route("/api", func(r chi.Router) {
		r.Post("/login", c.Login.Login)
		r.Post("/pdv/vendas", auth.Autenticar(c.VendaPDV.Finalizar))
		r.Post("/pdv/pre-vendas", auth.Autenticar(c.VendaPDV.SalvarPreVenda))
		r.Get("/pdv/pre-vendas", auth.Autenticar(c.VendaPDV.ListarPreVendas))
		r.Get("/pdv/pre-vendas/{id}", auth.Autenticar(c.VendaPDV.ObterPreVenda))
		r.Post("/pdv/vendas/{id}/cancelar", auth.Autenticar(c.VendaPDV.Cancelar))
		r.Get("/configuracoes/pdv", auth.Autenticar(c.ConfiguracaoPDV.Obter))
		r.Put("/configuracoes/pdv", auth.Autenticar(c.ConfiguracaoPDV.Salvar))
		r.Get("/caixas", auth.Autenticar(c.Caixa.Listar))
		r.Get("/caixas/usuarios", auth.Autenticar(c.Caixa.ListarUsuarios))
		r.Post("/caixas", auth.Autenticar(c.Caixa.Criar))
		r.Get("/caixas/atual", auth.Autenticar(c.Caixa.Atual))
		r.Post("/caixas/abrir", auth.Autenticar(c.Caixa.Abrir))
		r.Post("/caixas/fechar", auth.Autenticar(c.Caixa.Fechar))

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
		r.Get("/formas-pagamento/{id}", auth.Autenticar(c.FormasPagamento.BuscarPorID))
		r.Put("/formas-pagamento/{id}", auth.Autenticar(c.FormasPagamento.Atualizar))

		// Rotas Contas a Pagar
		r.Get("/contas-pagar/categorias", auth.Autenticar(c.CategoriasContasPagar.ListarCategorias))
		r.Post("/contas-pagar/categorias", auth.Autenticar(c.CategoriasContasPagar.CriarCategoria))

		r.Get("/contas-pagar", auth.Autenticar(c.ContasPagar.ListarContasPagar))
		r.Post("/contas-pagar", auth.Autenticar(c.ContasPagar.CriarContaPagar))
		r.Put("/contas-pagar/{id}", auth.Autenticar(c.ContasPagar.EditarContaPagar))
		r.Put("/contas-pagar/{id}/pagar", auth.Autenticar(c.ContasPagar.PagarContaPagar))

		r.Get("/contas-receber/categorias", auth.Autenticar(c.CategoriasContasReceber.ListarCategorias))
		r.Post("/contas-receber/categorias", auth.Autenticar(c.CategoriasContasReceber.CriarCategoria))
		r.Get("/contas-receber", auth.Autenticar(c.ContasReceber.ListarContasReceber))
		r.Post("/contas-receber", auth.Autenticar(c.ContasReceber.CriarContaReceber))
		r.Put("/contas-receber/{id}/receber", auth.Autenticar(c.ContasReceber.ReceberConta))

		// Rotas Condições de Pagamento
		r.Post("/condicoes-pagamento", auth.Autenticar(c.CondicoesPagamento.Criar))
		r.Get("/condicoes-pagamento", auth.Autenticar(c.CondicoesPagamento.Listar))
		r.Get("/condicoes-pagamento/{id}", auth.Autenticar(c.CondicoesPagamento.BuscarPorID))
		r.Put("/condicoes-pagamento/{id}", auth.Autenticar(c.CondicoesPagamento.Atualizar))

		// Rotas Estoques
		r.Get("/estoques", auth.Autenticar(c.Estoques.ListarEstoques))
		r.Post("/estoques", auth.Autenticar(c.Estoques.CriarEstoque))
		r.Get("/estoques/{id}/produtos", auth.Autenticar(c.Estoques.ListarProdutosDoEstoque))

		// Rotas de entradas de estoque
		r.Post("/entradas-estoque", auth.Autenticar(c.EntradaEstoque.RegistrarEntrada))
		r.Get("/entradas-estoque", auth.Autenticar(c.EntradaEstoque.ListarEntradas))
		r.Get("/entradas-estoque/{id}", auth.Autenticar(c.EntradaEstoque.ObterEntrada))
		r.Put("/entradas-estoque/{id}", auth.Autenticar(c.EntradaEstoque.EditarEntrada))
		r.Post("/entradas-estoque/{id}/aprovar", auth.Autenticar(c.EntradaEstoque.AprovarEntrada))
		r.Post("/entradas-estoque/{id}/cancelar", auth.Autenticar(c.EntradaEstoque.CancelarEntrada))
		r.Post("/saidas-estoque", auth.Autenticar(c.SaidaEstoque.RegistrarSaida))
		r.Get("/saidas-estoque", auth.Autenticar(c.SaidaEstoque.ListarSaidas))
		r.Get("/saidas-estoque/{id}", auth.Autenticar(c.SaidaEstoque.ObterSaida))
		r.Put("/saidas-estoque/{id}", auth.Autenticar(c.SaidaEstoque.EditarSaida))
		r.Post("/saidas-estoque/{id}/aprovar", auth.Autenticar(c.SaidaEstoque.AprovarSaida))
		r.Post("/saidas-estoque/{id}/cancelar", auth.Autenticar(c.SaidaEstoque.CancelarSaida))

		// Rotas Produtos
		r.Get("/produtos", auth.Autenticar(c.Produtos.ListarProdutos))
		r.Post("/produtos", auth.Autenticar(c.Produtos.CriarProduto))
		r.Get("/produtos/{id}", auth.Autenticar(c.Produtos.ObterProduto))
		r.Put("/produtos/{id}", auth.Autenticar(c.Produtos.AtualizarProduto))
		r.Get("/produtos/{id}/composicao", auth.Autenticar(c.Produtos.ListarComposicao))
		r.Put("/produtos/{id}/composicao", auth.Autenticar(c.Produtos.SalvarComposicao))
		r.Delete("/produtos/{id}", auth.Autenticar(c.Produtos.ExcluirProduto))
		r.Post("/produtos/{id}/estoques", auth.Autenticar(c.Produtos.VincularEstoque))
		r.Delete("/produtos/{id}/estoques/{id_estoque}", auth.Autenticar(c.Produtos.DesvincularEstoque))
		r.Get("/grupos-tributarios", auth.Autenticar(c.Produtos.ListarGruposTributarios))

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
