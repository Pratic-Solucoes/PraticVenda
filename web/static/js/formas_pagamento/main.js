import { checkAuth } from '../utils/auth.js';
import { carregarFormasPagamento } from './listarFormasPagamento.js';
import { setupCriarFormaPagamento } from './criarFormaPagamento.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    const tbody = document.getElementById('tabela_formas_pagamento_body');
    
    if (tbody) {
        carregarFormasPagamento();
    }

    setupCriarFormaPagamento();
});
