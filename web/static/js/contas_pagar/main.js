import { checkAuth } from '../utils/auth.js';
import { carregarDropdowns } from './dropdowns.js';
import { carregarContasPagar } from './listarContasPagar.js';
import { setupCriarContaPagar } from './criarContaPagar.js';
import { setupEditarContaPagar } from './editarContaPagar.js';
import { setupVisualizarContaPagar } from './visualizarContaPagar.js';
import { setupPagarContaPagar } from './pagarContaPagar.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    const formContaPagar = document.getElementById('formContaPagarAvulso');
    const formFiltro = document.getElementById('formFiltroDebitos');
    const tabelaContasPagar = document.getElementById('tabela_debitos_body');

    if (formContaPagar) {
        carregarDropdowns();
    }

    if (tabelaContasPagar) {
        carregarContasPagar();
    }

    if (formFiltro) {
        formFiltro.addEventListener('submit', (e) => {
            e.preventDefault();
            carregarContasPagar();
        });
    }

    // Inicializar listeners e handlers globais
    setupCriarContaPagar();
    setupEditarContaPagar();
    setupVisualizarContaPagar();
    setupPagarContaPagar();
});
