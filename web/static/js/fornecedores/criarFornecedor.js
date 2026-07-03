import { getToken } from '/static/js/utils/auth.js';
import { carregarFornecedores} from '/static/js/fornecedores/listarFornecedores.js';
import { showError } from '/static/js/utils/showError.js';
import { fecharModal } from '/static/js/utils/fecharModal.js';

async function criarFornecedorAPI(dados){

    const token = getToken();
    const res = await fetch('/api/fornecedores',{
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(dados)
    });
    if(!res.ok){
        const data = await res.json();
        throw new Error(data.erro || "Erro ao cadastrar fornecedor.");
    }
    return await res.json();
}

function recarregarDados(){
    const tbody = document.getElementById('tabela_fornecedores_body');
    if(tbody){
        carregarFornecedores();
    }else{
        window.location.reload();
    }
}

export function setupCriarFornecedor() {
    // 1. Busca o formulário na tela. Se não existir, interrompe a execução para evitar erros.
    const formNovo = document.getElementById('formNovoFornecedor');
    if (!formNovo) return;

    // 2. Intercepta o evento de envio (submit) do formulário
    formNovo.addEventListener('submit', async (event) => {
        event.preventDefault(); // Evita o envio padrão do formulário
        
        // 3. Coleta os dados do formulário
        const formData = new FormData(formNovo);
        const dadosFornecedor = {
            razao_social: formData.get('fornecedor_razao_social'),
            cnpj: formData.get('fornecedor_cnpj'),
            email: formData.get('fornecedor_email')
        };

        console.log('Dados do fornecedor a serem enviados:', dadosFornecedor); // Log para depuração

        try{
            // 4. Chama a função que faz a requisição para criar o fornecedor
            await criarFornecedorAPI(dadosFornecedor);
            // 5. Fecha o modal e recarrega os dados da tabela
            formNovo.reset(); // Limpa o formulário após o envio
            fecharModal('modalNovoFornecedor');
            recarregarDados();
        }catch(err){        
            showError(err.message); 
        }
    })
}
