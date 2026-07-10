import { getToken } from '../utils/auth.js';

export function setupCriarFornecedorModal() {
    const formModalFornecedor = document.getElementById('formModalFornecedor');
    if (!formModalFornecedor) return;

    formModalFornecedor.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const token = getToken();
        const payload = {
            razao_social: document.getElementById('modal_fornecedor_razao_social').value,
            cnpj: document.getElementById('modal_fornecedor_cnpj').value,
            inscricao_estadual: document.getElementById('modal_fornecedor_inscricao_estadual').value || null,
            email: document.getElementById('modal_fornecedor_email').value
        };

        try {
            const res = await fetch('/api/fornecedores', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(payload)
            });

            if (!res.ok) {
                const data = await res.json();
                alert(data.erro || "Erro ao cadastrar fornecedor.");
                return;
            }

            alert("Fornecedor cadastrado com sucesso!");
            
            const modalEl = document.getElementById('modalFornecedor');
            const modal = bootstrap.Modal.getInstance(modalEl);
            if (modal) modal.hide();
            
            formModalFornecedor.reset();
            
            // Reload page if needed to update lists or selects
            const selects = document.querySelectorAll('#id_fornecedor');
            if (selects.length > 0) {
                 // Try to reload choices without page reload
                 window.location.reload(); 
            } else {
                 window.location.reload();
            }

        } catch (err) {
            console.error(err);
            alert("Erro interno ao comunicar com o servidor.");
        }
    });
}
