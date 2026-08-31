import { checkAuth, getToken } from '../utils/auth.js';
import { showError, showSuccess } from '../utils/showError.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    const headers = { Authorization: `Bearer ${getToken()}`, 'Content-Type': 'application/json' };
    const form = document.getElementById('formCriarCaixa');
    const selectUsuario = document.getElementById('caixa_usuario');
    const listaCaixas = document.getElementById('listaCaixas');

    async function carregarCaixas() {
        const [usuariosRes, caixasRes] = await Promise.all([
            fetch('/api/caixas/usuarios', { headers }),
            fetch('/api/caixas', { headers }),
        ]);
        if (!usuariosRes.ok || !caixasRes.ok) {
            showError('Não foi possível carregar as configurações de caixa.');
            return;
        }

        const [usuarios, caixas] = await Promise.all([usuariosRes.json(), caixasRes.json()]);
        selectUsuario.innerHTML = '<option value="">Selecione...</option>';
        usuarios.forEach(usuario => {
            selectUsuario.innerHTML += `<option value="${usuario.id}">${usuario.nome}</option>`;
        });
        listaCaixas.innerHTML = caixas.length
            ? caixas.map(caixa => `<div class="border rounded p-2 mb-2"><i class="bi bi-cash me-2"></i>${caixa.nome}</div>`).join('')
            : 'Nenhum caixa vinculado ao seu usuário.';
    }

    form.addEventListener('submit', async event => {
        event.preventDefault();
        const response = await fetch('/api/caixas', {
            method: 'POST',
            headers,
            body: JSON.stringify({
                nome: document.getElementById('caixa_nome').value,
                id_usuario: Number(selectUsuario.value),
            }),
        });
        if (!response.ok) {
            const data = await response.json();
            showError(data.erro || 'Não foi possível criar o caixa.');
            return;
        }
        form.reset();
        showSuccess('Caixa criado com sucesso.');
        carregarCaixas();
    });

    carregarCaixas();
});
