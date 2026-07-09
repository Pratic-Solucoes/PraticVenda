import { getToken } from '../utils/auth.js';
import { carregarContasPagar } from './listarContasPagar.js';
import { state } from './state.js';

export function setupPagarContaPagar() {
    window.abrirModalPagamento = function (id) {
        // Busca a conta no state
        const conta = state.contasPagarCarregadas.find(c => c.id === id);
        if (!conta) {
            alert("Conta não encontrada.");
            return;
        }

        document.getElementById('pagar_conta_id').value = id;
        document.getElementById('pagar_saldo_restante').value = new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(conta.saldo_restante);
        
        // Define o max do input para não pagar mais do que o saldo
        const inputValor = document.getElementById('pagar_valor');
        inputValor.value = conta.saldo_restante.toFixed(2);
        inputValor.max = conta.saldo_restante.toFixed(2);

        const modal = new bootstrap.Modal(document.getElementById('modalPagarConta'));
        modal.show();
    };

    const formPagarConta = document.getElementById('formPagarConta');
    if (formPagarConta) {
        formPagarConta.addEventListener('submit', async function(e) {
            e.preventDefault();

            const id = document.getElementById('pagar_conta_id').value;
            const valor_pagamento = parseFloat(document.getElementById('pagar_valor').value);

            if (isNaN(valor_pagamento) || valor_pagamento <= 0) {
                alert("Por favor, insira um valor válido.");
                return;
            }

            const token = getToken();
            const btnSubmit = document.querySelector('button[form="formPagarConta"]');
            const originalText = btnSubmit.innerHTML;
            btnSubmit.innerHTML = 'Processando...';
            btnSubmit.disabled = true;

            try {
                const res = await fetch(`/api/contas-pagar/${id}/pagar`, {
                    method: 'PUT',
                    headers: { 
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ valor_pagamento })
                });

                if (!res.ok) {
                    const data = await res.json();
                    alert(data.erro || "Erro ao realizar o pagamento.");
                    return;
                }

                alert("Pagamento realizado com sucesso!");
                bootstrap.Modal.getInstance(document.getElementById('modalPagarConta')).hide();
                carregarContasPagar();
            } catch (err) {
                console.error(err);
                alert("Erro interno de comunicação.");
            } finally {
                btnSubmit.innerHTML = originalText;
                btnSubmit.disabled = false;
            }
        });
    }
}
