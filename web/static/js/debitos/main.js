import { checkAuth } from '../utils/auth.js';
import { carregarDropdowns } from './dropdowns.js';
import { carregarDebitos } from './listarDebitos.js';
import { setupCriarNovoDebito } from './criarNovoDebito.js';
import { setupEditarDebito } from './editarDebito.js';
import { setupVisualizarDebito } from './visualizarDebito.js';
import { setupPagarDebito } from './pagarDebito.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    const formDebito = document.getElementById('formDebitoAvulso');
    const formFiltro = document.getElementById('formFiltroDebitos');
    const tabelaDebitos = document.getElementById('tabela_debitos_body');

    if (formDebito) {
        carregarDropdowns();
    }

    if (tabelaDebitos) {
        carregarDebitos();
    }

    if (formFiltro) {
        formFiltro.addEventListener('submit', (e) => {
            e.preventDefault();
            carregarDebitos();
        });
    }

    // Inicializar listeners e handlers globais
    setupCriarNovoDebito();
    setupEditarDebito();
    setupVisualizarDebito();
    setupPagarDebito();
});
