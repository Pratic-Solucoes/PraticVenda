import { getToken } from '../utils/auth.js';

export async function carregarDropdowns() {
    const token = getToken();
    const selectFornecedor = document.getElementById('id_fornecedor');
    const selectCategoria = document.getElementById('id_categoria');
    const editSelectFornecedor = document.getElementById('edit_id_fornecedor');
    const editSelectCategoria = document.getElementById('edit_id_categoria');

    try {
        // Fetch fornecedores
        const resF = await fetch('/api/fornecedores', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (resF.ok) {
            const fornecedores = await resF.json();
            if (selectFornecedor) selectFornecedor.innerHTML = '<option value="" selected disabled>Selecione o Fornecedor...</option>';
            if (editSelectFornecedor) editSelectFornecedor.innerHTML = '<option value="" selected disabled>Selecione o Fornecedor...</option>';

            fornecedores.forEach(f => {
                const opt = document.createElement('option');
                opt.value = f.id;
                opt.textContent = `${f.razao_social} (CNPJ: ${f.cnpj})`;
                if (selectFornecedor) selectFornecedor.appendChild(opt);

                if (editSelectFornecedor) {
                    const opt2 = document.createElement('option');
                    opt2.value = f.id;
                    opt2.textContent = `${f.razao_social} (CNPJ: ${f.cnpj})`;
                    editSelectFornecedor.appendChild(opt2);
                }
            });
        } else {
            if (selectFornecedor) selectFornecedor.innerHTML = '<option value="" disabled>Erro ao carregar</option>';
        }

        // Fetch categorias
        const resC = await fetch('/api/categorias', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (resC.ok) {
            const categorias = await resC.json();
            if (selectCategoria) selectCategoria.innerHTML = '<option value="" selected>Sem Categoria (Opcional)</option>';
            if (editSelectCategoria) editSelectCategoria.innerHTML = '<option value="" selected>Sem Categoria (Opcional)</option>';

            categorias.forEach(c => {
                const opt = document.createElement('option');
                opt.value = c.id;
                opt.textContent = c.nome;
                if (selectCategoria) selectCategoria.appendChild(opt);

                if (editSelectCategoria) {
                    const opt2 = document.createElement('option');
                    opt2.value = c.id;
                    opt2.textContent = c.nome;
                    editSelectCategoria.appendChild(opt2);
                }
            });
        } else {
            if (selectCategoria) selectCategoria.innerHTML = '<option value="" selected>Erro ao carregar categorias</option>';
        }

    } catch (err) {
        console.error(err);
        if (selectFornecedor) selectFornecedor.innerHTML = '<option value="" disabled>Erro ao carregar</option>';
    }
}
