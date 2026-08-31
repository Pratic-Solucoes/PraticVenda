import { checkAuth } from '../utils/auth.js';
import { carregarProdutos } from './listarProdutos.js?v=2';
import { setupGerenciarProduto, abrirFormularioNovo, fecharFormulario } from './gerenciarProduto.js?v=2';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    // Elements for Inline Form
    const btnNovoProduto = document.getElementById("btnNovoProduto");
    const btnFecharFormProduto = document.getElementById("btnFecharFormProduto");
    const btnCancelarProduto = document.getElementById("btnCancelarProduto");
    const collapseFormProduto = document.getElementById("collapseFormProduto");
    const inputBusca = document.getElementById('buscaProduto');
    const btnLimparBusca = document.getElementById('btnLimparBusca');

    let bsCollapse = null;
    if (collapseFormProduto) {
        bsCollapse = new bootstrap.Collapse(collapseFormProduto, { toggle: false });
    }

    // Carregar lista de produtos inicialmente
    carregarProdutos();

    // Configurar busca em tempo real
    if (inputBusca) {
        let timeout = null;
        inputBusca.addEventListener('input', () => {
            clearTimeout(timeout);
            timeout = setTimeout(() => {
                carregarProdutos(inputBusca.value.trim());
            }, 300);
        });
    }

    if (btnLimparBusca) {
        btnLimparBusca.addEventListener('click', () => {
            if (inputBusca) inputBusca.value = '';
            carregarProdutos();
        });
    }

    if (btnNovoProduto) {
        btnNovoProduto.addEventListener("click", () => abrirFormularioNovo(bsCollapse));
    }

    if (btnFecharFormProduto) {
        btnFecharFormProduto.addEventListener("click", () => fecharFormulario(bsCollapse));
    }

    if (btnCancelarProduto) {
        btnCancelarProduto.addEventListener("click", () => fecharFormulario(bsCollapse));
    }

    // Exportar para escopo global para que possamos abrir a edição a partir da listagem
    window.abrirFormularioEditarProduto = () => {
        if (bsCollapse) bsCollapse.show();
    };

    setupGerenciarProduto();
});
