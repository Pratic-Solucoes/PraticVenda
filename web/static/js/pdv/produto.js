import { getToken } from '../utils/auth.js';
// Utilizamos uma função genérica de notificação caso exista, ou apenas o console por enquanto
function notify(msg, type) {
    console.log(`[${type}] ${msg}`);
    // Se existir o showError global no projeto:
    if (window.showError && type === 'error') window.showError(msg);
}

import { adicionarItem } from './carrinho.js';
import { renderizarTabela, atualizarTotalDisplay } from './ui.js';

export async function buscarEAdicionarProduto(query) {
    if (!query) return;

    try {
        const res = await fetch(`/api/produtos?busca=${encodeURIComponent(query)}`, {
            headers: { 'Authorization': `Bearer ${getToken()}` }
        });

        if (!res.ok) {
            notify("Erro ao buscar produto.", "error");
            return;
        }

        const produtos = await res.json();
        
        if (produtos && produtos.length > 0) {
            // Pega o primeiro produto retornado para agilizar a venda
            const produto = produtos[0]; 
            const produtoData = produto.produto || produto;
            
            adicionarItem(produto, 1);
            renderizarTabela();
            atualizarTotalDisplay();
            notify(`${produtoData.nome} adicionado!`, "success");
        } else {
            notify("Produto não encontrado.", "error");
        }

    } catch (err) {
        console.error(err);
        notify("Erro interno ao buscar produto.", "error");
    }
}
