import { checkAuth } from '../utils/auth.js';
import { carregarFornecedores } from './listarFornecedores.js';
import { setupCriarFornecedor } from './criarFornecedor.js';
import { setupEditarFornecedor } from './editarFornecedor.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    const tbody = document.getElementById('tabela_fornecedores_body');
    const formFiltro = document.getElementById('formFiltroFornecedores');
    const inputBusca = document.getElementById('filtro_busca');

    // Carregar fornecedores iniciais
    if (tbody) {
        carregarFornecedores();
    }

    // Filtro
    if (formFiltro) {
        formFiltro.addEventListener('submit', (e) => {
            e.preventDefault();
            carregarFornecedores(inputBusca.value);
        });
    }

    setupCriarFornecedor();
    setupEditarFornecedor();
});
