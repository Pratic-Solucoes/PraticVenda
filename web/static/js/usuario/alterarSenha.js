import { getToken } from '/static/js/utils/auth.js';
import { showError } from '/static/js/utils/showError.js';
import { validaRespostaRequisicao } from '/static/js/utils/resposta.js';

async function alterarSenhaAPI(senhaAtual, novaSenha, senhaConfirmacao) {
    const token = getToken();
    const res = await fetch(`/api/usuario/alterar-senha`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            senha_atual: senhaAtual,
            nova_senha: novaSenha,
            senha_confirmacao: senhaConfirmacao
        })
    });

    // Se a resposta não for bem-sucedida, delegamos para a função que extrai a mensagem de erro do JSON.
    if (!res.ok) {
        return await validaRespostaRequisicao(res);
    }

    // Se a resposta for bem-sucedida (200 OK), o backend não retorna corpo.
    // Não há nada para ser processado, então a função pode simplesmente terminar.
}

export async function alterarSenha(){

    const formAlterarSenha = document.getElementById('formAlterarSenha');
    if(!formAlterarSenha) return;

    formAlterarSenha.addEventListener('submit', async (e) => {
        e.preventDefault();

        const senhaAtualInput = document.getElementById('senha_atual');
        const novaSenhaInput = document.getElementById('nova_senha');
        const senhaConfirmacaoInput = document.getElementById('senha_confirmacao');

        const senhaAtual = senhaAtualInput.value;
        const novaSenha = novaSenhaInput.value;
        const senhaConfirmacao = senhaConfirmacaoInput.value;

        try{
            await alterarSenhaAPI(senhaAtual, novaSenha, senhaConfirmacao);
            senhaAtualInput.value = '';
            novaSenhaInput.value = '';
            senhaConfirmacaoInput.value = '';
            alert("Senha alterada com sucesso!");
        } catch(err){
            showError(err.message || "Erro ao alterar senha.");
        }
    });
}