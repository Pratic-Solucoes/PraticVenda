package repository

import (
	"context"
	"database/sql"
	"gestao/internal/model"
)

type Repository struct {
	Login interface {
		Login(ctx context.Context, username string) (uint64, string, string, error)
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
	ContasReceber interface {
		CriarParcelas(context.Context, *sql.Tx, *model.ContaReceberCriar, string, []float64, []string) error
		Listar(context.Context, *sql.Tx, string, string, string) ([]*model.ContaReceber, error)
		BuscarPorID(context.Context, *sql.Tx, int64) (*model.ContaReceber, error)
		Receber(context.Context, *sql.Tx, int64, float64, int64) error
	}
	CondicoesPagamento interface {
		Criar(ctx context.Context, tx *sql.Tx, cp *model.CondicaoPagamento) error
		Listar(ctx context.Context, tx *sql.Tx) ([]model.CondicaoPagamento, error)
		BuscarPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.CondicaoPagamento, error)
		Atualizar(ctx context.Context, tx *sql.Tx, cp *model.CondicaoPagamento) error
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
	CategoriasContasReceber interface {
		CriarCategoria(context.Context, *sql.Tx, *model.CategoriaContaReceber) (*model.CategoriaContaReceber, error)
		ListarCategorias(context.Context, *sql.Tx) ([]model.CategoriaContaReceber, error)
	}
	Estoques interface {
		CriarEstoque(ctx context.Context, tx *sql.Tx, e *model.Estoque) (*model.Estoque, error)
		ListarEstoques(ctx context.Context, tx *sql.Tx) ([]*model.Estoque, error)
		ListarProdutosDoEstoque(ctx context.Context, tx *sql.Tx, idEstoque int64) ([]*model.ProdutoEstoque, error)
	}
	EntradaEstoque interface {
		RegistrarEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error
		SalvarItensEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error
		RegistrarProdutosEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error
		ListarEntradasEstoque(ctx context.Context, tx *sql.Tx, filtro model.FiltroEntradaEstoque) ([]model.EntradaEstoque, error)
		ObterEntrada(ctx context.Context, tx *sql.Tx, id uint64) (*model.EntradaEstoque, error)
		AtualizarEntrada(ctx context.Context, tx *sql.Tx, entrada *model.EntradaEstoque) error
		AprovarEntrada(ctx context.Context, tx *sql.Tx, id uint64, usuarioID int64) error
		CancelarEntrada(ctx context.Context, tx *sql.Tx, id uint64, usuarioID int64) error
		ValidarProdutosFornecedor(ctx context.Context, tx *sql.Tx, idFornecedor uint64, produtos []model.ProdutoEntradaEstoque) error
	}
	SaidaEstoque interface {
		RegistrarSaida(ctx context.Context, tx *sql.Tx, saida *model.SaidaEstoque) error
		SalvarItens(ctx context.Context, tx *sql.Tx, saida *model.SaidaEstoque) error
		Listar(ctx context.Context, tx *sql.Tx, filtro model.FiltroSaidaEstoque) ([]model.SaidaEstoque, error)
		Obter(ctx context.Context, tx *sql.Tx, id uint64) (*model.SaidaEstoque, error)
		Atualizar(ctx context.Context, tx *sql.Tx, saida *model.SaidaEstoque) error
		Aprovar(ctx context.Context, tx *sql.Tx, id uint64, usuarioID int64) error
		CancelarComUsuario(ctx context.Context, tx *sql.Tx, id uint64, usuario int64) error
	}
	Produtos interface {
		CriarProduto(ctx context.Context, tx *sql.Tx, p *model.Produto, f *model.ProdutoFiscal) (*model.Produto, error)
		AtualizarProduto(ctx context.Context, tx *sql.Tx, id int64, p *model.Produto, f *model.ProdutoFiscal) error
		ListarProdutos(ctx context.Context, tx *sql.Tx, busca string, idFornecedor int64) ([]*model.ProdutoCompleto, error)
		ObterProdutoPorID(ctx context.Context, tx *sql.Tx, id int64) (*model.ProdutoCompleto, error)
		ExcluirProduto(ctx context.Context, tx *sql.Tx, id int64) error
		InativarProduto(ctx context.Context, tx *sql.Tx, id int64) error
		TemMovimentacaoEstoque(ctx context.Context, tx *sql.Tx, idProduto int64) (bool, error)
		TemMovimentacaoNoEstoqueEspecifico(ctx context.Context, tx *sql.Tx, idProduto, idEstoque int64) (bool, error)
		VincularAoEstoque(ctx context.Context, tx *sql.Tx, idProduto, idEstoque int64, quantidadeMinima, quantidade float64) error
		DesvincularDoEstoque(ctx context.Context, tx *sql.Tx, idProduto, idEstoque int64) error
		BuscarEstoqueVinculos(ctx context.Context, tx *sql.Tx, idProduto int64) ([]model.ProdutoEstoqueInfo, error)
		ListarGruposTributarios(ctx context.Context, tx *sql.Tx) ([]*model.GrupoTributario, error)
		ListarComposicao(ctx context.Context, tx *sql.Tx, idProduto int64) ([]model.ItemComposicaoProduto, error)
		SalvarComposicao(ctx context.Context, tx *sql.Tx, idProduto int64, itens []model.ItemComposicaoProduto) error
		SincronizarFornecedores(ctx context.Context, tx *sql.Tx, idProduto int64, ids []int64) error
		ListarFornecedoresProduto(ctx context.Context, tx *sql.Tx, idProduto int64) ([]int64, error)
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
		CategoriasContasPagar: &CategoriaContaPagarRepository{
			db: db,
		},
		CategoriasContasReceber: &CategoriaContaReceberRepository{db: db},
		ContasPagar: &ContaPagarRepository{
			db: db,
		},
		ContasReceber: &ContaReceberRepository{db: db},
		CondicoesPagamento: &CondicaoPagamentoRepository{
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
		SaidaEstoque: &SaidaEstoqueRepository{db: db},
		Produtos: &ProdutoRepository{
			db: db,
		},
		Dashboard: NewDashboardRepository(db),
	}
}
