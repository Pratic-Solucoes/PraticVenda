import { state } from './state.js';

export function setupVisualizarDebito() {
    window.abrirModalVisualizacao = function (id) {
        const debito = state.debitosCarregados.find(d => d.id === id);
        if (!debito) return;

        document.getElementById('view_fornecedor').textContent = debito.fornecedor ? debito.fornecedor.razao_social : '-';
        document.getElementById('view_categoria').textContent = debito.categoria ? debito.categoria.nome : 'Sem Categoria';
        document.getElementById('view_descricao').textContent = debito.descricao;
        document.getElementById('view_nr_documento').textContent = debito.nr_documento || '-';
        document.getElementById('view_nr_nota_fiscal').textContent = debito.nr_nota_fiscal || '-';
        document.getElementById('view_valor').textContent = new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(debito.valor);
        document.getElementById('view_dt_entrada').textContent = debito.dt_entrada ? debito.dt_entrada.substring(0, 10).split('-').reverse().join('/') : '-';
        document.getElementById('view_dt_vencimento').textContent = debito.dt_vencimento ? debito.dt_vencimento.substring(0, 10).split('-').reverse().join('/') : '-';
        document.getElementById('view_parcela').textContent = `${debito.nr_parcela} / ${debito.nr_total_parcelas}`;
        document.getElementById('view_status').textContent = debito.status;

        const modalEl = document.getElementById('modalVisualizarDebito');
        const modal = new bootstrap.Modal(modalEl);
        modal.show();
    };
}
