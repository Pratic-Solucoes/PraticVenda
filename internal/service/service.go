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
		CriarEndereco(ctx context.Context, idCliente int64, e *model.EnderecoCliente) (*model.EnderecoCliente, error)
		EditarEndereco(ctx context.Context, idCliente int64, idEndereco int64, e *model.EnderecoCliente) error
		BuscarEnderecoByID(ctx context.Context, idCliente int64, idEndereco int64) (*model.EnderecoCliente, error)
		CriarTelefone(ctx context.Context, idCliente int64, t *model.TelefoneCliente) (*model.TelefoneCliente, error)
		EditarTelefone(ctx context.Context, idCliente int64, idTelefone int64, t *model.TelefoneCliente) error
		BuscarTelefoneByID(ctx context.Context, idCliente int64, idTelefone int64) (*model.TelefoneCliente, error)
	}
	Fornecedores interface {
		CriarFornecedor(ctx context.Context, f *model.Fornecedor) (*model.Fornecedor, error)
		ListarFornecedores(ctx context.Context, busca string) ([]*model.Fornecedor, error)
		ObterFornecedorPorID(ctx context.Context, id int64) (*model.Fornecedor, error)
		AtualizarFornecedor(ctx context.Context, id int64, f *model.Fornecedor) error
	}
	ContasPagar interface {
		CriarContaPagar(ctx context.Context, contaPagar *model.ContaPagarCriar) error
		ListarContasPagar(ctx context.Context, busca, vencimento, status string) ([]*model.ContaPagar, error)
		PagarContaPagar(ctx context.Context, id int64, valorPagamento float64) error
		EditarContaPagar(ctx context.Context, id int64, contaPagar *model.ContaPagarCriar) error
	}
	CategoriasContasPagar interface {
		CriarCategoria(ctx context.Context, c *model.CategoriaContaPagar) (*model.CategoriaContaPagar, error)
		ListarCategorias(ctx context.Context) ([]*model.CategoriaContaPagar, error)
	}
	FormasPagamento interface {
		Criar(ctx context.Context, fp *model.FormaPagamento) (*model.FormaPagamento, error)
		Listar(ctx context.Context) ([]model.FormaPagamento, error)
		BuscarPorID(ctx context.Context, idFp int64) (*model.FormaPagamento, error)
		Atualizar(ctx context.Context, fp *model.FormaPagamento) (*model.FormaPagamento, error)
	}
	Estoques interface {
		CriarEstoque(ctx context.Context, input *model.EstoqueCriar) (*model.Estoque, error)
		ListarEstoques(ctx context.Context) ([]*model.Estoque, error)
		ListarProdutosDoEstoque(ctx context.Context, idEstoque int64) ([]*model.ProdutoEstoque, error)
	}
	Produtos interface {
		CriarProduto(ctx context.Context, input *model.ProdutoInput) (*model.ProdutoCompleto, error)
		ListarProdutos(ctx context.Context, busca string) ([]*model.ProdutoCompleto, error)
		ObterProdutoPorID(ctx context.Context, id int64) (*model.ProdutoCompleto, error)
		AtualizarProduto(ctx context.Context, id int64, input *model.ProdutoInput) error
		ExcluirOuInativarProduto(ctx context.Context, id int64) (string, error)
		VincularProdutoEstoque(ctx context.Context, idProduto, idEstoque int64, qtdMinima float64) error
		DesvincularProdutoEstoque(ctx context.Context, idProduto, idEstoque int64) error
		ListarGruposTributarios(ctx context.Context) ([]*model.GrupoTributario, error)
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
		CategoriasContasPagar: &CategoriaContaPagarService{
			repository: repository,
			db:         db,
		},
		ContasPagar: &ContaPagarService{
			repository: repository,
			db:         db,
		},
		FormasPagamento: &FormaPagamentoService{
			repository: repository,
			db:         db,
		},
		Estoques: &EstoqueService{
			repository: repository,
			db:         db,
		},
		Produtos: &ProdutoService{
			repository: repository,
			db:         db,
		},
		Dashboard: NewDashboardService(repository.Dashboard, db),
	}
}
