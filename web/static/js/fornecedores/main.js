import { checkAuth } from '../utils/auth.js';
import { carregarFornecedores } from './listarFornecedores.js';
import { setupEditarFornecedor } from './editarFornecedor.js';
import { setupCriarFornecedorModal } from './criarFornecedorModal.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    const tbody = document.getElementById('tabela_fornecedores_body');
    const formFiltro = document.getElementById('formFiltroFornecedores');
    const inputBusca = document.getElementById('filtro_busca');

    // Elements for Inline Form
    const btnNovoFornecedor = document.getElementById("btnNovoFornecedor");
    const btnFecharFormFornecedor = document.getElementById("btnFecharFormFornecedor");
    const btnCancelarFornecedor = document.getElementById("btnCancelarFornecedor");
    const collapseFormFornecedor = document.getElementById("collapseFormFornecedor");
    const formFornecedorTitulo = document.getElementById("formFornecedorTitulo");
    const formEditarFornecedor = document.getElementById("formInlineFornecedor");
    const editFornecedorId = document.getElementById("edit_fornecedor_id");

    let bsCollapse = null;
    if (collapseFormFornecedor) {
        bsCollapse = new bootstrap.Collapse(collapseFormFornecedor, { toggle: false });
    }

    // Carregar fornecedores iniciais
    if (tbody) {
        carregarFornecedores();
    }

    // Filtro
    if (formFiltro) {
        formFiltro.addEventListener('submit', (e) => {
            e.preventDefault();
            carregarFornecedores(inputBusca?.value.trim() || "");
        });
    }

    // Handle Form toggle
    function abrirFormularioNovo() {
        if(formFornecedorTitulo) formFornecedorTitulo.innerHTML = '<i class="bi bi-truck me-2"></i> Novo Fornecedor';
        if(formEditarFornecedor) formEditarFornecedor.reset();
        if(editFornecedorId) editFornecedorId.value = "";
        
        // Hide tabs except Dados Básicos for new creation
        const tabEnderecos = document.getElementById("enderecos-tab");
        const tabTelefones = document.getElementById("telefones-tab");
        if (tabEnderecos) tabEnderecos.parentElement.style.display = "none";
        if (tabTelefones) tabTelefones.parentElement.style.display = "none";

        const basicoTab = document.getElementById("basico-tab");
        if (basicoTab) bootstrap.Tab.getOrCreateInstance(basicoTab).show();

        if (bsCollapse) bsCollapse.show();
    }

    function fecharFormulario() {
        if (bsCollapse) bsCollapse.hide();
        if(formEditarFornecedor) formEditarFornecedor.reset();
        if(editFornecedorId) editFornecedorId.value = "";
    }

    if (btnNovoFornecedor) {
        btnNovoFornecedor.addEventListener("click", abrirFormularioNovo);
    }

    if (btnFecharFormFornecedor) {
        btnFecharFormFornecedor.addEventListener("click", fecharFormulario);
    }

    if (btnCancelarFornecedor) {
        btnCancelarFornecedor.addEventListener("click", fecharFormulario);
    }

    // Make fecharFormulario global for other files to use if needed
    window.fecharFormularioInlineFornecedor = fecharFormulario;
    window.abrirFormularioInlineFornecedor = function() {
        if (bsCollapse) bsCollapse.show();
        
        // Show all tabs when editing
        const tabEnderecos = document.getElementById("enderecos-tab");
        const tabTelefones = document.getElementById("telefones-tab");
        if (tabEnderecos) tabEnderecos.parentElement.style.display = "block";
        if (tabTelefones) tabTelefones.parentElement.style.display = "block";
    };

    setupEditarFornecedor();
    setupCriarFornecedorModal();
});
