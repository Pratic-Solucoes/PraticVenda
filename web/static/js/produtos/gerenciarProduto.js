import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';
import { carregarProdutos } from './listarProdutos.js';

let availableEstoques = [];

export async function loadGruposTributarios() {
    const selectGrupo = document.getElementById('prod_grupo_tributario');
    if (!selectGrupo) return;

    const token = getToken();
    try {
        const res = await fetch('/api/grupos-tributarios', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!res.ok) {
            throw new Error("Erro ao buscar grupos tributários");
        }

        const grupos = await res.json();
        selectGrupo.innerHTML = '<option value="" disabled selected>Selecione um grupo...</option>';
        grupos.forEach(g => {
            const opt = document.createElement('option');
            opt.value = g.id;
            opt.textContent = `${g.nome} (CFOP: ${g.cfop_padrao})`;
            selectGrupo.appendChild(opt);
        });
    } catch (err) {
        console.error(err);
        showError("Erro ao carregar grupos tributários.");
    }
}

export async function loadEstoquesDisponiveis() {
    const tbody = document.getElementById('listaEstoquesVinculoBody');
    if (!tbody) return;

    const token = getToken();
    try {
        const res = await fetch('/api/estoques', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!res.ok) {
            throw new Error("Erro ao buscar estoques");
        }

        availableEstoques = await res.json();
        renderEstoquesTable();
    } catch (err) {
        console.error(err);
        showError("Erro ao carregar locais de estoque disponíveis.");
    }
}

function renderEstoquesTable(vinculosAtuais = []) {
    const tbody = document.getElementById('listaEstoquesVinculoBody');
    if (!tbody) return;

    tbody.innerHTML = '';

    if (availableEstoques.length === 0) {
        tbody.innerHTML = '<tr><td colspan="4" class="text-center text-muted py-3">Nenhum local de estoque cadastrado.</td></tr>';
        return;
    }

    availableEstoques.forEach(est => {
        const vinculo = vinculosAtuais.find(v => v.id_estoque === est.id);
        const checks = vinculo ? 'checked' : '';
        const minEst = vinculo ? vinculo.estoque_minimo : 0;
        const qtd = vinculo ? vinculo.quantidade : 0;
        
        // Se for edição, a quantidade geralmente não deve ser editada livremente pelo cadastro do produto (já que tem movimentos de estoque),
        // mas vamos deixar como campo apenas de leitura ou desabilitado caso seja edição para evitar furos de estoque físicos,
        // ou seja, se for produto novo, permite digitar. Se for edição, mostramos e desabilitamos o input de quantidade (apenas leitura).
        const isEdit = document.getElementById('edit_produto_id').value !== "";
        const qtdDisabled = isEdit ? 'disabled' : '';

        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>
                <input type="checkbox" class="form-check-input check-vinculo" data-estoque-id="${est.id}" ${checks}>
            </td>
            <td class="fw-semibold">${est.nome}</td>
            <td>
                <input type="number" step="0.001" class="form-control form-control-sm input-min" style="width: 100px;" value="${minEst}">
            </td>
            <td>
                <input type="number" step="0.001" class="form-control form-control-sm input-qtd" style="width: 100px;" value="${qtd}" ${qtdDisabled}>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

export function abrirFormularioNovo(bsCollapse) {
    const form = document.getElementById('formInlineProduto');
    if (form) form.reset();

    document.getElementById('edit_produto_id').value = "";
    document.getElementById('formProdutoTitulo').innerHTML = '<i class="bi bi-box-seam me-2"></i> Novo Produto';
    
    // Reset tabs
    const basicoTab = document.getElementById("comercial-tab");
    if (basicoTab) bootstrap.Tab.getOrCreateInstance(basicoTab).show();

    // Renderizar tabela de estoques limpa
    renderEstoquesTable([]);

    if (bsCollapse) bsCollapse.show();
}

export function fecharFormulario(bsCollapse) {
    if (bsCollapse) bsCollapse.hide();
    const form = document.getElementById('formInlineProduto');
    if (form) form.reset();
    document.getElementById('edit_produto_id').value = "";
}

export async function carregarProdutoParaEdicao(id) {
    const token = getToken();
    try {
        const res = await fetch(`/api/produtos/${id}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!res.ok) {
            throw new Error("Erro ao buscar dados do produto");
        }

        const data = await res.json();
        
        document.getElementById('edit_produto_id').value = data.id;
        document.getElementById('formProdutoTitulo').innerHTML = `<i class="bi bi-box-seam me-2"></i> Editar Produto: ${data.nome}`;
        
        // Preecher comercial
        document.getElementById('prod_nome').value = data.nome;
        document.getElementById('prod_codigo_interno').value = data.codigo_interno_loja || '';
        document.getElementById('prod_codigo_barras').value = data.codigo_barras || '';
        document.getElementById('prod_preco_custo').value = data.preco_custo;
        document.getElementById('prod_preco_venda').value = data.preco_venda;
        document.getElementById('prod_unidade_estoque').value = data.unidade_estoque;
        document.getElementById('prod_unidade_venda').value = data.unidade_venda;
        document.getElementById('prod_peso_bruto').value = data.peso_bruto;
        document.getElementById('prod_peso_liquido').value = data.peso_liquido;
        document.getElementById('prod_ativo').checked = data.ativo;
        document.getElementById('prod_descricao').value = data.descricao || '';

        // Preencher fiscal
        if (data.fiscal) {
            document.getElementById('prod_ncm').value = data.fiscal.ncm;
            document.getElementById('prod_cest').value = data.fiscal.cest || '';
            document.getElementById('prod_grupo_tributario').value = data.fiscal.id_grupo_tributario;
        }

        // Renderizar tabela de estoques com vínculos atuais
        renderEstoquesTable(data.estoques || []);

        // Mostrar formulário
        if (window.abrirFormularioEditarProduto) {
            window.abrirFormularioEditarProduto();
        }
    } catch (err) {
        console.error(err);
        showError("Erro ao carregar detalhes do produto.");
    }
}

export async function excluirOuInativarProduto(id, nome) {
    if (!confirm(`Deseja realmente remover ou inativar o produto "${nome}"?`)) {
        return;
    }

    const token = getToken();
    try {
        const res = await fetch(`/api/produtos/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        const data = await res.json();
        if (!res.ok) {
            showError(data.erro || "Erro ao processar exclusão do produto.");
            return;
        }

        alert(`Produto processado com sucesso! Ação realizada: ${data.status.toUpperCase()}`);
        carregarProdutos();
    } catch (err) {
        console.error(err);
        showError("Erro interno ao tentar remover o produto.");
    }
}

export function setupGerenciarProduto() {
    loadGruposTributarios();
    loadEstoquesDisponiveis();

    const form = document.getElementById('formInlineProduto');
    if (!form) return;

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const token = getToken();
        const id = document.getElementById('edit_produto_id').value;

        // Validação manual de campos fiscais na aba oculta para evitar erro de foco
        const ncm = document.getElementById('prod_ncm').value.trim();
        const grupoTributario = document.getElementById('prod_grupo_tributario').value;

        if (!ncm) {
            showError("O campo NCM é obrigatório.");
            const fiscalTab = document.getElementById("fiscal-tab");
            if (fiscalTab) bootstrap.Tab.getOrCreateInstance(fiscalTab).show();
            document.getElementById('prod_ncm').focus();
            return;
        }

        if (!grupoTributario) {
            showError("Selecione um Grupo Tributário.");
            const fiscalTab = document.getElementById("fiscal-tab");
            if (fiscalTab) bootstrap.Tab.getOrCreateInstance(fiscalTab).show();
            document.getElementById('prod_grupo_tributario').focus();
            return;
        }

        // Montar a lista de estoques vinculados selecionados
        const estoquesVinculados = [];
        const rows = document.querySelectorAll('#listaEstoquesVinculoBody tr');
        rows.forEach(row => {
            const check = row.querySelector('.check-vinculo');
            if (check && check.checked) {
                const idEstoque = parseInt(check.dataset.estoqueId, 10);
                const inputMin = row.querySelector('.input-min');
                const inputQtd = row.querySelector('.input-qtd');
                
                estoquesVinculados.push({
                    id_estoque: idEstoque,
                    estoque_minimo: parseFloat(inputMin.value) || 0,
                    quantidade: parseFloat(inputQtd.value) || 0
                });
            }
        });

        const payload = {
            codigo_barras: document.getElementById('prod_codigo_barras').value.trim() || null,
            codigo_interno_loja: document.getElementById('prod_codigo_interno').value.trim() || null,
            nome: document.getElementById('prod_nome').value.trim(),
            descricao: document.getElementById('prod_descricao').value.trim() || null,
            preco_custo: parseFloat(document.getElementById('prod_preco_custo').value) || 0,
            preco_venda: parseFloat(document.getElementById('prod_preco_venda').value) || 0,
            unidade_estoque: document.getElementById('prod_unidade_estoque').value.trim(),
            unidade_venda: document.getElementById('prod_unidade_venda').value.trim(),
            peso_bruto: parseFloat(document.getElementById('prod_peso_bruto').value) || 0,
            peso_liquido: parseFloat(document.getElementById('prod_peso_liquido').value) || 0,
            ncm: document.getElementById('prod_ncm').value.trim(),
            cest: document.getElementById('prod_cest').value.trim() || null,
            id_grupo_tributario: parseInt(document.getElementById('prod_grupo_tributario').value, 10),
            estoques: estoquesVinculados
        };

        const url = id ? `/api/produtos/${id}` : '/api/produtos';
        const method = id ? 'PUT' : 'POST';

        try {
            const res = await fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(payload)
            });

            const data = await res.json();
            if (!res.ok) {
                showError(data.erro || "Erro ao salvar produto.");
                return;
            }

            alert(id ? "Produto atualizado com sucesso!" : "Produto cadastrado com sucesso!");
            
            // Fechar formulário
            const collapseElement = document.getElementById('collapseFormProduto');
            if (collapseElement) {
                const bsCollapse = bootstrap.Collapse.getInstance(collapseElement);
                if (bsCollapse) bsCollapse.hide();
            }
            form.reset();
            document.getElementById('edit_produto_id').value = "";

            carregarProdutos();
        } catch (err) {
            console.error(err);
            showError("Erro interno ao tentar salvar o produto.");
        }
    });
}
