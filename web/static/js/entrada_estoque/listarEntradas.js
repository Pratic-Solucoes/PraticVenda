import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';

// ─── Listar Entradas ─────────────────────────────────────────────────────────

/** Busca entradas da API com os filtros aplicados */
async function fetchEntradas(filtros = {}) {
    const token = getToken();
    const tbody = document.getElementById('tabelaEntradasBody');

    tbody.innerHTML = `
        <tr>
            <td colspan="7" class="text-center py-5">
                <div class="spinner-border spinner-border-sm text-primary" role="status"></div>
                <span class="ms-2 text-muted">Carregando...</span>
            </td>
        </tr>
    `;

    try {
        // Monta query string com filtros
        const params = new URLSearchParams();
        if (filtros.id) params.append('id', filtros.id);
        if (filtros.fornecedor) params.append('fornecedor', filtros.fornecedor);
        if (filtros.data) params.append('data', filtros.data);
        if (filtros.status) params.append('status', filtros.status);

        // TODO: ajustar endpoint quando backend estiver pronto
        const res = await fetch(`/api/estoques/entradas?${params.toString()}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!res.ok) {
            const data = await res.json();
            showError(data.erro || 'Erro ao carregar entradas de estoque.');
            tbody.innerHTML = `<tr><td colspan="7" class="text-center text-danger py-5">Erro ao carregar dados.</td></tr>`;
            return;
        }

        const entradas = await res.json();
        popularTabelaEntradas(entradas);

    } catch (err) {
        // Endpoint ainda não existe — mostra estado vazio sem erro visual
        tbody.innerHTML = `
            <tr>
                <td colspan="7" class="text-muted py-5 text-center">
                    <i class="bi bi-inbox fs-3 d-block mb-2 text-muted"></i>
                    Nenhuma entrada encontrada.
                </td>
            </tr>
        `;
    }
}

/** Popula a tabela com os dados recebidos da API */
function popularTabelaEntradas(entradas) {
    const tbody = document.getElementById('tabelaEntradasBody');
    tbody.innerHTML = '';

    if (!entradas || entradas.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="7" class="text-muted py-5 text-center">
                    <i class="bi bi-inbox fs-3 d-block mb-2 text-muted"></i>
                    Nenhuma entrada encontrada.
                </td>
            </tr>
        `;
        return;
    }

    entradas.forEach(entrada => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td class="fw-semibold text-muted">#${entrada.id}</td>
            <td class="text-start fw-semibold">${entrada.fornecedor || '—'}</td>
            <td>R$ ${(entrada.valor_total || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</td>
            <td>${entrada.criado_em ? new Date(entrada.criado_em).toLocaleDateString('pt-BR') : '—'}</td>
            <td>${badgeStatus(entrada.status)}</td>
            <td>${entrada.usuario || '—'}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary" title="Ver detalhes" onclick="alert('Detalhes — ID ${entrada.id}')">
                    <i class="bi bi-eye"></i>
                </button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

function badgeStatus(status) {
    const map = {
        'ABERTA':    '<span class="badge bg-warning text-dark">Em Aberto</span>',
        'CONCLUIDA': '<span class="badge bg-success">Concluída</span>',
        'CANCELADA': '<span class="badge bg-danger">Cancelada</span>',
    };
    return map[status] || `<span class="badge bg-secondary">${status || '—'}</span>`;
}

// ─── Setup ───────────────────────────────────────────────────────────────────

export function setupListarEntradas() {
    // Carrega ao iniciar (estado vazio tolerado enquanto backend não existe)
    fetchEntradas();

    // Submit do formulário de filtros
    const form = document.getElementById('formFiltroEntradas');
    form?.addEventListener('submit', (e) => {
        e.preventDefault();
        const filtros = {
            id:         document.getElementById('filtro_entrada_id').value.trim(),
            fornecedor: document.getElementById('filtro_entrada_fornecedor').value.trim(),
            data:       document.getElementById('filtro_entrada_data').value,
            status:     document.getElementById('filtro_entrada_status').value,
        };
        fetchEntradas(filtros);
    });
}
