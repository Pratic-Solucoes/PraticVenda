package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type Repository struct {
	Login interface {
		Login(ctx context.Context, email string) (uint64, string, string, string, error)
		LoginAdministrativo(ctx context.Context, email string) (uint64, string, string, bool, error)
	}
	Admin interface {
		CarregarOrganizacoes(ctx context.Context) ([]model.Organizacao, error)
		CarregarUsuarios(ctx context.Context) ([]model.Usuario, error)
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
		CriarEndereco(ctx context.Context, tx *sql.Tx, idCliente int64, e *model.EnderecoCliente) (*model.EnderecoCliente, error)
		EditarEndereco(ctx context.Context, tx *sql.Tx, idCliente int64, idEndereco int64, e *model.EnderecoCliente) error
		BuscarEnderecoByID(ctx context.Context, tx *sql.Tx, idCliente int64, idEndereco int64) (*model.EnderecoCliente, error)
		CriarTelefone(ctx context.Context, tx *sql.Tx, idCliente int64, t *model.TelefoneCliente) (*model.TelefoneCliente, error)
		EditarTelefone(ctx context.Context, tx *sql.Tx, idCliente int64, idTelefone int64, t *model.TelefoneCliente) error
		BuscarTelefoneByID(ctx context.Context, tx *sql.Tx, idCliente int64, idTelefone int64) (*model.TelefoneCliente, error)
	}
	Fornecedores interface {
		CriarFornecedor(ctx context.Context, tx *sql.Tx, f *model.Fornecedor) (*model.Fornecedor, error)
		ListarFornecedores(ctx context.Context, tx *sql.Tx, busca string) ([]*model.Fornecedor, error)
		ObterFornecedorPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.Fornecedor, error)
		AtualizarFornecedor(ctx context.Context, tx *sql.Tx, id int64, f *model.Fornecedor) error
	}
	ContasPagar interface {
		CriarContaPagar(ctx context.Context, tx *sql.Tx, contaPagar *model.ContaPagarCriar) error
		ListarContasPagar(ctx context.Context, tx *sql.Tx, busca, vencimento, status string) ([]*model.ContaPagar, error)
		PagarContaPagar(ctx context.Context, tx *sql.Tx, id int64, valorPagamento float64) error
		EditarContaPagar(ctx context.Context, tx *sql.Tx, id int64, contaPagar *model.ContaPagarCriar) error
		BuscarPorID(ctx context.Context, tx *sql.Tx, ID int64) (*model.ContaPagar, error)
	}
	FormasPagamento interface {
		Criar(ctx context.Context, tx *sql.Tx, fp *model.FormaPagamento) (*model.FormaPagamento, error)
		Listar(ctx context.Context, tx *sql.Tx) ([]model.FormaPagamento, error)
		BuscarPorID(ctx context.Context, tx *sql.Tx, idFp int64) (*model.FormaPagamento, error)
		Atualizar(ctx context.Context, tx *sql.Tx, fp *model.FormaPagamento) (*model.FormaPagamento, error)
	}
	CategoriasContasPagar interface {
		CriarCategoria(ctx context.Context, tx *sql.Tx, c *model.CategoriaContaPagar) (*model.CategoriaContaPagar, error)
		ListarCategorias(ctx context.Context, tx *sql.Tx) ([]model.CategoriaContaPagar, error)
	}
	Estoques interface {
		CriarEstoque(ctx context.Context, tx *sql.Tx, e *model.Estoque) (*model.Estoque, error)
		ListarEstoques(ctx context.Context, tx *sql.Tx) ([]*model.Estoque, error)
		ListarProdutosDoEstoque(ctx context.Context, tx *sql.Tx, idEstoque int64) ([]*model.ProdutoEstoque, error)
	}
	EntradaEstoque interface {
		RegistrarEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error
		RegistrarProdutosEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error
	}
	Produtos interface {
		CriarProduto(ctx context.Context, tx *sql.Tx, p *model.Produto, f *model.ProdutoFiscal) (*model.Produto, error)
		AtualizarProduto(ctx context.Context, tx *sql.Tx, id int64, p *model.Produto, f *model.ProdutoFiscal) error
		ListarProdutos(ctx context.Context, tx *sql.Tx, busca string) ([]*model.ProdutoCompleto, error)
		ObterProdutoPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.ProdutoCompleto, error)
		ExcluirProduto(ctx context.Context, tx *sql.Tx, id int64) error
		InativarProduto(ctx context.Context, tx *sql.Tx, id int64) error
		TemMovimentacaoEstoque(ctx context.Context, tx *sql.Tx, idProduto int64) (bool, error)
		TemMovimentacaoNoEstoqueEspecifico(ctx context.Context, tx *sql.Tx, idProduto, idEstoque int64) (bool, error)
		VincularAoEstoque(ctx context.Context, tx *sql.Tx, idProduto, idEstoque int64, quantidadeMinima, quantidade float64) error
		DesvincularDoEstoque(ctx context.Context, tx *sql.Tx, idProduto, idEstoque int64) error
		BuscarEstoqueVinculos(ctx context.Context, tx *sql.Tx, idProduto int64) ([]model.ProdutoEstoqueInfo, error)
		ListarGruposTributarios(ctx context.Context, tx *sql.Tx) ([]*model.GrupoTributario, error)
	}
	Dashboard *DashboardRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Login: &LoginRepository{
			db: db,
		},
		Admin: &AdminRepository{
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
		CategoriasContasPagar: &CategoriaContaPagarRepository{
			db: db,
		},
		ContasPagar: &ContaPagarRepository{
			db: db,
		},
		FormasPagamento: &FormaPagamentoRepository{
			db: db,
		},
		Estoques: &EstoqueRepository{
			db: db,
		},
		EntradaEstoque: &EntradaEstoqueRepository{
			db: db,
		},
		Produtos: &ProdutoRepository{
			db: db,
		},
		Dashboard: NewDashboardRepository(db),
	}
}
