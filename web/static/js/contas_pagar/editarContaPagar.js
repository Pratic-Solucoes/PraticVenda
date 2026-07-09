import { getToken } from '../utils/auth.js';
import { state } from './state.js';
import { carregarContasPagar } from './listarContasPagar.js';

export function setupEditarContaPagar() {
    window.abrirModalEdicao = function (id) {
        const debito = state.contasPagarCarregadas.find(d => d.id === id);
        if (!debito) return;

        document.getElementById('edit_id_conta_pagar').value = debito.id;
        document.getElementById('edit_id_fornecedor').value = debito.id_fornecedor;
        document.getElementById('edit_id_categoria').value = debito.id_categoria || '';
        document.getElementById('edit_descricao').value = debito.descricao;
        document.getElementById('edit_nr_documento').value = debito.nr_documento || '';
        document.getElementById('edit_nr_nota_fiscal').value = debito.nr_nota_fiscal || '';
        document.getElementById('edit_valor').value = debito.valor_original;
        document.getElementById('edit_dt_entrada').value = debito.dt_entrada.substring(0, 10);
        document.getElementById('edit_dt_vencimento').value = debito.dt_vencimento.substring(0, 10);
        document.getElementById('edit_nr_parcela').value = debito.nr_parcela;
        document.getElementById('edit_nr_total_parcelas').value = debito.nr_total_parcelas;

        const modalEl = document.getElementById('modalEditarContaPagar');
        const modal = new bootstrap.Modal(modalEl);
        modal.show();
    };

    const formEditarContaPagar = document.getElementById('formEditarContaPagar');
    if (formEditarContaPagar) {
        formEditarContaPagar.addEventListener('submit', async (e) => {
            e.preventDefault();

            const token = getToken();
            const id = document.getElementById('edit_id_conta_pagar').value;
            const id_fornecedor = parseInt(document.getElementById('edit_id_fornecedor').value, 10);
            let id_categoria = parseInt(document.getElementById('edit_id_categoria').value, 10);
            if (isNaN(id_categoria)) id_categoria = null;

            const payload = {
                id_fornecedor: id_fornecedor,
                descricao: document.getElementById('edit_descricao').value,
                nr_documento: document.getElementById('edit_nr_documento').value,
                nr_nota_fiscal: document.getElementById('edit_nr_nota_fiscal').value,
                valor_original: parseFloat(document.getElementById('edit_valor').value),
                dt_entrada: document.getElementById('edit_dt_entrada').value,
                dt_vencimento: document.getElementById('edit_dt_vencimento').value,
                nr_parcela: parseInt(document.getElementById('edit_nr_parcela').value, 10),
                nr_total_parcelas: parseInt(document.getElementById('edit_nr_total_parcelas').value, 10)
            };

            if (id_categoria) {
                payload.id_categoria = id_categoria;
            }

            try {
                const res = await fetch(`/api/contas-pagar/${id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(payload)
                });

                if (!res.ok) {
                    const data = await res.json();
                    alert(data.erro || "Erro ao editar débito.");
                    return;
                }

                alert("Débito atualizado com sucesso!");
                const modalEl = document.getElementById('modalEditarContaPagar');
                const modal = bootstrap.Modal.getInstance(modalEl);
                if (modal) modal.hide();

                carregarContasPagar();
            } catch (err) {
                console.error(err);
                alert("Erro interno ao comunicar com o servidor.");
            }
        });
    }
}
