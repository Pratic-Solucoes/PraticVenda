import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';

// ─── Estado interno do formulário ────────────────────────────────────────────

/** @type {{ idProduto: number, nome: string, quantidade: number, valorUnitario: number, icmsSt: number, ipi: number, desconto: number, rateio: number, valorCusto: number, valorTotal: number }[]} */
let itensDaEntrada = [];

/** Produto atualmente carregado no card de busca */
let produtoSelecionado = null;

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatBRL(valor) {
    return valor.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

/** Retorna a despesa adicional informada no cabeçalho */
function getDespesaAdicional() {
    return parseFloat(document.getElementById('entrada_despesa_adicional').value) || 0;
}

/**
 * Recalcula o rateio de despesa adicional para cada item e atualiza a tabela.
 * Rateio proporcional ao valor (unitario * quantidade) de cada item.
 */
function recalcularRateios() {
    const despesa = getDespesaAdicional();
    const totalBruto = itensDaEntrada.reduce((acc, it) => acc + (it.valorUnitario * it.quantidade), 0);

    itensDaEntrada.forEach(item => {
        const proporcao = totalBruto > 0 ? (item.valorUnitario * item.quantidade) / totalBruto : 0;
        item.rateio = parseFloat((despesa * proporcao).toFixed(4));
        item.valorCusto = calcularValorCusto(item);
        item.valorTotal = parseFloat((item.valorCusto * item.quantidade).toFixed(2));
    });

    renderizarTabelaItens();
}

/**
 * Fórmula: ((vlr_unitario * qtd) + icms_st + ipi - desconto + rateio) / qtd
 */
function calcularValorCusto(item) {
    const { valorUnitario, quantidade, icmsSt, ipi, desconto, rateio } = item;
    if (!quantidade || quantidade <= 0) return 0;
    const custo = ((valorUnitario * quantidade) + icmsSt + ipi - desconto + rateio) / quantidade;
    return parseFloat(custo.toFixed(4));
}

/** Recalcula o valor de custo e total a partir dos campos do card */
function recalcularCamposCard() {
    const qtd = parseFloat(document.getElementById('ps_quantidade').value) || 0;
    const vlrUnit = parseFloat(document.getElementById('ps_valor_unitario').value) || 0;
    const icmsSt = parseFloat(document.getElementById('ps_icms_st').value) || 0;
    const ipi = parseFloat(document.getElementById('ps_ipi').value) || 0;
    const desconto = parseFloat(document.getElementById('ps_desconto').value) || 0;

    // Rateio provisório: despesa adicional total / (itens já adicionados + 1)
    const despesa = getDespesaAdicional();
    const divisor = itensDaEntrada.length + 1;
    const rateioProvisorio = parseFloat((despesa / divisor).toFixed(4));

    document.getElementById('ps_rateio').value = rateioProvisorio.toFixed(2);

    const itemTemp = { valorUnitario: vlrUnit, quantidade: qtd, icmsSt, ipi, desconto, rateio: rateioProvisorio };
    const custo = calcularValorCusto(itemTemp);
    const total = parseFloat((custo * qtd).toFixed(2));

    document.getElementById('ps_valor_custo').value = custo > 0 ? custo.toFixed(4) : '';
    document.getElementById('ps_valor_total').value = total > 0 ? total.toFixed(2) : '';
}

// ─── Renderização da tabela de itens ─────────────────────────────────────────

function renderizarTabelaItens() {
    const tbody = document.getElementById('tabelaProdutosEntradaBody');
    const trVazia = document.getElementById('trEntradaVazia');
    const badge = document.getElementById('badge_qtd_itens');
    const totalEl = document.getElementById('totalGeralEntrada');

    // Limpar linhas de produto (manter a linha vazia)
    tbody.querySelectorAll('tr.linha-item').forEach(tr => tr.remove());

    badge.textContent = `${itensDaEntrada.length} ${itensDaEntrada.length === 1 ? 'item' : 'itens'}`;

    if (itensDaEntrada.length === 0) {
        trVazia.classList.remove('d-none');
        totalEl.textContent = 'R$ 0,00';
        return;
    }

    trVazia.classList.add('d-none');

    let totalGeral = 0;

    itensDaEntrada.forEach((item, idx) => {
        totalGeral += item.valorTotal;
        const tr = document.createElement('tr');
        tr.className = 'linha-item';
        tr.innerHTML = `
            <td class="text-muted small">${item.idProduto}</td>
            <td class="fw-semibold text-start">${item.nome}</td>
            <td class="text-center">${item.quantidade % 1 === 0 ? item.quantidade.toFixed(0) : item.quantidade}</td>
            <td class="text-end">R$ ${formatBRL(item.valorUnitario)}</td>
            <td class="text-end text-muted small">R$ ${formatBRL(item.icmsSt)}</td>
            <td class="text-end text-muted small">R$ ${formatBRL(item.ipi)}</td>
            <td class="text-end text-muted small">R$ ${formatBRL(item.desconto)}</td>
            <td class="text-end text-muted small">R$ ${formatBRL(item.rateio)}</td>
            <td class="text-end fw-bold text-success">R$ ${item.valorCusto.toFixed(4)}</td>
            <td class="text-end fw-bold text-primary">R$ ${formatBRL(item.valorTotal)}</td>
            <td class="text-center">
                <button type="button" class="btn btn-sm btn-outline-danger btn-remover-item" data-idx="${idx}" title="Remover item">
                    <i class="bi bi-trash"></i>
                </button>
            </td>
        `;
        tbody.appendChild(tr);
    });

    totalEl.textContent = `R$ ${formatBRL(totalGeral)}`;

    // Delegação de eventos para remover
    tbody.querySelectorAll('.btn-remover-item').forEach(btn => {
        btn.addEventListener('click', () => {
            const i = parseInt(btn.dataset.idx, 10);
            itensDaEntrada.splice(i, 1);
            recalcularRateios();
        });
    });
}

// ─── Busca de Produto ────────────────────────────────────────────────────────

async function buscarProduto() {
    const busca = document.getElementById('entrada_busca_produto').value.trim();
    if (!busca) {
        showError('Informe o nome ou código de barras do produto para buscar.');
        return;
    }

    const token = getToken();
    try {
        const res = await fetch(`/api/produtos?busca=${encodeURIComponent(busca)}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!res.ok) {
            const data = await res.json();
            showError(data.erro || 'Erro ao buscar produto.');
            return;
        }

        const produtos = await res.json();

        if (!produtos || produtos.length === 0) {
            showError('Nenhum produto encontrado com esse nome ou código.');
            esconderCardProduto();
            return;
        }

        // Usa o primeiro resultado (pode ser expandido com autocomplete futuramente)
        const p = produtos[0];
        produtoSelecionado = p;
        preencherCardProduto(p);

    } catch (err) {
        console.error(err);
        showError('Erro interno ao buscar produto.');
    }
}

function preencherCardProduto(p) {
    document.getElementById('ps_id').textContent = p.id || '—';
    document.getElementById('ps_codigo_barras').textContent = p.codigo_barras || '—';
    document.getElementById('ps_codigo_interno').textContent = p.codigo_interno_loja || '—';
    document.getElementById('ps_nome').textContent = p.nome || '—';

    // Limpar campos de valor
    ['ps_quantidade', 'ps_valor_unitario', 'ps_icms_st', 'ps_ipi', 'ps_desconto'].forEach(id => {
        const el = document.getElementById(id);
        if (id === 'ps_icms_st' || id === 'ps_ipi' || id === 'ps_desconto') el.value = '0';
        else el.value = '';
    });
    document.getElementById('ps_rateio').value = '0';
    document.getElementById('ps_valor_custo').value = '';
    document.getElementById('ps_valor_total').value = '';

    document.getElementById('cardProdutoSelecionado').classList.remove('d-none');
    document.getElementById('ps_quantidade').focus();
}

function esconderCardProduto() {
    document.getElementById('cardProdutoSelecionado').classList.add('d-none');
    produtoSelecionado = null;
}

// ─── Adicionar item à entrada ────────────────────────────────────────────────

function adicionarItemEntrada() {
    if (!produtoSelecionado) {
        showError('Nenhum produto selecionado. Busque um produto antes de adicionar.');
        return;
    }

    const qtd = parseFloat(document.getElementById('ps_quantidade').value);
    const vlrUnit = parseFloat(document.getElementById('ps_valor_unitario').value);

    if (!qtd || qtd <= 0) {
        showError('Informe uma quantidade válida (maior que zero).');
        document.getElementById('ps_quantidade').focus();
        return;
    }

    if (!vlrUnit || vlrUnit < 0) {
        showError('Informe um valor unitário válido.');
        document.getElementById('ps_valor_unitario').focus();
        return;
    }

    const icmsSt = parseFloat(document.getElementById('ps_icms_st').value) || 0;
    const ipi = parseFloat(document.getElementById('ps_ipi').value) || 0;
    const desconto = parseFloat(document.getElementById('ps_desconto').value) || 0;

    const novoItem = {
        idProduto: produtoSelecionado.id,
        nome: produtoSelecionado.nome,
        quantidade: qtd,
        valorUnitario: vlrUnit,
        icmsSt,
        ipi,
        desconto,
        rateio: 0,
        valorCusto: 0,
        valorTotal: 0,
    };

    itensDaEntrada.push(novoItem);
    recalcularRateios(); // Recalcula rateios de todos os itens e re-renderiza

    // Limpar card e campo de busca
    esconderCardProduto();
    document.getElementById('entrada_busca_produto').value = '';
}

// ─── Carrega Fornecedores no Select ─────────────────────────────────────────

async function carregarFornecedoresSelect(selectId) {
    const token = getToken();
    const select = document.getElementById(selectId);
    if (!select) return;

    try {
        const res = await fetch('/api/fornecedores', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!res.ok) return;

        const fornecedores = await res.json();
        select.innerHTML = '<option value="" disabled selected>Selecione o fornecedor...</option>';
        fornecedores.forEach(f => {
            const opt = document.createElement('option');
            opt.value = f.id;
            opt.textContent = f.razao_social;
            select.appendChild(opt);
        });
    } catch (err) {
        console.error('Erro ao carregar fornecedores:', err);
    }
}

// ─── Resetar Formulário ──────────────────────────────────────────────────────

function resetarFormEntrada() {
    itensDaEntrada = [];
    produtoSelecionado = null;
    document.getElementById('entrada_id_fornecedor').value = '';
    document.getElementById('entrada_despesa_adicional').value = '0';
    document.getElementById('entrada_busca_produto').value = '';
    esconderCardProduto();
    renderizarTabelaItens();
}

// ─── Setup Principal ─────────────────────────────────────────────────────────

export function setupBotaoNovaEntrada() {
    const btn = document.getElementById('btnNovaEntrada');
    if (!btn) return;

    btn.addEventListener('click', () => {
        const collapse = bootstrap.Collapse.getOrCreateInstance(
            document.getElementById('collapseFormEntrada')
        );
        collapse.show();
        carregarFornecedoresSelect('entrada_id_fornecedor');
        btn.scrollIntoView({ behavior: 'smooth' });
    });
}

export function setupFormEntrada() {
    // Fechar painel
    document.getElementById('btnFecharFormEntrada')?.addEventListener('click', () => {
        bootstrap.Collapse.getOrCreateInstance(
            document.getElementById('collapseFormEntrada')
        ).hide();
        resetarFormEntrada();
    });

    document.getElementById('btnCancelarEntrada')?.addEventListener('click', () => {
        bootstrap.Collapse.getOrCreateInstance(
            document.getElementById('collapseFormEntrada')
        ).hide();
        resetarFormEntrada();
    });

    // Busca de produto
    document.getElementById('btnBuscarProdutoEntrada')?.addEventListener('click', buscarProduto);
    document.getElementById('entrada_busca_produto')?.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') { e.preventDefault(); buscarProduto(); }
    });

    // Recalcular ao mudar campos do card
    ['ps_quantidade', 'ps_valor_unitario', 'ps_icms_st', 'ps_ipi', 'ps_desconto'].forEach(id => {
        document.getElementById(id)?.addEventListener('input', recalcularCamposCard);
    });

    // Recalcular rateios ao mudar despesa adicional
    document.getElementById('entrada_despesa_adicional')?.addEventListener('input', () => {
        recalcularRateios();
        recalcularCamposCard();
    });

    // Adicionar produto
    document.getElementById('btnAdicionarProdutoEntrada')?.addEventListener('click', adicionarItemEntrada);

    // Salvar (em aberto) — apenas UI por enquanto, backend será implementado
    document.getElementById('btnSalvarEntrada')?.addEventListener('click', () => {
        if (itensDaEntrada.length === 0) {
            showError('Adicione ao menos um produto antes de salvar.');
            return;
        }
        const payload = montarPayload('ABERTA');
        console.log('Salvar entrada (ABERTA):', payload);
        // TODO: Chamar API quando backend estiver pronto
        showError('Funcionalidade de salvar será implementada com o backend.');
    });

    // Concluir
    document.getElementById('btnConcluirEntrada')?.addEventListener('click', () => {
        if (itensDaEntrada.length === 0) {
            showError('Adicione ao menos um produto antes de concluir.');
            return;
        }
        const payload = montarPayload('CONCLUIDA');
        console.log('Concluir entrada:', payload);
        // TODO: Chamar API quando backend estiver pronto
        showError('Funcionalidade de concluir será implementada com o backend.');
    });
}

// ─── Montar payload para envio ────────────────────────────────────────────────

export function montarPayload(status) {
    return {
        id_fornecedor: parseInt(document.getElementById('entrada_id_fornecedor').value) || null,
        despesa_adicional: getDespesaAdicional(),
        status,
        produtos: itensDaEntrada.map(item => ({
            id_produto: item.idProduto,
            quantidade: item.quantidade,
            valor_unitario: item.valorUnitario,
            valor_icms_st: item.icmsSt,
            valor_ipi: item.ipi,
            valor_desconto: item.desconto,
            rateio_despesa_adicional: item.rateio,
            valor_custo: item.valorCusto,
            valor_total: item.valorTotal,
        })),
    };
}

/** Expõe o total geral para o módulo de financeiro */
export function getTotalGeralEntrada() {
    return itensDaEntrada.reduce((acc, it) => acc + it.valorTotal, 0);
}

/** Expõe o id_fornecedor selecionado */
export function getIdFornecedorSelecionado() {
    return document.getElementById('entrada_id_fornecedor')?.value || '';
}
