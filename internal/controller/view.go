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
		"web/template/components/errorModal.html",
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
		"web/template/components/modalDebitoAvulso.html",
		"web/template/components/modalFornecedor.html",
		"web/template/components/errorModal.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar dashboard: %v\n", err)
		http.Error(w, "Erro ao renderizar dashboard", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao renderizar dashboard: %v\n", err)
		http.Error(w, "Erro ao renderizar dashboard", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarDebitosPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/debitos.html",
		"web/template/components/sidebar.html",
		"web/template/components/modalDebitoAvulso.html",
		"web/template/components/modalEditarDebito.html",
		"web/template/components/modalVisualizarDebito.html",
	)
	if err != nil {
		fmt.Printf("Erro ao renderizar página de débitos: %v\n", err)
		http.Error(w, "Erro interno ao renderizar página", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Erro ao executar template de débitos: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
}

func (c *ViewController) RenderizarClientesPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/template/pages/clientes.html",
		"web/template/components/sidebar.html",
		"web/template/components/modalCliente.html",
		"web/template/components/modalEditarCliente.html",
		"web/template/components/errorModal.html",
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
		"web/template/components/modalFornecedor.html",
		"web/template/components/modalEditarFornecedor.html",
		"web/template/components/errorModal.html",
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
		"web/template/pages/categorias_debito.html",
		"web/template/components/sidebar.html",
		"web/template/components/modalCategoria.html",
		"web/template/components/errorModal.html",
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
