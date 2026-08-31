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
    const dropdown = document.getElementById('dropdown_produto');
    const inputProduto = document.getElementById('pdv_produto');

    if (!query) {
        if (dropdown) dropdown.classList.add('d-none');
        return;
    }

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
            if (produtos.length === 1 && (produtos[0].produto.codigo_barras === query || produtos[0].produto.id.toString() === query)) {
                // Se for bipe exato (1 resultado) adiciona direto
                adicionarProduto(produtos[0].produto || produtos[0]);
                inputProduto.value = '';
                if (dropdown) dropdown.classList.add('d-none');
            } else {
                // Preenche o dropdown
                if (dropdown) {
                    dropdown.innerHTML = '';
                    produtos.forEach(p => {
                        const prod = p.produto || p;
                        const btn = document.createElement('button');
                        btn.type = 'button';
                        btn.className = 'list-group-item list-group-item-action py-2';
                        btn.innerHTML = `<strong>${prod.codigo_barras || prod.id}</strong> - ${prod.nome} <span class="float-end text-success fw-bold">R$ ${(prod.preco_venda || 0).toFixed(2).replace('.', ',')}</span>`;
                        btn.addEventListener('click', () => {
                            adicionarProduto(prod);
                            inputProduto.value = '';
                            dropdown.classList.add('d-none');
                        });
                        dropdown.appendChild(btn);
                    });
                    dropdown.classList.remove('d-none');
                }
            }
        } else {
            if (dropdown) {
                dropdown.innerHTML = '<div class="list-group-item text-muted">Nenhum produto encontrado.</div>';
                dropdown.classList.remove('d-none');
            }
        }

    } catch (err) {
        console.error(err);
        notify("Erro interno ao buscar produto.", "error");
    }
}

function adicionarProduto(p) {
    const produtoData = p.produto || p;
    adicionarItem(produtoData, 1);
    renderizarTabela();
    atualizarTotalDisplay();
    notify(`${produtoData.nome} adicionado!`, "success");
}
