import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';

export async function carregarFormasPagamento() {
    const tbody = document.getElementById('tabela_formas_pagamento_body');
    if (!tbody) return;
    
    const token = getToken();
    tbody.innerHTML = `<tr><td colspan="3" class="text-muted py-5 text-center">Carregando...</td></tr>`;
    
    try {
        const res = await fetch('/api/formas-pagamento', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (res.status === 401) {
            window.location.href = '/';
            return;
        }

        if (!res.ok) {
            showError("Erro ao carregar lista de formas de pagamento.");
            tbody.innerHTML = `<tr><td colspan="3" class="text-danger py-5 text-center">Erro ao carregar.</td></tr>`;
            return;
        }

        const dados = await res.json();
        renderTabela(dados);

    } catch (err) {
        console.error(err);
        tbody.innerHTML = `<tr><td colspan="3" class="text-danger py-5 text-center">Erro de comunicação.</td></tr>`;
    }
}

function renderTabela(formas) {
    const tbody = document.getElementById('tabela_formas_pagamento_body');
    if (!tbody) return;

    if (!formas || formas.length === 0) {
        tbody.innerHTML = `<tr><td colspan="3" class="text-muted py-5 text-center">Nenhuma forma de pagamento encontrada.</td></tr>`;
        return;
    }

    tbody.innerHTML = '';
    formas.forEach(f => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>#${f.id}</td>
            <td class="text-start fw-bold">${f.descricao}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary" title="Editar" onclick="abrirEditarFormaPagamento(${f.id})">
                    <i class="bi bi-pencil"></i>
                </button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}
