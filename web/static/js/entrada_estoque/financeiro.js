import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';
import { getTotalGeralEntrada, getIdFornecedorSelecionado } from './formEntrada.js';

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatDataBR(date) {
    const d = date instanceof Date ? date : new Date();
    return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`;
}

function formatDateInput(date) {
    const d = date instanceof Date ? date : new Date();
    return d.toISOString().slice(0, 10);
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

    const valor = parseFloat(document.getElementById('fin_valor').value);
    if (!valor || valor <= 0) {
        showError('Informe um valor válido para o lançamento.');
        return;
    }

    const payload = {
        id_fornecedor:      idFornecedor,
        id_categoria:       parseInt(document.getElementById('fin_id_categoria').value) || null,
        descricao:          document.getElementById('fin_descricao').value.trim(),
        nr_documento:       document.getElementById('fin_nr_documento').value.trim() || null,
        nr_nota_fiscal:     document.getElementById('fin_nr_nota_fiscal').value.trim() || null,
        valor,
        dt_entrada:         document.getElementById('fin_dt_entrada').value,
        dt_vencimento:      document.getElementById('fin_dt_vencimento').value,
        nr_parcela:         parseInt(document.getElementById('fin_nr_parcela').value) || 1,
        nr_total_parcelas:  parseInt(document.getElementById('fin_nr_total_parcelas').value) || 1,
    };

    try {
        const res = await fetch('/api/contas-pagar', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });

        if (!res.ok) {
            const data = await res.json();
            showError(data.erro || 'Erro ao lançar conta a pagar.');
            return;
        }

        // Fechar modal
        const modalEl = document.getElementById('modalFinanceiroEntrada');
        bootstrap.Modal.getInstance(modalEl)?.hide();

        alert('Lançamento financeiro realizado com sucesso!');

    } catch (err) {
        console.error(err);
        showError('Erro interno ao comunicar com o servidor.');
    }
}

// ─── Setup ───────────────────────────────────────────────────────────────────

export function setupFinanceiro() {
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
