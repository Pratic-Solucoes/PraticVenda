package service

import (
	"context"
	"database/sql"
	"gestao/internal/model"
	"gestao/internal/repository"
)

type Service struct {
	VendaPDV        *VendaPDVService
	ConfiguracaoPDV *ConfiguracaoPDVService
	Caixa           *CaixaService
	Usuarios        interface {
		CriarUsuario(ctx context.Context, usuario *model.UsuarioCriar) (*model.UsuarioBasico, error)
		BuscarUsuarioPorID(ctx context.Context, usuarioID int) (*model.Usuario, error)
		AlterarSenha(ctx context.Context, usuarioID int64, senhaAtual, senhaNova, senhaNovaConfirmacao string) error
	}
	Login interface {
		Login(ctx context.Context, usuario *model.UsuarioLogin) (uint64, string, error)
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
	ContasReceber interface {
		CriarContaReceber(context.Context, *model.ContaReceberCriar) error
		ListarContasReceber(context.Context, string, string, string) ([]*model.ContaReceber, error)
		ReceberConta(context.Context, int64, float64, int64) error
	}
	CondicoesPagamento interface {
		Criar(ctx context.Context, cp *model.CondicaoPagamento) error
		Listar(ctx context.Context) ([]model.CondicaoPagamento, error)
		BuscarPorID(ctx context.Context, id int64) (*model.CondicaoPagamento, error)
		Atualizar(ctx context.Context, cp *model.CondicaoPagamento) error
	}
	CategoriasContasPagar interface {
		CriarCategoria(ctx context.Context, c *model.CategoriaContaPagar) (*model.CategoriaContaPagar, error)
		ListarCategorias(ctx context.Context) ([]model.CategoriaContaPagar, error)
	}
	CategoriasContasReceber interface {
		CriarCategoria(context.Context, *model.CategoriaContaReceber) (*model.CategoriaContaReceber, error)
		ListarCategorias(context.Context) ([]model.CategoriaContaReceber, error)
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
	EntradaEstoque interface {
		RegistrarEntrada(ctx context.Context, entrada *model.EntradaEstoque) error
		ListarEntradas(ctx context.Context, filtro model.FiltroEntradaEstoque) ([]model.EntradaEstoque, error)
		ObterEntrada(ctx context.Context, id uint64) (*model.EntradaEstoque, error)
		EditarEntrada(ctx context.Context, entrada *model.EntradaEstoque) error
		AprovarEntrada(ctx context.Context, id uint64, usuarioID int64) error
		CancelarEntrada(ctx context.Context, id uint64, usuarioID int64) error
	}
	SaidaEstoque interface {
		RegistrarSaida(ctx context.Context, saida *model.SaidaEstoque) error
		ListarSaidas(ctx context.Context, filtro model.FiltroSaidaEstoque) ([]model.SaidaEstoque, error)
		ObterSaida(ctx context.Context, id uint64) (*model.SaidaEstoque, error)
		EditarSaida(ctx context.Context, saida *model.SaidaEstoque) error
		AprovarSaida(ctx context.Context, id uint64, usuarioID int64) error
		CancelarSaida(ctx context.Context, id uint64, usuarioID int64) error
	}
	Produtos interface {
		CriarProduto(ctx context.Context, input *model.ProdutoInput) (*model.ProdutoCompleto, error)
		ListarProdutos(ctx context.Context, busca string, idFornecedor int64) ([]*model.ProdutoCompleto, error)
		ObterProdutoPorID(ctx context.Context, id int64) (*model.ProdutoCompleto, error)
		AtualizarProduto(ctx context.Context, id int64, input *model.ProdutoInput) error
		ExcluirOuInativarProduto(ctx context.Context, id int64) (string, error)
		VincularProdutoEstoque(ctx context.Context, idProduto, idEstoque int64, qtdMinima float64) error
		DesvincularProdutoEstoque(ctx context.Context, idProduto, idEstoque int64) error
		ListarGruposTributarios(ctx context.Context) ([]*model.GrupoTributario, error)
		ListarComposicao(ctx context.Context, idProduto int64) ([]model.ItemComposicaoProduto, error)
		SalvarComposicao(ctx context.Context, idProduto int64, itens []model.ItemComposicaoProduto) error
	}
	Dashboard *DashboardService
}

func NewService(repository *repository.Repository, db *sql.DB) *Service {
	return &Service{
		VendaPDV:        &VendaPDVService{db: db},
		ConfiguracaoPDV: &ConfiguracaoPDVService{db: db},
		Caixa:           &CaixaService{db: db},
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
		CategoriasContasReceber: &CategoriaContaReceberService{repository: repository, db: db},
		ContasPagar: &ContaPagarService{
			repository: repository,
			db:         db,
		},
		ContasReceber: &ContaReceberService{repository: repository, db: db},
		CondicoesPagamento: &CondicaoPagamentoService{
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
		EntradaEstoque: &EntradaEstoqueService{
			repositorio: repository,
			db:          db,
		},
		SaidaEstoque: &SaidaEstoqueService{repository: repository, db: db},
		Produtos: &ProdutoService{
			repository: repository,
			db:         db,
		},
		Dashboard: NewDashboardService(repository.Dashboard, db),
	}
}
