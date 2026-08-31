package model

type ConfiguracaoPDV struct {
	IDEstoquePadrao         int64   `json:"id_estoque_padrao"`
	IDCategoriaCredito      int64   `json:"id_categoria_credito"`
	FormasPagamento         []int64 `json:"formas_pagamento"`
	CondicoesPagamento      []int64 `json:"condicoes_pagamento"`
	ExigirClientePrazo      bool    `json:"exigir_cliente_prazo"`
	PermitirDescontoManual  bool    `json:"permitir_desconto_manual"`
	PermitirAlterarPreco    bool    `json:"permitir_alterar_preco"`
	GerarFinanceiroRecebido bool    `json:"gerar_financeiro_recebido"`
}
