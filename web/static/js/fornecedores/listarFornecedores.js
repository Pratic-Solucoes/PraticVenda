import { getToken } from '/static/js/utils/auth.js';
import { showError } from '/static/js/utils/showError.js';
import { validaRespostaRequisicao } from '/static/js/utils/resposta.js';

async function carregarFornecedoresAPI(busca = ""){

    const token = getToken();

    let url = '/api/fornecedores';
    if(busca){
        url += `?busca=${encodeURIComponent(busca)}`;
    }

    const res = await fetch(url, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })

    return await validaRespostaRequisicao(res);
}

export async function carregarFornecedores(busca = "") {
    const tbody = document.getElementById('tabela_fornecedores_body');
    if (!tbody) return;
    
    tbody.innerHTML = `<tr><td colspan="5" class="text-muted py-5 text-center">Carregando...</td></tr>`;
    
    try {
        
        const dados = await carregarFornecedoresAPI(busca);
        renderTabela(dados);

    } catch (err) {
        showError(err.message || 'Erro ao carregar fornecedores. Tente novamente mais tarde.');
    }
}

function renderTabela(fornecedores) {
    const tbody = document.getElementById('tabela_fornecedores_body');
    if (!tbody) return;

    if (!fornecedores || fornecedores.length === 0) {
        tbody.innerHTML = `<tr><td colspan="5" class="text-muted py-5 text-center">Nenhum fornecedor encontrado.</td></tr>`;
        return;
    }

    tbody.innerHTML = '';
    fornecedores.forEach(f => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>#${f.id}</td>
            <td class="text-start fw-bold">${f.razao_social}</td>
            <td>${f.cnpj}</td>
            <td>${f.email || '-'}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary" title="Editar" onclick="window.abrirModalEditarFornecedor(${f.id})"><i class="bi bi-pencil"></i></button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}
