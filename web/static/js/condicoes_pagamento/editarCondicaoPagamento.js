import { getToken } from '../utils/auth.js';
import { carregarCondicoesPagamento } from './listarCondicoesPagamento.js';
import { showError } from '../utils/showError.js';

export function setupEditarCondicaoPagamento() {

    // Abre o formulário inline preenchido com os dados da condição buscada por ID
    window.abrirEditarCondicaoPagamento = async function(id) {
        const token = getToken();
        try {
            const res = await fetch(`/api/condicoes-pagamento/${id}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (!res.ok) {
                showError("Erro ao buscar condição de pagamento.");
                return;
            }

            const cp = await res.json();

            document.getElementById('edit_cp_id').value = cp.id;
            document.getElementById('edit_cp_descricao').value = cp.descricao || '';
            document.getElementById('edit_cp_qtd_parcelas').value = cp.qtd_parcelas || '';
            document.getElementById('edit_cp_dias_venc').value = cp.dias_primeiro_venc || '';
            document.getElementById('edit_cp_intervalo').value = cp.intervalo_parcelas || '';

            // Reseta todos checkboxes
            document.querySelectorAll('.fp-checkbox').forEach(cb => cb.checked = false);
            // Marca os checkboxes que vieram na condição
            if (cp.formas_pagamento && cp.formas_pagamento.length > 0) {
                cp.formas_pagamento.forEach(fp_id => {
                    const cb = document.getElementById(`fp_check_${fp_id}`);
                    if (cb) cb.checked = true;
                });
            }

            // Atualiza título do painel
            const titulo = document.getElementById('formCondicaoPagamentoTitulo');
            if (titulo) {
                titulo.innerHTML = `<i class="bi bi-pencil-square me-2"></i> Editando Condição de Pagamento #${id}`;
            }

            if (window.abrirFormularioInlineCondicaoPagamento) {
                window.abrirFormularioInlineCondicaoPagamento();
            }
        } catch (err) {
            console.error(err);
            showError("Erro de comunicação ao buscar condição de pagamento.");
        }
    };

    // Listener de submit: cria (POST) ou atualiza (PUT) conforme houver ID no campo hidden
    const formEditar = document.getElementById('formInlineCondicaoPagamento');
    if (formEditar) {
        formEditar.addEventListener('submit', async (e) => {
            e.preventDefault();
            const token = getToken();
            const id = document.getElementById('edit_cp_id').value;

            const checkboxes = document.querySelectorAll('.fp-checkbox:checked');
            const formasPagamento = Array.from(checkboxes).map(cb => parseInt(cb.value, 10));

            const payload = {
                descricao: document.getElementById('edit_cp_descricao').value,
                formas_pagamento: formasPagamento,
                qtd_parcelas: parseInt(document.getElementById('edit_cp_qtd_parcelas').value, 10),
                dias_primeiro_venc: parseInt(document.getElementById('edit_cp_dias_venc').value, 10),
                intervalo_parcelas: parseInt(document.getElementById('edit_cp_intervalo').value, 10)
            };

            const method   = id ? 'PUT' : 'POST';
            const endpoint = id ? `/api/condicoes-pagamento/${id}` : '/api/condicoes-pagamento';

            try {
                const res = await fetch(endpoint, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(payload)
                });

                if (!res.ok) {
                    const data = await res.json().catch(() => ({}));
                    showError(data.erro || "Erro ao salvar condição de pagamento.");
                    return;
                }

                if (window.fecharFormularioInlineCondicaoPagamento) {
                    window.fecharFormularioInlineCondicaoPagamento();
                }

                carregarCondicoesPagamento();
            } catch (err) {
                console.error(err);
                showError("Erro interno ao comunicar com o servidor.");
            }
        });
    }
}
