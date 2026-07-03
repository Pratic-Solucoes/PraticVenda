export async function validaRespostaRequisicao(resposta) {
    
    if (resposta.ok) {
        return await resposta.json();
    }

    let mensagemRetorno = '';
    try{
        const dataRetorno = await resposta.json();
        mensagemRetorno = dataRetorno.erro || dataRetorno.message || '';
    }catch{}
 
    switch (resposta.status) {
        case 400:
            throw new Error(mensagemRetorno || 'Requisição inválida (400). Verifique os dados enviados.');
        case 401:
            throw new Error(mensagemRetorno || 'Não autorizado (401). Faça login novamente.');
        case 404:
            throw new Error(mensagemRetorno || 'Recurso não encontrado (404).');
        case 500:
            throw new Error(mensagemRetorno || 'Erro interno no servidor (500). Tente mais tarde.');
        default:
            throw new Error(`Erro inesperado: Status ${resposta.status}`);
        }
}