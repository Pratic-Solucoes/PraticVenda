import { checkAuth } from '../utils/auth.js';
import { renderizarTabela, atualizarTotalDisplay } from './ui.js';
import { buscarEAdicionarProduto } from './produto.js';
import { removerItem, obterItens } from './carrinho.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    const inputProduto = document.getElementById('pdv_produto');
    const btnAddProduto = document.getElementById('btn_adicionar_produto');
    const tbody = document.getElementById('tabela_pdv_body');

    // 1. Inicializar UI limpa
    renderizarTabela();
    atualizarTotalDisplay();

    // 2. Handler para buscar o produto
    async function handleAddProduto() {
        const query = inputProduto.value.trim();
        if (query) {
            await buscarEAdicionarProduto(query);
            inputProduto.value = '';
            inputProduto.focus();
        }
    }

    // 3. Eventos do Input de Busca de Produto
    if (btnAddProduto) {
        btnAddProduto.addEventListener('click', handleAddProduto);
    }

    if (inputProduto) {
        inputProduto.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                handleAddProduto();
            }
        });
    }

    // 4. Delegação de Eventos para os botões dentro da Tabela
    if (tbody) {
        // Remover Item
        tbody.addEventListener('click', (e) => {
            const btnRemover = e.target.closest('.btn-remover');
            if (btnRemover) {
                const index = parseInt(btnRemover.getAttribute('data-index'), 10);
                removerItem(index);
                renderizarTabela();
                atualizarTotalDisplay();
            }
        });

        // Alterar Quantidade Manulamente
        tbody.addEventListener('change', (e) => {
            if (e.target.classList.contains('input-qtd')) {
                const index = parseInt(e.target.getAttribute('data-index'), 10);
                const novaQtd = parseInt(e.target.value, 10);
                
                if (novaQtd > 0) {
                    const itens = obterItens();
                    const item = itens[index];
                    
                    const produtoData = item.produto.produto || item.produto;
                    const preco = produtoData.preco_venda || 0;

                    item.quantidade = novaQtd;
                    item.subtotal = novaQtd * preco;
                    
                    renderizarTabela();
                    atualizarTotalDisplay();
                } else {
                    e.target.value = 1; // Reseta se colocar número inválido
                }
            }
        });
    }
});
