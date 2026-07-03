import { getToken } from '../utils/auth.js';
import { carregarCategorias, showError } from './listarCategorias.js';

export function setupCriarCategoria() {
    const formNovo = document.getElementById('formNovaCategoria');
    if (!formNovo) return;

    formNovo.addEventListener('submit', async (e) => {
        e.preventDefault();
        const token = getToken();
        const nome = document.getElementById('categoria_nome').value;

        try {
            const res = await fetch('/api/categorias', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ nome })
            });

            if (!res.ok) {
                const data = await res.json();
                showError(data.erro || "Erro ao cadastrar categoria.");
                return;
            }

            const modalEl = document.getElementById('modalCategoria');
            const modal = bootstrap.Modal.getInstance(modalEl);
            if (modal) modal.hide();

            formNovo.reset();
            carregarCategorias();

        } catch (err) {
            console.error(err);
            showError("Erro interno ao comunicar com o servidor.");
        }
    });
}
