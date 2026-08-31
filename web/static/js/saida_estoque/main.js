import { getToken, checkAuth } from '../utils/auth.js';
import { showError, showSuccess } from '../utils/showError.js';

let itens = [];
let produtosDoEstoque = [];
let saidaEmEdicaoId = null;

const api = (url, opt = {}) => fetch(url, { ...opt, headers: { Authorization: `Bearer ${getToken()}`, 'Content-Type': 'application/json', ...(opt.headers || {}) } });
const brl = valor => Number(valor || 0).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
const numero = valor => Number(valor || 0);
const sufixoQuery = parametros => parametros.toString() ? `?${parametros}` : '';
const formatarStatus = status => ({ ABERTO: 'Em aberto', CONCLUIDA: 'Concluída', CANCELADA: 'Cancelada' })[status] || status;

async function respostaJson(response) {
    const corpo = await response.json();
    if (!response.ok) throw new Error(corpo.erro || 'Não foi possível concluir a operação.');
    return corpo;
}

async function carregarEstoques() {
    const estoques = await respostaJson(await api('/api/estoques'));
    saidaEstoque.innerHTML = (estoques || []).map(estoque => `<option value="${estoque.id}">${estoque.nome}</option>`).join('');
    await carregarProdutosDoEstoque();
}

async function carregarProdutosDoEstoque() {
    produtosDoEstoque = saidaEstoque.value ? await respostaJson(await api(`/api/estoques/${saidaEstoque.value}/produtos`)) : [];
    saidaProduto.innerHTML = produtosDoEstoque.length
        ? produtosDoEstoque.map(item => `<option value="${item.id_produto}" data-saldo="${item.quantidade}" data-preco="${item.produto?.preco_venda || 0}" data-custo="${item.produto?.preco_custo || 0}">${item.produto?.nome || 'Produto'} (${item.quantidade} disponível)</option>`).join('')
        : '<option value="">Nenhum produto vinculado a este estoque</option>';
    saidaProduto.disabled = !produtosDoEstoque.length;
    atualizarSaldoProduto();
}

function atualizarSaldoProduto() {
    const produto = saidaProduto.selectedOptions[0];
    saldoProduto.textContent = produto?.value ? `Saldo disponível: ${produto.dataset.saldo}` : 'Nenhum produto disponível neste estoque.';
}

async function carregarSaidas() {
    const filtros = new URLSearchParams();
    if (filtroStatus.value) filtros.set('status', filtroStatus.value);
    if (filtroData.value) filtros.set('data', filtroData.value);
    const saidas = await respostaJson(await api(`/api/saidas-estoque${sufixoQuery(filtros)}`));
    listaSaidas.innerHTML = (saidas || []).map(saida => `<tr><td>#${saida.id}</td><td>${saida.estoque}</td><td>${brl(saida.valor_total)}</td><td>${formatarStatus(saida.status)}</td><td>${new Date(saida.criado_em).toLocaleDateString('pt-BR')}</td><td class="text-end"><button class="btn btn-sm btn-outline-primary ver" data-id="${saida.id}">Ver</button>${saida.status === 'ABERTO' ? ` <button class="btn btn-sm btn-outline-secondary editar" data-id="${saida.id}">Editar</button> <button class="btn btn-sm btn-success aprovar" data-id="${saida.id}">Aprovar</button>` : ''}${saida.status !== 'CANCELADA' ? ` <button class="btn btn-sm btn-outline-danger cancelar" data-id="${saida.id}">Cancelar</button>` : ''}</td></tr>`).join('') || '<tr><td colspan="6" class="text-center">Nenhuma saída encontrada.</td></tr>';
    document.querySelectorAll('.aprovar').forEach(botao => botao.onclick = () => aprovar(botao.dataset.id));
    document.querySelectorAll('.cancelar').forEach(botao => botao.onclick = () => cancelar(botao.dataset.id));
    document.querySelectorAll('.editar').forEach(botao => botao.onclick = () => editar(botao.dataset.id));
    document.querySelectorAll('.ver').forEach(botao => botao.onclick = () => ver(botao.dataset.id));
}

function renderizarItens() {
    itensSaida.innerHTML = itens.map((item, indice) => `<tr><td class="fw-semibold">${item.nome}</td><td class="text-center">${item.saldo}</td><td class="text-center">${item.quantidade}</td><td class="text-end">${brl(item.valor_unitario)}</td><td class="text-end fw-bold">${brl(item.valor_total)}</td><td class="text-end"><button class="btn btn-sm btn-outline-danger remover" data-indice="${indice}"><i class="bi bi-trash"></i></button></td></tr>`).join('') || '<tr><td colspan="6" class="text-center text-muted py-4">Nenhum item adicionado.</td></tr>';
    totalSaida.textContent = brl(itens.reduce((total, item) => total + item.valor_total, 0));
    document.querySelectorAll('.remover').forEach(botao => botao.onclick = () => { itens.splice(botao.dataset.indice, 1); renderizarItens(); });
}

function adicionarItem() {
    const opcao = saidaProduto.selectedOptions[0];
    const quantidade = numero(saidaQuantidade.value);
    if (!opcao?.value || quantidade <= 0) return showError('Informe um produto e uma quantidade válida.');
    if (itens.some(item => item.id_produto === numero(opcao.value))) return showError('Este produto já foi adicionado. Altere a quantidade ou remova a linha existente.');
    if (quantidade > numero(opcao.dataset.saldo)) return showError(`Quantidade maior que o saldo disponível (${opcao.dataset.saldo}).`);
    itens.push({ id_produto: numero(opcao.value), nome: opcao.text.split(' (')[0], saldo: numero(opcao.dataset.saldo), quantidade, valor_unitario: numero(opcao.dataset.preco), valor_custo: numero(opcao.dataset.custo), valor_total: numero(opcao.dataset.preco) * quantidade });
    saidaQuantidade.value = '';
    renderizarItens();
}

async function salvar(concluir = false) {
    if (!itens.length) return showError('Adicione ao menos um item.');
    const saida = { id_estoque: numero(saidaEstoque.value), status: concluir ? 'CONCLUIDA' : 'ABERTO', produtos: itens.map(({ saldo, nome, ...item }) => item) };
    try {
        const url = saidaEmEdicaoId ? `/api/saidas-estoque/${saidaEmEdicaoId}` : '/api/saidas-estoque';
        await respostaJson(await api(url, { method: saidaEmEdicaoId ? 'PUT' : 'POST', body: JSON.stringify(saida) }));
        showSuccess(saidaEmEdicaoId ? 'Saída atualizada.' : concluir ? 'Saída salva e concluída.' : 'Saída salva em aberto.');
        fecharFormulario();
        await carregarSaidas();
    } catch (erro) { showError(erro.message); }
}

async function aprovar(id) {
    if (!confirm('Aprovar a saída e baixar o saldo?')) return;
    try { showSuccess((await respostaJson(await api(`/api/saidas-estoque/${id}/aprovar`, { method: 'POST' }))).mensagem); await carregarSaidas(); } catch (erro) { showError(erro.message); }
}

async function cancelar(id) {
    if (!confirm('Cancelar esta saída? Uma saída concluída terá seu saldo estornado.')) return;
    try { showSuccess((await respostaJson(await api(`/api/saidas-estoque/${id}/cancelar`, { method: 'POST' }))).mensagem); await carregarSaidas(); } catch (erro) { showError(erro.message); }
}

async function obter(id) { return respostaJson(await api(`/api/saidas-estoque/${id}`)); }
async function ver(id) {
    try { const saida = await obter(id); alert(`Saída #${saida.id}\nStatus: ${formatarStatus(saida.status)}\nItens:\n${saida.produtos.map(produto => `${produto.nome_produto}: ${produto.quantidade}`).join('\n')}`); } catch (erro) { showError(erro.message); }
}

async function editar(id) {
    try {
        const saida = await obter(id);
        saidaEmEdicaoId = id;
        formSaida.classList.remove('d-none');
        saidaEstoque.value = String(saida.id_estoque);
        await carregarProdutosDoEstoque();
        itens = saida.produtos.map(produto => {
            const estoque = produtosDoEstoque.find(item => numero(item.id_produto) === numero(produto.id_produto));
            return { ...produto, nome: produto.nome_produto, saldo: numero(estoque?.quantidade) || 0 };
        });
        tituloSaida.innerHTML = '<i class="bi bi-pencil-square me-2"></i>Editar Saída';
        textoSalvarSaida.textContent = 'Salvar alterações';
		btnConcluirSaida.classList.add('d-none');
        renderizarItens();
    } catch (erro) { showError(erro.message); }
}

function novaSaida() { saidaEmEdicaoId = null; itens = []; tituloSaida.innerHTML = '<i class="bi bi-box-arrow-up-right me-2"></i>Nova Saída'; textoSalvarSaida.textContent = 'Salvar em aberto'; btnConcluirSaida.classList.remove('d-none'); formSaida.classList.remove('d-none'); renderizarItens(); }
function fecharFormulario() { itens = []; saidaEmEdicaoId = null; btnConcluirSaida.classList.remove('d-none'); formSaida.classList.add('d-none'); renderizarItens(); }

document.addEventListener('DOMContentLoaded', async () => {
    if (!checkAuth()) return;
    try { await carregarEstoques(); await carregarSaidas(); } catch (erro) { showError(erro.message); }
    btnNovaSaida.onclick = novaSaida;
    btnCancelarSaida.onclick = fecharFormulario;
    btnAdicionarItem.onclick = adicionarItem;
    btnSalvarSaida.onclick = () => salvar(false);
    btnConcluirSaida.onclick = () => salvar(true);
    saidaEstoque.onchange = () => carregarProdutosDoEstoque().catch(erro => showError(erro.message));
    saidaProduto.onchange = atualizarSaldoProduto;
    filtroStatus.onchange = () => carregarSaidas().catch(erro => showError(erro.message));
    filtroData.onchange = () => carregarSaidas().catch(erro => showError(erro.message));
    btnLimparFiltros.onclick = () => { filtroStatus.value = ''; filtroData.value = ''; carregarSaidas().catch(erro => showError(erro.message)); };
});
