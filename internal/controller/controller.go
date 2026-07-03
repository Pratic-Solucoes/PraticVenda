package controller

import (
	"gestao/internal/service"
	"net/http"
)

type Controller struct {
	View interface {
		RenderizarLoginPage(w http.ResponseWriter, r *http.Request)
		RenderizarDashboardPage(w http.ResponseWriter, r *http.Request)
		RenderizarDebitosPage(w http.ResponseWriter, r *http.Request)
		RenderizarClientesPage(w http.ResponseWriter, r *http.Request)
		RenderizarFornecedoresPage(w http.ResponseWriter, r *http.Request)
		RenderizarCategoriasPage(w http.ResponseWriter, r *http.Request)
		RenderizarConfiguracaoUsuarioPage(w http.ResponseWriter, r *http.Request)
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
	}
	Fornecedores interface {
		CriarFornecedor(w http.ResponseWriter, r *http.Request)
		ListarFornecedores(w http.ResponseWriter, r *http.Request)
		ObterFornecedor(w http.ResponseWriter, r *http.Request)
		AtualizarFornecedor(w http.ResponseWriter, r *http.Request)
	}
	Debitos interface {
		CriarDebitoAvulso(w http.ResponseWriter, r *http.Request)
		ListarDebitos(w http.ResponseWriter, r *http.Request)
		PagarDebito(w http.ResponseWriter, r *http.Request)
		EditarDebito(w http.ResponseWriter, r *http.Request)
	}
	Categorias interface {
		CriarCategoria(w http.ResponseWriter, r *http.Request)
		ListarCategorias(w http.ResponseWriter, r *http.Request)
	}
	Dashboard interface {
		ResumoDashboard(w http.ResponseWriter, r *http.Request)
	}
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		View:         &ViewController{},
		Login:        &LoginController{service: service},
		Usuarios:     &UsuarioController{service: service},
		Clientes:     &ClienteController{service: service},
		Fornecedores: &FornecedorController{service: service},
		Categorias:   &CategoriaController{service: service},
		Debitos:      &DebitoController{service: service},
		Dashboard:    &DashboardController{service: service.Dashboard},
	}
}
