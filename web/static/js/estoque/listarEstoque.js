import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';

export async function carregarEstoques(selecionarId = null) {
    const selectEstoque = document.getElementById('selectEstoque');
    if (!selectEstoque) return;

    const token = getToken();

    try {
        const res = await fetch('/api/estoques', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!res.ok) {
            const data = await res.json();
            showError(data.erro || "Erro ao carregar locais de estoque.");
            return;
        }

        const estoques = await res.json();

        // Limpar dropdown
        selectEstoque.innerHTML = '';

        if (estoques.length === 0) {
            selectEstoque.innerHTML = '<option value="" disabled selected>Nenhum estoque cadastrado</option>';
            document.getElementById('descEstoqueSelecionado').innerText = 'Cadastre um estoque para começar.';
            return;
        }

        estoques.forEach(est => {
            const opt = document.createElement('option');
            opt.value = est.id;
            opt.textContent = est.nome;
            opt.dataset.descricao = est.descricao || 'Sem descrição';
            opt.dataset.ativo = est.ativo;
            selectEstoque.appendChild(opt);
        });

        // Selecionar o ID desejado ou o primeiro
        if (selecionarId) {
            selectEstoque.value = selecionarId;
        } else {
            selectEstoque.selectedIndex = 0;
        }

        // Atualizar informações e carregar os produtos
        atualizarDetalhesEstoque();
        const idSelecionado = selectEstoque.value;
        if (idSelecionado) {
            await carregarProdutosDoEstoque(idSelecionado);
        }

    } catch (err) {
        console.error(err);
        showError("Erro interno ao carregar estoques.");
    }
}

export async function carregarProdutosDoEstoque(idEstoque) {
    const tbody = document.getElementById('tabelaProdutosEstoqueBody');
    if (!tbody) return;

    tbody.innerHTML = '<tr><td colspan="7" class="text-center py-4"><div class="spinner-border text-primary" role="status"></div></td></tr>';

    const token = getToken();

    try {
        const res = await fetch(`/api/estoques/${idEstoque}/produtos`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!res.ok) {
            const data = await res.json();
            showError(data.erro || "Erro ao buscar produtos do estoque.");
            tbody.innerHTML = '<tr><td colspan="7" class="text-center text-danger py-4">Erro ao carregar dados.</td></tr>';
            return;
        }

        const produtosEstoque = await res.json();

        tbody.innerHTML = '';

        if (produtosEstoque.length === 0) {
            tbody.innerHTML = '<tr><td colspan="7" class="text-center text-muted py-4">Nenhum produto cadastrado neste estoque.</td></tr>';
            return;
        }

        produtosEstoque.forEach(pe => {
            const prod = pe.produto || {};
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${prod.codigo_interno_loja || '-'}</td>
                <td>${prod.codigo_barras || '-'}</td>
                <td class="fw-semibold">${prod.nome || 'Produto sem nome'}</td>
                <td>R$ ${(prod.preco_venda || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                <td>${(pe.estoque_minimo || 0).toFixed(0)}</td>
                <td class="text-center fw-bold text-primary">${(pe.quantidade || 0).toFixed(0)}</td>
                <td><span class="badge bg-secondary">${prod.unidade_estoque || 'UN'}</span></td>
            `;
            tbody.appendChild(tr);
        });

    } catch (err) {
        console.error(err);
        showError("Erro interno ao buscar produtos do estoque.");
        tbody.innerHTML = '<tr><td colspan="7" class="text-center text-danger py-4">Erro interno ao carregar dados.</td></tr>';
    }
}

export function atualizarDetalhesEstoque() {
    const selectEstoque = document.getElementById('selectEstoque');
    const badgeStatus = document.getElementById('badgeStatusEstoque');
    const descEstoque = document.getElementById('descEstoqueSelecionado');

    if (!selectEstoque || !badgeStatus || !descEstoque) return;

    const optSelecionada = selectEstoque.options[selectEstoque.selectedIndex];
    if (optSelecionada) {
        descEstoque.innerText = optSelecionada.dataset.descricao || 'Sem descrição';
        const ativo = optSelecionada.dataset.ativo === 'true';
        badgeStatus.className = ativo ? 'badge bg-success p-2 fs-6' : 'badge bg-danger p-2 fs-6';
        badgeStatus.innerText = ativo ? 'Ativo' : 'Inativo';
    }
}

export function setupListarEstoque() {
    const selectEstoque = document.getElementById('selectEstoque');
    if (!selectEstoque) return;

    selectEstoque.addEventListener('change', () => {
        atualizarDetalhesEstoque();
        const idEstoque = selectEstoque.value;
        if (idEstoque) {
            carregarProdutosDoEstoque(idEstoque);
        }
    });
}
