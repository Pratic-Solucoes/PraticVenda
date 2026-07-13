import { checkAuth } from '../utils/auth.js';
import { carregarEstoques, setupListarEstoque } from './listarEstoque.js';
import { setupCriarEstoque } from './criarEstoque.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    // Carregar estoques iniciais no dropdown
    carregarEstoques();

    // Configurar listeners
    setupListarEstoque();
    setupCriarEstoque();
});
