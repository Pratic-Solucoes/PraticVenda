import { getToken } from '../utils/auth.js';
import { showError, showSuccess } from '../utils/showError.js';
import { abrirEntradaParaEdicao } from './formEntrada.js';

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

        const res = await fetch(`/api/entradas-estoque?${params.toString()}`, {
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
        console.error(err);
        tbody.innerHTML = `
            <tr>
                <td colspan="7" class="text-danger py-5 text-center">
                    Erro de comunicação ao carregar entradas de estoque.
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
            <td class="text-nowrap">
                <button class="btn btn-sm btn-outline-primary btn-ver-entrada" data-id="${entrada.id}" title="Visualizar"><i class="bi bi-eye"></i></button>
                ${entrada.status === 'ABERTO' ? `<button class="btn btn-sm btn-outline-secondary btn-editar-entrada" data-id="${entrada.id}" title="Editar"><i class="bi bi-pencil"></i></button>
                <button class="btn btn-sm btn-success btn-aprovar-entrada" data-id="${entrada.id}" title="Aprovar"><i class="bi bi-check-lg"></i></button>` : ''}
            </td>
        `;
		tr.querySelector('.btn-ver-entrada').addEventListener('click', () => visualizarEntrada(entrada.id));
		tr.querySelector('.btn-editar-entrada')?.addEventListener('click', () => editarEntrada(entrada.id));
		tr.querySelector('.btn-aprovar-entrada')?.addEventListener('click', () => aprovarEntrada(entrada.id));
        tbody.appendChild(tr);
    });
}

/** Busca uma entrada completa, necessária para visualizar ou editar seus itens. */
async function buscarEntrada(id) {
	const res = await fetch(`/api/entradas-estoque/${id}`, { headers: { Authorization: `Bearer ${getToken()}` } });
	if (!res.ok) { const data = await res.json(); throw new Error(data.erro || 'Não foi possível carregar a entrada.'); }
	return res.json();
}

function escapar(valor) { return String(valor ?? '').replace(/[&<>"']/g, char => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#039;' }[char])); }
function moeda(valor) { return Number(valor || 0).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }); }

/** Exibe o resumo e os itens da entrada no modal de detalhes. */
async function visualizarEntrada(id) {
	try {
		const entrada = await buscarEntrada(id);
		document.getElementById('tituloDetalheEntrada').textContent = `Entrada #${entrada.id}`;
		document.getElementById('corpoDetalheEntrada').innerHTML = `<div class="row g-2 mb-3"><div class="col-md-6"><strong>Fornecedor:</strong> ${escapar(entrada.fornecedor)}</div><div class="col-md-3"><strong>Status:</strong> ${badgeStatus(entrada.status)}</div><div class="col-md-3"><strong>Total:</strong> ${moeda(entrada.valor_total)}</div></div><div class="table-responsive"><table class="table table-sm"><thead><tr><th>Produto</th><th>Qtd.</th><th class="text-end">Custo</th><th class="text-end">Total</th></tr></thead><tbody>${(entrada.produtos || []).map(p => `<tr><td>${escapar(p.nome_produto || `Produto #${p.id_produto}`)}</td><td>${p.quantidade}</td><td class="text-end">${moeda(p.valor_custo)}</td><td class="text-end">${moeda(p.valor_total)}</td></tr>`).join('') || '<tr><td colspan="4" class="text-center text-muted">Sem itens.</td></tr>'}</tbody></table></div>`;
		const acoes = document.getElementById('acoesDetalheEntrada');
		acoes.innerHTML = '<button type="button" class="btn btn-light" data-bs-dismiss="modal">Fechar</button>';
		if (entrada.status === 'ABERTO') acoes.insertAdjacentHTML('afterbegin', '<button type="button" class="btn btn-outline-primary" id="btnEditarDetalhe">Editar</button><button type="button" class="btn btn-success" id="btnAprovarDetalhe">Aprovar</button>');
		if (entrada.status !== 'CANCELADA') acoes.insertAdjacentHTML('afterbegin', '<button type="button" class="btn btn-outline-danger" id="btnCancelarDetalhe">Cancelar</button>');
		document.getElementById('btnEditarDetalhe')?.addEventListener('click', () => editarEntrada(id));
		document.getElementById('btnAprovarDetalhe')?.addEventListener('click', () => aprovarEntrada(id));
		document.getElementById('btnCancelarDetalhe')?.addEventListener('click', () => cancelarEntrada(id));
		bootstrap.Modal.getOrCreateInstance(document.getElementById('modalDetalheEntrada')).show();
	} catch (err) { showError(err.message); }
}

async function cancelarEntrada(id) { if (!confirm('Cancelar esta entrada? Uma entrada concluída terá seu saldo estornado.')) return; try { const res=await fetch(`/api/entradas-estoque/${id}/cancelar`,{method:'POST',headers:{Authorization:`Bearer ${getToken()}`}});const data=await res.json();if(!res.ok)throw new Error(data.erro||'Não foi possível cancelar a entrada.');bootstrap.Modal.getInstance(document.getElementById('modalDetalheEntrada'))?.hide();showSuccess(data.mensagem);fetchEntradas(); } catch(err){showError(err.message)} }

async function editarEntrada(id) {
	try {
		const entrada = await buscarEntrada(id);
		bootstrap.Modal.getInstance(document.getElementById('modalDetalheEntrada'))?.hide();
		await abrirEntradaParaEdicao(entrada);
	} catch (err) { showError(err.message); }
}

/** Confirma e aprova a entrada; a aprovação atualiza o saldo do estoque. */
async function aprovarEntrada(id) {
	if (!window.confirm('Aprovar esta entrada? Essa ação atualizará o saldo do estoque e não poderá ser desfeita.')) return;
	try {
		const res = await fetch(`/api/entradas-estoque/${id}/aprovar`, { method: 'POST', headers: { Authorization: `Bearer ${getToken()}` } });
		const data = await res.json();
		if (!res.ok) throw new Error(data.erro || 'Não foi possível aprovar a entrada.');
		bootstrap.Modal.getInstance(document.getElementById('modalDetalheEntrada'))?.hide();
		showSuccess(data.mensagem || 'Entrada aprovada com sucesso.');
		fetchEntradas();
	} catch (err) { showError(err.message); }
}

function badgeStatus(status) {
    const map = {
        'ABERTO':    '<span class="badge bg-warning text-dark">Em Aberto</span>',
        'CONCLUIDA': '<span class="badge bg-success">Concluída</span>',
        'CANCELADA': '<span class="badge bg-danger">Cancelada</span>',
    };
    return map[status] || `<span class="badge bg-secondary">${status || '—'}</span>`;
}

// ─── Setup ───────────────────────────────────────────────────────────────────

	export function setupListarEntradas() {
    // Carrega ao iniciar.
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
	window.addEventListener('entradas-atualizadas', () => fetchEntradas());
}
