import { getToken } from '../utils/auth.js';
import { state } from './state.js';

export async function carregarContasPagar() {
    const tabelaDebitos = document.getElementById('tabela_debitos_body');
    if (!tabelaDebitos) return;
    
    const token = getToken();
    tabelaDebitos.innerHTML = '<tr><td colspan="7" class="text-muted py-5 text-center">Carregando...</td></tr>';

    const busca = document.getElementById('filtro_fornecedor')?.value || '';
    const vencimento = document.getElementById('filtro_vencimento')?.value || '';
    const status = document.getElementById('filtro_status')?.value || '';

    try {
        let url = '/api/contas-pagar?';
        const params = new URLSearchParams();
        if (busca) params.append('busca', busca);
        if (vencimento) params.append('vencimento', vencimento);
        if (status) params.append('status', status);
        url += params.toString();

        const res = await fetch(url, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (res.status === 401) {
            window.location.href = '/';
            return;
        }

        if (!res.ok) {
            tabelaDebitos.innerHTML = '<tr><td colspan="7" class="text-danger py-5 text-center">Erro ao carregar débitos.</td></tr>';
            return;
        }

        const dados = await res.json();
        state.contasPagarCarregadas = dados;

        const totalDebitos = dados.reduce((total, d) => total + d.saldo_restante, 0);
        const totalDebitosElement = document.getElementById('total_debitos');
        if (totalDebitosElement) {
            totalDebitosElement.textContent = new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(totalDebitos);
        }

        renderTabelaContasPagar(dados);
    } catch (err) {
        console.error(err);
        tabelaDebitos.innerHTML = '<tr><td colspan="7" class="text-danger py-5 text-center">Erro de comunicação.</td></tr>';
    }
}

function renderTabelaContasPagar(debitos) {
    const tabelaDebitos = document.getElementById('tabela_debitos_body');
    if (!tabelaDebitos) return;

    if (!debitos || debitos.length === 0) {
        tabelaDebitos.innerHTML = '<tr><td colspan="7" class="text-muted py-5 text-center">Nenhum débito encontrado.</td></tr>';
        return;
    }

    tabelaDebitos.innerHTML = '';
    debitos.forEach(d => {
        const dataVencimento = d.dt_vencimento ? d.dt_vencimento.substring(0, 10).split('-').reverse().join('/') : '-';
        const valorFormatado = new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(d.saldo_restante);

        let badgeClass = 'bg-secondary';
        if (d.status === 'PENDENTE') badgeClass = 'bg-warning text-dark';
        if (d.status === 'PAGO') badgeClass = 'bg-success';
        if (d.status === 'CANCELADO') badgeClass = 'bg-danger';

        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>#${d.id}</td>
            <td class="text-start fw-bold">${d.fornecedor ? d.fornecedor.razao_social : '-'}</td>
            <td class="text-start text-muted">${d.descricao}</td>
            <td>${dataVencimento}</td>
            <td class="fw-bold text-danger">${valorFormatado}</td>
            <td><span class="badge ${badgeClass} px-3 py-2 rounded-pill">${d.status}</span></td>
            <td>
                <button class="btn btn-sm btn-outline-primary" title="Visualizar" onclick="window.abrirModalVisualizacao(${d.id})"><i class="bi bi-eye"></i></button>
                <button class="btn btn-sm btn-outline-warning" title="Editar" onclick="window.abrirModalEdicao(${d.id})" ${d.status === 'PAGO' ? 'disabled' : ''}><i class="bi bi-pencil"></i></button>
                ${(d.status === 'PENDENTE' || d.status === 'PAGO_PARCIAL') ? `<button class="btn btn-sm btn-outline-success" title="Dar Baixa" onclick="window.abrirModalPagamento(${d.id})"><i class="bi bi-check2-circle"></i></button>` : ''}
            </td>
        `;
        tabelaDebitos.appendChild(tr);
    });
}
