package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type Service struct {
	Usuarios interface {
		CriarUsuario(ctx context.Context, usuario *model.UsuarioCriar) (*model.UsuarioBasico, error)
		BuscarUsuarioPorID(ctx context.Context, usuarioID int) (*model.Usuario, error)
		AlterarSenha(ctx context.Context, usuarioID int64, senhaAtual, senhaNova, senhaNovaConfirmacao string) error
	}
	Login interface {
		Login(ctx context.Context, usuario *model.UsuarioLogin) (uint64, string, string, error)
	}
	Clientes interface {
		CriarCliente(ctx context.Context, c *model.Cliente) (*model.Cliente, error)
		ListarClientes(ctx context.Context, busca string) ([]model.Cliente, error)
		ObterClientePorID(ctx context.Context, id int64) (*model.Cliente, error)
		AtualizarCliente(ctx context.Context, id int64, c *model.Cliente) error
	}
	Fornecedores interface {
		CriarFornecedor(ctx context.Context, f *model.Fornecedor) (*model.Fornecedor, error)
		ListarFornecedores(ctx context.Context, busca string) ([]*model.Fornecedor, error)
		ObterFornecedorPorID(ctx context.Context, id int64) (*model.Fornecedor, error)
		AtualizarFornecedor(ctx context.Context, id int64, f *model.Fornecedor) error
	}
	Debitos interface {
		LancarDebito(ctx context.Context, debito *model.DebitoAvulsoCriar) error
		ListarDebitos(ctx context.Context, busca, vencimento, status string) ([]*model.Debito, error)
		PagarDebito(ctx context.Context, id int64) error
		EditarDebito(ctx context.Context, id int64, debito *model.DebitoAvulsoCriar) error
	}
	Categorias interface {
		CriarCategoria(ctx context.Context, c *model.CategoriaDebito) (*model.CategoriaDebito, error)
		ListarCategorias(ctx context.Context) ([]*model.CategoriaDebito, error)
	}
	Dashboard *DashboardService
}

func NewService(repository *repository.Repository, db *sql.DB) *Service {
	return &Service{
		Usuarios: &UsuarioService{
			repository: repository,
			db:         db,
		},
		Login: &LoginService{
			repository: repository,
			db:         db,
		},
		Clientes: &ClienteService{
			repository: repository,
			db:         db,
		},
		Fornecedores: &FornecedorService{
			repository: repository,
			db:         db,
		},
		Categorias: &CategoriaService{
			repository: repository,
			db:         db,
		},
		Debitos: &DebitoService{
			repository: repository,
			db:         db,
		},
		Dashboard: NewDashboardService(repository.Dashboard, db),
	}
}
