import { getToken } from '../utils/auth.js';
import { carregarEstoques } from './listarEstoque.js';
import { showError } from '../utils/showError.js';

export function setupCriarEstoque() {
    const formNovo = document.getElementById('formNovoEstoque');
    if (!formNovo) return;

    formNovo.addEventListener('submit', async (e) => {
        e.preventDefault();
        const token = getToken();
        
        const nome = document.getElementById('estoque_nome').value;
        const descricaoEl = document.getElementById('estoque_descricao');
        const descricao = descricaoEl ? descricaoEl.value : '';

        try {
            const res = await fetch('/api/estoques', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ nome, descricao: descricao || null })
            });

            if (!res.ok) {
                const data = await res.json();
                showError(data.erro || "Erro ao cadastrar local de estoque.");
                return;
            }

            const modalEl = document.getElementById('modalEstoque');
            const modal = bootstrap.Modal.getInstance(modalEl);
            if (modal) modal.hide();

            formNovo.reset();
            // Recarregar os estoques no dropdown
            const data = await res.json();
            await carregarEstoques(data.id);

        } catch (err) {
            console.error(err);
            showError("Erro interno ao comunicar com o servidor.");
        }
    });
}
