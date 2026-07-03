import { getToken } from '/static/js/utils/auth.js';
import { showError } from '/static/js/utils/showError.js'; 
import { carregarClientes } from '/static/js/clientes/listarClientes.js'; // Garanta que a p

export function setupCriarCliente() {
    const formNovo = document.getElementById('formNovoCliente');
    if (!formNovo) return;

    formNovo.addEventListener('submit', async (e) => {
        e.preventDefault();

        const token = getToken();

        const nome = document.getElementById('cliente_nome').value;
        const tipo = document.getElementById('cliente_tipo').value;
        const cpf_cnpj = document.getElementById('cliente_cpf_cnpj').value;
        const email = document.getElementById('cliente_email').value;
        const telefone = document.getElementById('cliente_telefone').value;

        const payload = {
            nome,
            tipo,
            email,
            telefone
        };

        if (tipo === 'PF') {
            payload.cpf = cpf_cnpj;
        } else if (tipo === 'PJ') {
            payload.cnpj = cpf_cnpj;
        }

        try {
            const res = await fetch('/api/clientes', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(payload)
            });

            if (!res.ok) {
                const data = await res.json();
                showError(data.erro || "Erro ao cadastrar cliente.")
                return;
            }

            const modalEl = document.getElementById('modalCliente');
            const modal = bootstrap.Modal.getInstance(modalEl);
            if (modal) modal.hide();

            formNovo.reset();

            const tbody = document.getElementById('tabela_clientes_body');
            if (tbody) {
                carregarClientes();
            } else {
                window.location.reload();
            }
        } catch (err) {
            console.error(err);
            showError("Erro interno ao comunicar com servidor")
        }
    });
}