import { checkAuth } from '../utils/auth.js';
import { setupFormEntrada, setupBotaoNovaEntrada } from './formEntrada.js';
import { setupListarEntradas } from './listarEntradas.js';
import { setupFinanceiro } from './financeiro.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    // Inicializa o painel de nova entrada (collapse + listeners do form)
    setupBotaoNovaEntrada();
    setupFormEntrada();

    // Inicializa a tabela de entradas existentes
    setupListarEntradas();

    // Inicializa o modal de gerar financeiro
    setupFinanceiro();
});
