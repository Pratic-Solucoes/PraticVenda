import { getToken } from '../utils/auth.js';
import { carregarDebitos } from './listarDebitos.js';

export function setupPagarDebito() {
    window.pagarDebito = async function (id) {
        if (!confirm("Tem certeza que deseja dar baixa (pagar) este débito?")) return;

        const token = getToken();

        try {
            const res = await fetch(`/api/debitos/${id}/pagar`, {
                method: 'PUT',
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (!res.ok) {
                const data = await res.json();
                alert(data.erro || "Erro ao pagar débito.");
                return;
            }

            alert("Débito pago com sucesso!");
            carregarDebitos();
        } catch (err) {
            console.error(err);
            alert("Erro interno.");
        }
    };
}
