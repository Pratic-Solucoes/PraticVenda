import { checkAuth } from '../utils/auth.js';
import { carregarCategorias } from './listarCategorias.js';
import { setupCriarCategoria } from './criarCategoria.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    const tbody = document.getElementById('tabela_categorias_body');
    
    if (tbody) {
        carregarCategorias();
    }

    setupCriarCategoria();
});
