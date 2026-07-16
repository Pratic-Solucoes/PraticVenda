package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

type ViewController struct {
}

func (c *ViewController) RenderizarLoginPage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/template/pages/login.html",
		"web/template/components/loginForm.html",
		"web/template/components/toastContainer.html",
	)

	if err != nil {
		fmt.Printf("Erro ao renderizar página de login: %v\n", err)
		http.Error(w, "Erro ao renderizar página de login", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao renderizar página de login: %v\n", err)
		http.Error(w, "Erro ao renderizar página de login", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarDashboardPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/dashboard.html",
		"web/template/components/sidebar.html",
		"web/template/components/modalContaPagarAvulso.html",
		"web/template/components/modalFornecedor.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar dashboard: %v\n", err)
		http.Error(w, "Erro ao renderizar dashboard: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao renderizar dashboard: %v\n", err)
		http.Error(w, "Erro ao renderizar dashboard", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarContasPagarPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/contas_pagar.html",
		"web/template/components/sidebar.html",
		"web/template/components/modalContaPagarAvulso.html",
		"web/template/components/modalEditarContaPagar.html",
		"web/template/components/modalVisualizarContaPagar.html",
		"web/template/components/modalPagarConta.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de contas a pagar: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de contas a pagar: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarClientesPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/clientes.html",
		"web/template/components/sidebar.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de clientes: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de clientes: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarFornecedoresPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/fornecedores.html",
		"web/template/components/sidebar.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de fornecedores: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de fornecedores: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarCategoriasPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/categorias_contas_pagar.html",
		"web/template/components/sidebar.html",
		"web/template/components/modalCategoria.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de categorias: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de categorias: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarConfiguracaoUsuarioPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/configuracao_usuario.html",
		"web/template/components/sidebar.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de configuração do usuário: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de configuração do usuário: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarPdvPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/pdv.html",
		"web/template/components/sidebar.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de PDV: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de PDV: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarFormasPagamentoPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/formas_pagamento.html",
		"web/template/components/sidebar.html",
		"web/template/components/modalFormaPagamento.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de formas de pagamento: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de formas de pagamento: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarEstoquesPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/estoques.html",
		"web/template/components/sidebar.html",
		"web/template/components/modalEstoque.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de estoques: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de estoques: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarProdutosPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/produtos.html",
		"web/template/components/sidebar.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de produtos: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de produtos: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarEntradaEstoquePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/entrada_estoque.html",
		"web/template/components/sidebar.html",
		"web/template/components/modalFinanceiroEntrada.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de entrada de estoque: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de entrada de estoque: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarCondicoesPagamentoPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/condicoes_pagamento.html",
		"web/template/components/sidebar.html",
		"web/template/components/toastContainer.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de condições de pagamento: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de condições de pagamento: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

