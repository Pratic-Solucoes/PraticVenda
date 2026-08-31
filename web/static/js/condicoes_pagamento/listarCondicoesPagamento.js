import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';

export async function carregarCondicoesPagamento() {
    const tbody = document.getElementById('tabela_condicoes_pagamento_body');
    if (!tbody) return;

    const token = getToken();
    tbody.innerHTML = `<tr><td colspan="6" class="text-muted py-5 text-center">Carregando...</td></tr>`;

    try {
        const res = await fetch('/api/condicoes-pagamento', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (res.status === 401) {
            window.location.href = '/';
            return;
        }

        if (!res.ok) {
            showError("Erro ao carregar lista de condições de pagamento.");
            tbody.innerHTML = `<tr><td colspan="6" class="text-danger py-5 text-center">Erro ao carregar.</td></tr>`;
            return;
        }

        const dados = await res.json();
        renderTabela(dados);

    } catch (err) {
        console.error(err);
        tbody.innerHTML = `<tr><td colspan="6" class="text-danger py-5 text-center">Erro de comunicação.</td></tr>`;
    }
}

function renderTabela(condicoes) {
    const tbody = document.getElementById('tabela_condicoes_pagamento_body');
    if (!tbody) return;

    if (!condicoes || condicoes.length === 0) {
        tbody.innerHTML = `<tr><td colspan="6" class="text-muted py-5 text-center">Nenhuma condição de pagamento encontrada.</td></tr>`;
        return;
    }

    tbody.innerHTML = '';
    condicoes.forEach(c => {
        const tr = document.createElement('tr');
        const valores = [`#${c.id}`, c.descricao, c.qtd_parcelas, c.dias_primeiro_venc, c.intervalo_parcelas];
        valores.forEach((valor, indice) => {
            const td = document.createElement('td');
            td.textContent = valor ?? '-';
            if (indice === 1) td.className = 'text-start fw-bold';
            tr.appendChild(td);
        });

        const tdAcoes = document.createElement('td');
        const botaoEditar = document.createElement('button');
        botaoEditar.type = 'button';
        botaoEditar.className = 'btn btn-sm btn-outline-primary';
        botaoEditar.title = 'Editar';
        botaoEditar.innerHTML = '<i class="bi bi-pencil"></i>';
        botaoEditar.addEventListener('click', () => window.abrirEditarCondicaoPagamento(c.id));
        tdAcoes.appendChild(botaoEditar);
        tr.appendChild(tdAcoes);
        tbody.appendChild(tr);
    });
}
