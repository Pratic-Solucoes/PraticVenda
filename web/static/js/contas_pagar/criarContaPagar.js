import { getToken } from '../utils/auth.js';
import { carregarContasPagar } from './listarContasPagar.js';

export function setupCriarContaPagar() {
    const formContaPagar = document.getElementById('formContaPagarAvulso');
    if (!formContaPagar) return;

    formContaPagar.addEventListener('submit', async (e) => {
        e.preventDefault();

        const token = getToken();
        const selectFornecedor = document.getElementById('id_fornecedor');
        const selectCategoria = document.getElementById('id_categoria');

        const id_fornecedor = parseInt(selectFornecedor.value, 10);
        let id_categoria = parseInt(selectCategoria.value, 10);
        if (isNaN(id_categoria)) id_categoria = null;

        const payload = {
            id_fornecedor: id_fornecedor,
            descricao: document.getElementById('descricao').value,
            nr_documento: document.getElementById('nr_documento').value,
            nr_nota_fiscal: document.getElementById('nr_nota_fiscal').value,
            valor_original: parseFloat(document.getElementById('valor').value),
            dt_entrada: document.getElementById('dt_entrada').value,
            dt_vencimento: document.getElementById('dt_vencimento').value,
            nr_parcela: parseInt(document.getElementById('nr_parcela').value, 10),
            nr_total_parcelas: parseInt(document.getElementById('nr_total_parcelas').value, 10)
        };

        if (id_categoria) {
            payload.id_categoria = id_categoria;
        }

        try {
            const res = await fetch('/api/contas-pagar', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(payload)
            });

            if (!res.ok) {
                const data = await res.json();
                alert(data.erro || "Erro ao cadastrar débito.");
                return;
            }

            alert("Pagamento avulso cadastrado com sucesso!");
            const modalEl = document.getElementById('modalContaPagarAvulso');
            const modal = bootstrap.Modal.getInstance(modalEl);
            if (modal) modal.hide();

            formContaPagar.reset();
            const tabelaDebitos = document.getElementById('tabela_debitos_body');
            if (tabelaDebitos) {
                carregarContasPagar();
            } else {
                window.location.reload();
            }
        } catch (err) {
            console.error(err);
            alert("Erro interno ao comunicar com o servidor.");
        }
    });
}
