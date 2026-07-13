import { getToken } from '../utils/auth.js';
import { carregarFormasPagamento } from './listarFormasPagamento.js';
import { showError } from '../utils/showError.js';

export function setupEditarFormaPagamento() {

    // Abre o formulário inline preenchido com os dados da forma de pagamento buscada por ID
    window.abrirEditarFormaPagamento = async function(id) {
        const token = getToken();
        try {
            const res = await fetch(`/api/formas-pagamento/${id}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (!res.ok) {
                showError("Erro ao buscar forma de pagamento.");
                return;
            }

            const fp = await res.json();

            document.getElementById('edit_fp_id').value = fp.id;
            document.getElementById('edit_fp_descricao').value = fp.descricao || '';

            // Atualiza título do painel
            const titulo = document.getElementById('formFormaPagamentoTitulo');
            if (titulo) {
                titulo.innerHTML = `<i class="bi bi-pencil-square me-2"></i> Editando Forma de Pagamento #${id}`;
            }

            if (window.abrirFormularioInlineFormaPagamento) {
                window.abrirFormularioInlineFormaPagamento();
            }
        } catch (err) {
            console.error(err);
            showError("Erro de comunicação ao buscar forma de pagamento.");
        }
    };

    // Listener de submit: cria (POST) ou atualiza (PUT) conforme houver ID no campo hidden
    const formEditar = document.getElementById('formInlineFormaPagamento');
    if (formEditar) {
        formEditar.addEventListener('submit', async (e) => {
            e.preventDefault();
            const token = getToken();
            const id = document.getElementById('edit_fp_id').value;

            const payload = {
                descricao: document.getElementById('edit_fp_descricao').value
            };

            const method   = id ? 'PUT'                        : 'POST';
            const endpoint = id ? `/api/formas-pagamento/${id}` : '/api/formas-pagamento';

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
                    showError(data.erro || "Erro ao salvar forma de pagamento.");
                    return;
                }

                if (window.fecharFormularioInlineFormaPagamento) {
                    window.fecharFormularioInlineFormaPagamento();
                }

                carregarFormasPagamento();
            } catch (err) {
                console.error(err);
                showError("Erro interno ao comunicar com o servidor.");
            }
        });
    }
}
