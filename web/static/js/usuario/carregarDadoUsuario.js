import { getToken } from '/static/js/utils/auth.js';
import { showError } from '/static/js/utils/showError.js';
// A função validaRespostaRequisicao precisa ser importada de algum lugar.
// Assumindo que ela existe em '/static/js/utils/resposta.js'.
import { validaRespostaRequisicao } from '/static/js/utils/resposta.js';

async function carregarDadosUsuarioAPI() {
    const token = getToken();
    // A rota /api/usuario já está configurada para buscar o usuário pelo token.
    const res = await fetch(`/api/usuario`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });

    return await validaRespostaRequisicao(res);
}

export async function carregarDadosUsuario() {
    const form = document.getElementById('formInfoPessoais');
    if (!form) return;

    // Seleciona os inputs pelos IDs corretos do HTML
    const nomeInput = document.getElementById('nome');
    const emailInput = document.getElementById('email');
    const telefoneInput = document.getElementById('telefone');

    try {
        const dadosUsuario = await carregarDadosUsuarioAPI();

        if (dadosUsuario) {
            nomeInput.value = dadosUsuario.nome || '';
            emailInput.value = dadosUsuario.email || '';
            // O backend retorna um ponteiro para string, que pode ser null.
            // O `|| ''` garante que o campo não exiba "null".
            telefoneInput.value = dadosUsuario.telefone || '';
        }

    } catch (err) {
        showError(err.message || "Erro ao carregar dados do usuário.");
        nomeInput.value = "Erro ao carregar";
        emailInput.value = "Erro ao carregar";
        telefoneInput.placeholder = "Erro ao carregar";
    }
}