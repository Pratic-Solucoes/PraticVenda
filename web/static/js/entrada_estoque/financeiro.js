import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';
import { getTotalGeralEntrada, getIdFornecedorSelecionado } from './formEntrada.js';

let parcelas = [];

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatDataBR(date) {
    const d = date instanceof Date ? date : new Date();
    return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`;
}

function formatDateInput(date) {
    const d = date instanceof Date ? date : new Date();
    return d.toISOString().slice(0, 10);
}

function moeda(valor) { return Number(valor || 0).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }); }

function renderizarParcelas() {
    const body = document.getElementById('finParcelasBody');
    if (!body) return;
    body.innerHTML = parcelas.map((parcela, indice) => `<tr><td>${indice + 1}</td><td><input class="form-control form-control-sm fin-parcela-data" data-indice="${indice}" type="date" value="${parcela.dt_vencimento}"></td><td><input class="form-control form-control-sm text-end fin-parcela-valor" data-indice="${indice}" type="number" min="0.01" step="0.01" value="${parcela.valor.toFixed(2)}"></td><td><button class="btn btn-sm btn-outline-danger fin-remover-parcela" data-indice="${indice}" type="button" ${parcelas.length === 1 ? 'disabled' : ''}><i class="bi bi-trash"></i></button></td></tr>`).join('');
    document.querySelectorAll('.fin-parcela-data').forEach(input => input.onchange = () => { parcelas[input.dataset.indice].dt_vencimento = input.value; });
    document.querySelectorAll('.fin-parcela-valor').forEach(input => input.oninput = () => { parcelas[input.dataset.indice].valor = Number(input.value || 0); atualizarTotalParcelas(); });
    document.querySelectorAll('.fin-remover-parcela').forEach(botao => botao.onclick = () => { parcelas.splice(botao.dataset.indice, 1); renderizarParcelas(); atualizarTotalParcelas(); });
    atualizarTotalParcelas();
}

function atualizarTotalParcelas() {
    const totalEntrada = getTotalGeralEntrada();
    const totalParcelas = parcelas.reduce((total, parcela) => total + Number(parcela.valor || 0), 0);
    document.getElementById('finTotalParcelas').textContent = moeda(totalParcelas);
    const diferenca = totalEntrada - totalParcelas;
    const aviso = document.getElementById('finDiferencaParcelas');
    aviso.textContent = Math.abs(diferenca) < 0.005 ? 'Valores conferem com o total da entrada.' : `Falta distribuir ${moeda(diferenca)} nas parcelas.`;
    aviso.className = Math.abs(diferenca) < 0.005 ? 'small text-success' : 'small text-danger';
}

async function carregarFornecedoresModal() {
    const token = getToken();
    const select = document.getElementById('fin_id_fornecedor');
    if (!select) return;

    try {
        const res = await fetch('/api/fornecedores', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!res.ok) return;
        const fornecedores = await res.json();
        select.innerHTML = '<option value="" disabled>Selecione o fornecedor...</option>';
        fornecedores.forEach(f => {
            const opt = document.createElement('option');
            opt.value = f.id;
            opt.textContent = f.razao_social;
            select.appendChild(opt);
        });
    } catch (err) {
        console.error('Erro ao carregar fornecedores no modal financeiro:', err);
    }
}

async function carregarCategoriasModal() {
    const token = getToken();
    const select = document.getElementById('fin_id_categoria');
    if (!select) return;

    try {
        const res = await fetch('/api/contas-pagar/categorias', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!res.ok) return;
        const cats = await res.json();
        select.innerHTML = '<option value="">Sem Categoria (Opcional)</option>';
        cats.forEach(c => {
            const opt = document.createElement('option');
            opt.value = c.id;
            opt.textContent = c.nome;
            select.appendChild(opt);
        });
    } catch (err) {
        console.error('Erro ao carregar categorias no modal financeiro:', err);
    }
}

/** Pré-preenche o modal com os dados da entrada em andamento */
function preencherModalFinanceiro() {
    const hoje = new Date();

    // Descrição automática
    const descEl = document.getElementById('fin_descricao');
    if (descEl) descEl.value = `Entrada Estoque - ${formatDataBR(hoje)}`;

    // Valor total da entrada
    const valorEl = document.getElementById('fin_valor');
    if (valorEl) valorEl.value = getTotalGeralEntrada().toFixed(2);

    parcelas = [{ dt_vencimento: formatDateInput(hoje), valor: getTotalGeralEntrada() }];
    renderizarParcelas();

    // Data de entrada = hoje
    const dtEntradaEl = document.getElementById('fin_dt_entrada');
    if (dtEntradaEl) dtEntradaEl.value = formatDateInput(hoje);

    // Pré-selecionar fornecedor se selecionado no form principal
    const idFornecedor = getIdFornecedorSelecionado();
    if (idFornecedor) {
        const selectFin = document.getElementById('fin_id_fornecedor');
        if (selectFin) {
            // Aguarda carregamento antes de selecionar
            const tentarSelecionar = setInterval(() => {
                if (selectFin.querySelector(`option[value="${idFornecedor}"]`)) {
                    selectFin.value = idFornecedor;
                    clearInterval(tentarSelecionar);
                }
            }, 100);
            setTimeout(() => clearInterval(tentarSelecionar), 3000);
        }
    }
}

/** Envia o lançamento financeiro para a API de contas a pagar */
async function enviarFinanceiro(e) {
    e.preventDefault();
    const token = getToken();

    const idFornecedor = parseInt(document.getElementById('fin_id_fornecedor').value);
    if (!idFornecedor) {
        showError('Selecione um fornecedor para gerar o financeiro.');
        return;
    }

    const idCategoria = parseInt(document.getElementById('fin_id_categoria').value);
    if (!idCategoria) {
        showError('Selecione uma categoria para o lançamento.');
        return;
    }
	const totalEntrada = getTotalGeralEntrada();
	const totalParcelas = parcelas.reduce((total, parcela) => total + Number(parcela.valor || 0), 0);
	if (!parcelas.length || parcelas.some(parcela => !parcela.dt_vencimento || parcela.valor <= 0) || Math.abs(totalEntrada - totalParcelas) >= 0.005) {
		showError('Informe vencimento e valor de cada parcela. A soma deve ser igual ao total da entrada.');
		return;
	}

    const dadosComuns = {
        id_fornecedor:      idFornecedor,
        id_categoria:       idCategoria,
        descricao:          document.getElementById('fin_descricao').value.trim(),
        nr_documento:       document.getElementById('fin_nr_documento').value.trim() || null,
        nr_nota_fiscal:     document.getElementById('fin_nr_nota_fiscal').value.trim() || null,
        dt_entrada:         document.getElementById('fin_dt_entrada').value,
    };

    try {
        for (let indice = 0; indice < parcelas.length; indice += 1) {
            const res = await fetch('/api/contas-pagar', { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify({ ...dadosComuns, valor_original: parcelas[indice].valor, dt_vencimento: parcelas[indice].dt_vencimento, nr_parcela: indice + 1, nr_total_parcelas: parcelas.length }) });
            if (!res.ok) { const data = await res.json(); throw new Error(data.erro || 'Erro ao lançar conta a pagar.'); }
        }

        // Fechar modal
        const modalEl = document.getElementById('modalFinanceiroEntrada');
        bootstrap.Modal.getInstance(modalEl)?.hide();

        alert('Lançamento financeiro realizado com sucesso!');

    } catch (err) {
        console.error(err);
        showError(err.message || 'Erro interno ao comunicar com o servidor.');
    }
}

// ─── Setup ───────────────────────────────────────────────────────────────────

export function setupFinanceiro() {
	document.getElementById('btnAdicionarParcela')?.addEventListener('click', () => { parcelas.push({ dt_vencimento: formatDateInput(new Date()), valor: 0 }); renderizarParcelas(); });
    // Botão "Gerar Financeiro" abre o modal e pré-preenche
    document.getElementById('btnGerarFinanceiro')?.addEventListener('click', () => {
        const totalEntrada = getTotalGeralEntrada();
        if (totalEntrada <= 0) {
            showError('Adicione produtos à entrada antes de gerar o financeiro.');
            return;
        }

        // Carregar selects do modal
        carregarFornecedoresModal();
        carregarCategoriasModal();
        preencherModalFinanceiro();

        const modalEl = document.getElementById('modalFinanceiroEntrada');
        bootstrap.Modal.getOrCreateInstance(modalEl).show();
    });

    // Submit do form do modal
    document.getElementById('formFinanceiroEntrada')?.addEventListener('submit', enviarFinanceiro);
}
