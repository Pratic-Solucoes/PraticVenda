import { getToken } from '../utils/auth.js';
import { showError, showSuccess } from '../utils/showError.js';
import { carregarProdutos } from './listarProdutos.js';

let availableEstoques = [];
let componentes = [];
let produtosParaComposicao = [];

function moeda(valor) { return Number(valor || 0).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' }); }
function renderizarComposicao() {
    const body = document.getElementById('listaComposicao');
    if (!body) return;
    body.innerHTML = componentes.map((item, indice) => `<tr><td>${item.nome_produto}</td><td>${item.unidade_estoque}</td><td class="text-end">${item.quantidade}</td><td class="text-end">${moeda(item.preco_custo * item.quantidade)}</td><td class="text-end"><button class="btn btn-sm btn-outline-danger remover-componente" data-indice="${indice}" type="button"><i class="bi bi-trash"></i></button></td></tr>`).join('') || '<tr><td colspan="5" class="text-center text-muted">Nenhuma matéria-prima adicionada.</td></tr>';
    document.getElementById('custoComposicao').textContent = moeda(componentes.reduce((total, item) => total + item.preco_custo * item.quantidade, 0));
    document.querySelectorAll('.remover-componente').forEach(botao => botao.onclick = () => { componentes.splice(botao.dataset.indice, 1); renderizarComposicao(); });
}
async function carregarComposicao(id) {
    const [produtosRes, composicaoRes] = await Promise.all([fetch('/api/produtos', { headers: { Authorization: `Bearer ${getToken()}` } }), fetch(`/api/produtos/${id}/composicao`, { headers: { Authorization: `Bearer ${getToken()}` } })]);
    produtosParaComposicao = await produtosRes.json(); componentes = await composicaoRes.json();
    document.getElementById('compProduto').innerHTML = produtosParaComposicao.filter(p => p.id !== Number(id) && p.materia_prima).map(p => `<option value="${p.id}">${p.nome} (${p.unidade_estoque})</option>`).join('');
    avisoComposicao.classList.add('d-none'); conteudoComposicao.classList.remove('d-none'); renderizarComposicao();
}

async function prepararComposicao(id = 0) {
    const res = await fetch('/api/produtos', { headers: { Authorization: `Bearer ${getToken()}` } });
    if (!res.ok) return showError('Não foi possível carregar as matérias-primas.');
    produtosParaComposicao = await res.json();
    document.getElementById('compProduto').innerHTML = produtosParaComposicao.filter(p => p.id !== Number(id) && p.materia_prima).map(p => `<option value="${p.id}">${p.nome} (${p.unidade_estoque})</option>`).join('');
    document.getElementById('avisoComposicao').classList.add('d-none');
    document.getElementById('conteudoComposicao').classList.remove('d-none');
    renderizarComposicao();
}

async function loadFornecedores() {
    const container = document.getElementById('prod_fornecedores');
    if (!container) return;
    try {
        const res = await fetch('/api/fornecedores', { headers: { Authorization: `Bearer ${getToken()}` } });
        if (!res.ok) throw new Error();
        const fornecedores = await res.json();
		container.innerHTML = fornecedores.length ? fornecedores.map(fornecedor => `<div class="form-check"><input class="form-check-input fornecedor-produto" type="checkbox" value="${fornecedor.id}" id="prod_fornecedor_${fornecedor.id}"><label class="form-check-label" for="prod_fornecedor_${fornecedor.id}">${fornecedor.razao_social}</label></div>`).join('') : '<span class="text-muted small">Nenhum fornecedor cadastrado.</span>';
    } catch (_) { showError('Erro ao carregar fornecedores.'); }
}

function atualizarPercentualLucro() {
    const custo = parseFloat(document.getElementById('prod_preco_custo').value) || 0;
    const venda = parseFloat(document.getElementById('prod_preco_venda').value) || 0;
    document.getElementById('prod_percentual_lucro').value = custo > 0 ? (((venda - custo) / custo) * 100).toFixed(2) : '';
}

function atualizarAbaComposicao() {
    const composto = document.getElementById('prod_composto')?.checked;
	const materiaPrima = document.getElementById('prod_materia_prima')?.checked;
	const salvo = Boolean(document.getElementById('edit_produto_id')?.value);
	const fornecedores = document.getElementById('campo_fornecedores');
	const estoquesNav = document.getElementById('estoques-nav');
	fornecedores?.classList.toggle('d-none', composto);
	estoquesNav?.classList.toggle('d-none', composto);
	document.getElementById('campo_preco_venda')?.classList.toggle('d-none', materiaPrima);
	document.getElementById('campo_percentual_lucro')?.classList.toggle('d-none', materiaPrima);
	document.getElementById('prod_preco_venda').required = !materiaPrima;
    document.getElementById('composicao-tab')?.classList.toggle('disabled', !composto);
    if (!composto) {
        document.getElementById('avisoComposicao')?.classList.remove('d-none');
        document.getElementById('conteudoComposicao')?.classList.add('d-none');
	} else if (salvo) {
		document.getElementById('avisoComposicao')?.classList.add('d-none');
		document.getElementById('conteudoComposicao')?.classList.remove('d-none');
    }
}

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
	componentes = [];
	atualizarAbaComposicao();
	componentes = []; avisoComposicao?.classList.remove('d-none'); conteudoComposicao?.classList.add('d-none');

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
		await loadFornecedores();
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
		const fornecedoresSelecionados = data.composto ? [] : (data.ids_fornecedores || [data.id_fornecedor]);
		document.querySelectorAll('.fornecedor-produto').forEach(input => input.checked = fornecedoresSelecionados.includes(Number(input.value)));
        document.getElementById('prod_codigo_interno').value = data.codigo_interno_loja || '';
        document.getElementById('prod_codigo_barras').value = data.codigo_barras || '';
        document.getElementById('prod_preco_custo').value = data.preco_custo;
        document.getElementById('prod_preco_venda').value = data.preco_venda;
        document.getElementById('prod_unidade_estoque').value = data.unidade_estoque;
        document.getElementById('prod_unidade_venda').value = data.unidade_venda;
        document.getElementById('prod_peso_bruto').value = data.peso_bruto;
        document.getElementById('prod_peso_liquido').value = data.peso_liquido;
        document.getElementById('prod_ativo').checked = data.ativo;
		document.getElementById('prod_composto').checked = Boolean(data.composto);
		document.getElementById('prod_materia_prima').checked = Boolean(data.materia_prima);
        document.getElementById('prod_descricao').value = data.descricao || '';
		atualizarPercentualLucro();

        // Preencher fiscal
        if (data.fiscal) {
            document.getElementById('prod_ncm').value = data.fiscal.ncm;
            document.getElementById('prod_cest').value = data.fiscal.cest || '';
            document.getElementById('prod_grupo_tributario').value = data.fiscal.id_grupo_tributario;
        }

        // Renderizar tabela de estoques com vínculos atuais
        renderEstoquesTable(data.estoques || []);
        await carregarComposicao(data.id);
		atualizarAbaComposicao();

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
	loadFornecedores();

    const form = document.getElementById('formInlineProduto');
    if (!form) return;
	const adicionarComponente = () => { const produto = produtosParaComposicao.find(p => p.id === Number(compProduto.value)); const quantidade = Number(compQuantidade.value); if (!produto || quantidade <= 0) return showError('Informe uma matéria-prima e quantidade válida.'); if (componentes.some(item => item.id_produto_componente === produto.id)) return showError('Matéria-prima já adicionada.'); componentes.push({ id_produto_componente: produto.id, nome_produto: produto.nome, unidade_estoque: produto.unidade_estoque, preco_custo: produto.preco_custo, quantidade }); compQuantidade.value = ''; renderizarComposicao(); };
	document.getElementById('btnAdicionarComponente')?.addEventListener('click', adicionarComponente);
	document.getElementById('btnSalvarComposicao')?.addEventListener('click', async () => { const id = document.getElementById('edit_produto_id').value; if (!id || !componentes.length) return showError('Adicione ao menos uma matéria-prima.'); const res = await fetch(`/api/produtos/${id}/composicao`, { method:'PUT', headers:{'Content-Type':'application/json', Authorization:`Bearer ${getToken()}`}, body:JSON.stringify(componentes.map(({id_produto_componente,quantidade})=>({id_produto_componente,quantidade}))) }); const data = await res.json(); if(!res.ok)return showError(data.erro||'Erro ao salvar ficha técnica.'); showSuccess('Ficha técnica salva com sucesso.'); await carregarComposicao(id); });

	['prod_preco_custo', 'prod_preco_venda'].forEach(id => document.getElementById(id)?.addEventListener('input', atualizarPercentualLucro));
	document.getElementById('prod_composto')?.addEventListener('change', async event => { if (event.target.checked) document.getElementById('prod_materia_prima').checked=false; atualizarAbaComposicao(); if (event.target.checked) await prepararComposicao(document.getElementById('edit_produto_id').value); });
	document.getElementById('prod_materia_prima')?.addEventListener('change', event => { if (event.target.checked) { document.getElementById('prod_composto').checked=false; document.getElementById('prod_preco_venda').value='0.00'; } atualizarAbaComposicao(); });

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const token = getToken();
        const id = document.getElementById('edit_produto_id').value;

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

		const composto = document.getElementById('prod_composto').checked;
		const materiaPrima = document.getElementById('prod_materia_prima').checked;
		if (composto && componentes.length === 0) {
			showError('Adicione ao menos uma matéria-prima à ficha técnica.');
			bootstrap.Tab.getOrCreateInstance(document.getElementById('composicao-tab')).show();
			return;
		}
        if (!composto && estoquesVinculados.length === 0) {
            showError("Vincule o produto a pelo menos um estoque.");
            const estoquesTab = document.getElementById("estoques-tab");
            if (estoquesTab) bootstrap.Tab.getOrCreateInstance(estoquesTab).show();
            return;
        }

		const fornecedoresSelecionados = [...document.querySelectorAll('.fornecedor-produto:checked')].map(input => Number(input.value));
		if (!composto && fornecedoresSelecionados.length === 0) {
			showError('Selecione ao menos um fornecedor para o produto.');
			return;
		}

        const payload = {
			composto,
			materia_prima: materiaPrima,
			ids_fornecedores: composto ? [] : fornecedoresSelecionados,
            codigo_barras: document.getElementById('prod_codigo_barras').value.trim() || null,
            codigo_interno_loja: document.getElementById('prod_codigo_interno').value.trim() || null,
            nome: document.getElementById('prod_nome').value.trim(),
            descricao: document.getElementById('prod_descricao').value.trim() || null,
            preco_custo: parseFloat(document.getElementById('prod_preco_custo').value) || 0,
			preco_venda: materiaPrima ? 0 : (parseFloat(document.getElementById('prod_preco_venda').value) || 0),
            unidade_estoque: document.getElementById('prod_unidade_estoque').value.trim(),
            unidade_venda: document.getElementById('prod_unidade_venda').value.trim(),
            peso_bruto: parseFloat(document.getElementById('prod_peso_bruto').value) || 0,
            peso_liquido: parseFloat(document.getElementById('prod_peso_liquido').value) || 0,
            ncm: document.getElementById('prod_ncm').value.trim(),
            cest: document.getElementById('prod_cest').value.trim() || null,
            id_grupo_tributario: parseInt(document.getElementById('prod_grupo_tributario').value, 10) || 0,
			estoques: composto ? [] : estoquesVinculados,
			composicao: composto ? componentes.map(({ id_produto_componente, quantidade }) => ({ id_produto_componente, quantidade })) : []
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
