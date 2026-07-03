package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type Repository struct {
	Login interface {
		Login(ctx context.Context, email string) (uint64, string, string, string, error)
	}
	Usuarios interface {
		CriarUsuario(ctx context.Context, tx *sql.Tx, usuario *model.UsuarioCriar) (*model.UsuarioBasico, error)
		BuscarUsuarioPorID(ctx context.Context, usuarioID int) (*model.Usuario, error)
		BuscarSenhaUsuario(ctx context.Context, usuarioID int64) (*string, error)
		AtualizarSenhaUsuario(ctx context.Context, tx *sql.Tx, usuarioID int64, novaSenha string) error
	}
	Clientes interface {
		CriarCliente(ctx context.Context, tx *sql.Tx, c *model.Cliente) (*model.Cliente, error)
		ListarClientes(ctx context.Context, tx *sql.Tx, busca string) ([]model.Cliente, error)
		ObterClientePorID(ctx context.Context, tx *sql.Tx, id int64) (*model.Cliente, error)
		AtualizarCliente(ctx context.Context, tx *sql.Tx, id int64, c *model.Cliente) error
	}
	Fornecedores interface {
		CriarFornecedor(ctx context.Context, tx *sql.Tx, f *model.Fornecedor) (*model.Fornecedor, error)
		ListarFornecedores(ctx context.Context, tx *sql.Tx, busca string) ([]*model.Fornecedor, error)
		ObterFornecedorPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.Fornecedor, error)
		AtualizarFornecedor(ctx context.Context, tx *sql.Tx, id int64, f *model.Fornecedor) error
	}
	Debitos interface {
		LancarDebito(ctx context.Context, tx *sql.Tx, debito *model.DebitoAvulsoCriar) error
		ListarDebitos(ctx context.Context, tx *sql.Tx, busca, vencimento, status string) ([]*model.Debito, error)
		PagarDebito(ctx context.Context, tx *sql.Tx, id int64) error
		EditarDebito(ctx context.Context, tx *sql.Tx, id int64, debito *model.DebitoAvulsoCriar) error
	}
	Categorias interface {
		CriarCategoria(ctx context.Context, tx *sql.Tx, c *model.CategoriaDebito) (*model.CategoriaDebito, error)
		ListarCategorias(ctx context.Context, tx *sql.Tx) ([]*model.CategoriaDebito, error)
	}
	Dashboard *DashboardRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Login: &LoginRepository{
			db: db,
		},
		Usuarios: &UsuarioRepository{
			db: db,
		},
		Fornecedores: &FornecedorRepository{
			db: db,
		},
		Clientes: &ClienteRepository{
			db: db,
		},
		Categorias: NovoCategoriaRepository(db),
		Debitos: &DebitoRepository{
			db: db,
		},
		Dashboard: NewDashboardRepository(db),
	}
}
