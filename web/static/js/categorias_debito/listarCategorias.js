import { getToken } from '../utils/auth.js';

function showError(message) {
    const modalBody = document.getElementById('errorModalBody');
    if (modalBody) {
        modalBody.textContent = message;
        const errorModalElement = document.getElementById('errorModal');
        const errorModal = bootstrap.Modal.getOrCreateInstance(errorModalElement);
        errorModal.show();
    } else {
        alert(message);
    }
}

export async function carregarCategorias() {
    const tbody = document.getElementById('tabela_categorias_body');
    if (!tbody) return;
    
    const token = getToken();
    tbody.innerHTML = `<tr><td colspan="3" class="text-muted py-5 text-center">Carregando...</td></tr>`;
    
    try {
        const res = await fetch('/api/categorias', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (res.status === 401) {
            window.location.href = '/';
            return;
        }

        if (!res.ok) {
            showError("Erro ao carregar lista de categorias.");
            tbody.innerHTML = `<tr><td colspan="3" class="text-danger py-5 text-center">Erro ao carregar.</td></tr>`;
            return;
        }

        const dados = await res.json();
        renderTabela(dados);

    } catch (err) {
        console.error(err);
        tbody.innerHTML = `<tr><td colspan="3" class="text-danger py-5 text-center">Erro de comunicação.</td></tr>`;
    }
}

function renderTabela(categorias) {
    const tbody = document.getElementById('tabela_categorias_body');
    if (!tbody) return;

    if (!categorias || categorias.length === 0) {
        tbody.innerHTML = `<tr><td colspan="3" class="text-muted py-5 text-center">Nenhuma categoria encontrada.</td></tr>`;
        return;
    }

    tbody.innerHTML = '';
    categorias.forEach(c => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>#${c.id}</td>
            <td class="text-start fw-bold">${c.nome}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary" title="Editar"><i class="bi bi-pencil"></i></button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

export { showError };
