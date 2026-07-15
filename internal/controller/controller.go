package controller

import (
	"gestao/internal/service"
	"net/http"
)

type Controller struct {
	View interface {
		RenderizarLoginPage(w http.ResponseWriter, r *http.Request)
		RenderizarDashboardPage(w http.ResponseWriter, r *http.Request)
		RenderizarContasPagarPage(w http.ResponseWriter, r *http.Request)
		RenderizarClientesPage(w http.ResponseWriter, r *http.Request)
		RenderizarFornecedoresPage(w http.ResponseWriter, r *http.Request)
		RenderizarCategoriasPage(w http.ResponseWriter, r *http.Request)
		RenderizarConfiguracaoUsuarioPage(w http.ResponseWriter, r *http.Request)
		RenderizarPdvPage(w http.ResponseWriter, r *http.Request)
		RenderizarFormasPagamentoPage(w http.ResponseWriter, r *http.Request)
		RenderizarEstoquesPage(w http.ResponseWriter, r *http.Request)
		RenderizarEntradaEstoquePage(w http.ResponseWriter, r *http.Request)
		RenderizarProdutosPage(w http.ResponseWriter, r *http.Request)
	}
	Login interface {
		Login(w http.ResponseWriter, r *http.Request)
	}
	Empresa interface {
		CriarEmpresa(w http.ResponseWriter, r *http.Request)
	}
	Usuarios interface {
		CriarUsuario(w http.ResponseWriter, r *http.Request)
		BuscarUsuarioPorID(w http.ResponseWriter, r *http.Request)
		EditarUsuario(w http.ResponseWriter, r *http.Request)
		AlterarSenha(w http.ResponseWriter, r *http.Request)
	}
	Clientes interface {
		CriarCliente(w http.ResponseWriter, r *http.Request)
		ListarClientes(w http.ResponseWriter, r *http.Request)
		ObterCliente(w http.ResponseWriter, r *http.Request)
		AtualizarCliente(w http.ResponseWriter, r *http.Request)
		CriarEndereco(w http.ResponseWriter, r *http.Request)
		BuscarEnderecoByID(w http.ResponseWriter, r *http.Request)
		AtualizarEndereco(w http.ResponseWriter, r *http.Request)
		CriarTelefone(w http.ResponseWriter, r *http.Request)
		BuscarTelefoneByID(w http.ResponseWriter, r *http.Request)
		AtualizarTelefone(w http.ResponseWriter, r *http.Request)
	}
	Fornecedores interface {
		CriarFornecedor(w http.ResponseWriter, r *http.Request)
		ListarFornecedores(w http.ResponseWriter, r *http.Request)
		ObterFornecedor(w http.ResponseWriter, r *http.Request)
		AtualizarFornecedor(w http.ResponseWriter, r *http.Request)
	}
	ContasPagar interface {
		CriarContaPagar(w http.ResponseWriter, r *http.Request)
		ListarContasPagar(w http.ResponseWriter, r *http.Request)
		PagarContaPagar(w http.ResponseWriter, r *http.Request)
		EditarContaPagar(w http.ResponseWriter, r *http.Request)
	}
	CategoriasContasPagar interface {
		CriarCategoria(w http.ResponseWriter, r *http.Request)
		ListarCategorias(w http.ResponseWriter, r *http.Request)
	}
	FormasPagamento interface {
		Criar(w http.ResponseWriter, r *http.Request)
		Listar(w http.ResponseWriter, r *http.Request)
		BuscarPorID(w http.ResponseWriter, r *http.Request)
		Atualizar(w http.ResponseWriter, r *http.Request)
	}
	Estoques interface {
		CriarEstoque(w http.ResponseWriter, r *http.Request)
		ListarEstoques(w http.ResponseWriter, r *http.Request)
		ListarProdutosDoEstoque(w http.ResponseWriter, r *http.Request)
		EntradaEstoque(w http.ResponseWriter, r *http.Request)
	}
	EntradaEstoque interface {
		RegistrarEntrada(w http.ResponseWriter, r *http.Request)
	}
	Produtos interface {
		CriarProduto(w http.ResponseWriter, r *http.Request)
		ListarProdutos(w http.ResponseWriter, r *http.Request)
		ObterProduto(w http.ResponseWriter, r *http.Request)
		AtualizarProduto(w http.ResponseWriter, r *http.Request)
		ExcluirProduto(w http.ResponseWriter, r *http.Request)
		VincularEstoque(w http.ResponseWriter, r *http.Request)
		DesvincularEstoque(w http.ResponseWriter, r *http.Request)
		ListarGruposTributarios(w http.ResponseWriter, r *http.Request)
	}
	Dashboard interface {
		ResumoDashboard(w http.ResponseWriter, r *http.Request)
	}
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		View:                  &ViewController{},
		Login:                 &LoginController{service: service},
		Usuarios:              &UsuarioController{service: service},
		Clientes:              &ClienteController{service: service},
		Fornecedores:          &FornecedorController{service: service},
		CategoriasContasPagar: &CategoriaContaPagarController{service: service},
		ContasPagar:           &ContaPagarController{service: service},
		FormasPagamento:       &FormaPagamentoController{service: service},
		Estoques:              &EstoqueController{service: service},
		Produtos:              &ProdutoController{service: service},
		Dashboard:             &DashboardController{service: service.Dashboard},
	}
}
