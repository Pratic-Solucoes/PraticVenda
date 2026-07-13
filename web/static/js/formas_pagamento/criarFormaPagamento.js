import { getToken } from '../utils/auth.js';
import { carregarFormasPagamento } from './listarFormasPagamento.js';
import { showError } from '../utils/showError.js';

export function setupCriarFormaPagamento() {
    const formNovo = document.getElementById('formNovaFormaPagamento');
    if (!formNovo) return;

    formNovo.addEventListener('submit', async (e) => {
        e.preventDefault();
        const token = getToken();
        const descricao = document.getElementById('forma_pagamento_descricao').value;

        try {
            const res = await fetch('/api/formas-pagamento', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ descricao })
            });

            if (!res.ok) {
                const data = await res.json();
                showError(data.erro || "Erro ao cadastrar forma de pagamento.");
                return;
            }

            const modalEl = document.getElementById('modalFormaPagamento');
            const modal = bootstrap.Modal.getInstance(modalEl);
            if (modal) modal.hide();

            formNovo.reset();
            carregarFormasPagamento();

        } catch (err) {
            console.error(err);
            showError("Erro interno ao comunicar com o servidor.");
        }
    });
}
