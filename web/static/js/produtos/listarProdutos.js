import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';
import { carregarProdutoParaEdicao, excluirOuInativarProduto } from './gerenciarProduto.js';

export async function carregarProdutos(busca = "") {
    const tbody = document.getElementById('tabelaProdutosBody');
    if (!tbody) return;

    tbody.innerHTML = '<tr><td colspan="8" class="text-center py-4"><div class="spinner-border text-primary" role="status"></div></td></tr>';

    const token = getToken();

    try {
        const url = busca ? `/api/produtos?busca=${encodeURIComponent(busca)}` : '/api/produtos';
        const res = await fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!res.ok) {
            const data = await res.json();
            showError(data.erro || "Erro ao buscar produtos.");
            tbody.innerHTML = '<tr><td colspan="8" class="text-center text-danger py-4">Erro ao carregar dados.</td></tr>';
            return;
        }

        const produtos = await res.json();

        tbody.innerHTML = '';

        if (produtos.length === 0) {
            tbody.innerHTML = '<tr><td colspan="8" class="text-center text-muted py-4">Nenhum produto cadastrado.</td></tr>';
            return;
        }

        produtos.forEach(p => {
            const tr = document.createElement('tr');
            
            // Estoques Vinculados badges
            let estoquesBadges = '<span class="text-muted small">Nenhum</span>';
            if (p.estoques && p.estoques.length > 0) {
                estoquesBadges = p.estoques.map(est => 
                    `<span class="badge bg-light text-dark border me-1 mb-1" title="Estoque Mínimo: ${est.estoque_minimo}">
                        ${est.nome_estoque}: <strong>${est.quantidade}</strong>
                     </span>`
                ).join('');
            }

            const badgeStatus = p.produto.ativo 
                ? '<span class="badge bg-success">Ativo</span>' 
                : '<span class="badge bg-danger">Inativo</span>';

            tr.innerHTML = `
                <td>${p.produto.codigo_interno_loja || '-'}</td>
                <td class="fw-semibold text-dark">${p.produto.nome}</td>
                <td>${p.produto.codigo_barras || '-'}</td>
                <td>R$ ${(p.produto.preco_custo || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 4 })}</td>
                <td>R$ ${(p.produto.preco_venda || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                <td>${estoquesBadges}</td>
                <td>${badgeStatus}</td>
                <td class="text-end">
                    <button class="btn btn-sm btn-outline-primary btn-editar-prod me-1" data-id="${p.produto.id}" title="Editar Produto">
                        <i class="bi bi-pencil-square"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-danger btn-excluir-prod" data-id="${p.produto.id}" title="Excluir/Inativar">
                        <i class="bi bi-trash"></i>
                    </button>
                </td>
            `;

            // Bind events
            tr.querySelector('.btn-editar-prod').addEventListener('click', () => {
                carregarProdutoParaEdicao(p.produto.id);
            });

            tr.querySelector('.btn-excluir-prod').addEventListener('click', () => {
                excluirOuInativarProduto(p.produto.id, p.produto.nome);
            });

            tbody.appendChild(tr);
        });

    } catch (err) {
        console.error(err);
        showError("Erro interno ao carregar a listagem de produtos.");
        tbody.innerHTML = '<tr><td colspan="8" class="text-center text-danger py-4">Erro interno ao carregar dados.</td></tr>';
    }
}
