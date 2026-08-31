import { checkAuth } from '../utils/auth.js';
import { renderizarTabela, atualizarTotalDisplay, configurarEdicaoPreco } from './ui.js';
import { buscarEAdicionarProduto } from './produto.js';
import { removerItem, obterItens, limparCarrinho, definirItens } from './carrinho.js';
import { getToken } from '../utils/auth.js';
import { showError, showSuccess } from '../utils/showError.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

	const inputProduto = document.getElementById('pdv_produto');
	const btnAddProduto = document.getElementById('btn_adicionar_produto');
	const tbody = document.getElementById('tabela_pdv_body');
	const btnAbrirCaixa = document.getElementById('btnAbrirCaixa');
	const btnFecharCaixa = document.getElementById('btnFecharCaixa');
	const selectForma = document.getElementById('pdv_forma_pagamento');
	const selectCondicao = document.getElementById('pdv_condicao_pagamento');
	let configuracaoPDV;
	let formasPagamento = [];
	let condicoesPagamento = [];
	let preVendaAtual = null;
	let emPagamento = false;
	const iniciarPagamento = () => { emPagamento=true; document.getElementById('etapaPagamento').classList.remove('d-none'); document.getElementById('btnFinalizarVenda').innerHTML='<i class="bi bi-check-circle me-1 fs-5"></i> Concluir venda'; };
	const voltarParaProdutos = () => { emPagamento=false; document.getElementById('etapaPagamento').classList.add('d-none'); document.getElementById('btnFinalizarVenda').innerHTML='<i class="bi bi-cash-coin me-1 fs-5"></i> Pagamento'; };

	const headers = () => ({ 'Content-Type':'application/json', Authorization:`Bearer ${getToken()}` });
	async function carregarConfiguracao(){const r=await fetch('/api/configuracoes/pdv',{headers:headers()});if(!r.ok){showError('Configure o PDV antes de realizar vendas.');return false}configuracaoPDV=await r.json();const [formas,condicoes]=await Promise.all([fetch('/api/formas-pagamento',{headers:headers()}).then(r=>r.json()),fetch('/api/condicoes-pagamento',{headers:headers()}).then(r=>r.json())]);formasPagamento=formas.filter(f=>configuracaoPDV.formas_pagamento.includes(Number(f.id)));condicoesPagamento=condicoes.filter(c=>configuracaoPDV.condicoes_pagamento.includes(Number(c.id)));selectCondicao.innerHTML='<option value="">Selecione...</option>';condicoesPagamento.forEach(c=>selectCondicao.innerHTML+=`<option value="${c.id}">${c.descricao}</option>`);document.getElementById('campo_desconto').classList.toggle('d-none',!configuracaoPDV.permitir_desconto_manual);configurarEdicaoPreco(configuracaoPDV.permitir_alterar_preco);renderizarTabela();return true}
	function carregarFormas(){const condicao=condicoesPagamento.find(c=>Number(c.id)===Number(selectCondicao.value));selectForma.innerHTML='<option value="">Selecione...</option>';if(!condicao)return;formasPagamento.filter(f=>condicao.formas_pagamento.includes(Number(f.id))).forEach(f=>selectForma.innerHTML+=`<option value="${f.id}">${f.descricao}</option>`)}
	async function carregarClientes(){const r=await fetch('/api/clientes',{headers:{Authorization:`Bearer ${getToken()}`}});if(!r.ok)return;const s=document.getElementById('pdv_cliente');s.innerHTML='<option value="">Consumidor final</option>';(await r.json()).forEach(c=>s.innerHTML+=`<option value="${c.id}">${c.nome}</option>`)}
	async function carregarCaixas(){const r=await fetch('/api/caixas',{headers:{Authorization:`Bearer ${getToken()}`}});const caixas=await r.json();const s=document.getElementById('abertura_caixa');s.innerHTML=caixas.length?'':'<option value="">Nenhum caixa criado</option>';caixas.forEach(c=>s.innerHTML+=`<option value="${c.id}">${c.nome}</option>`)}
	function bloquearPdv(bloquear){document.querySelectorAll('#pdv_operacao input, #pdv_operacao select, #pdv_operacao button').forEach(el=>el.disabled=bloquear);document.getElementById('pdv_operacao').classList.toggle('opacity-50',bloquear)}
	function atualizarAcoesCaixa(aberto){btnAbrirCaixa.classList.toggle('d-none',aberto);btnFecharCaixa.classList.toggle('d-none',!aberto)}
	async function verificarCaixa(){bloquearPdv(true);const r=await fetch('/api/caixas/atual',{headers:{Authorization:`Bearer ${getToken()}`}});if(!r.ok)return showError('Não foi possível verificar o caixa.');const d=await r.json();atualizarAcoesCaixa(d.aberto);if(!d.aberto){await carregarCaixas();return}bloquearPdv(false)}
	document.getElementById('btnCriarCaixa').onclick=async()=>{const r=await fetch('/api/caixas',{method:'POST',headers:headers(),body:JSON.stringify({nome:'Caixa principal'})});if(!r.ok)return showError('Não foi possível criar o caixa.');await carregarCaixas()};
	btnAbrirCaixa.onclick=async()=>{await carregarCaixas();new bootstrap.Modal(document.getElementById('modalAbrirCaixa')).show()};
	document.getElementById('formAbrirCaixa').onsubmit=async e=>{e.preventDefault();const r=await fetch('/api/caixas/abrir',{method:'POST',headers:headers(),body:JSON.stringify({id_caixa:+abertura_caixa.value,valor_abertura:+abertura_valor.value,senha_acesso:abertura_senha.value})});if(!r.ok){const d=await r.json();return showError(d.erro||'Não foi possível abrir o caixa.')}e.target.reset();bootstrap.Modal.getInstance(document.getElementById('modalAbrirCaixa')).hide();bloquearPdv(false);atualizarAcoesCaixa(true);showSuccess('Caixa aberto com sucesso.')};
	btnFecharCaixa.onclick=()=>new bootstrap.Modal(document.getElementById('modalFecharCaixa')).show();
	document.getElementById('formFecharCaixa').onsubmit=async e=>{e.preventDefault();const r=await fetch('/api/caixas/fechar',{method:'POST',headers:headers(),body:JSON.stringify({valor_dinheiro:+fechamento_dinheiro.value,valor_cartao:+fechamento_cartao.value,senha_acesso:fechamento_senha.value})});if(!r.ok){const d=await r.json();return showError(d.erro||'Não foi possível fechar o caixa.')}e.target.reset();bootstrap.Modal.getInstance(document.getElementById('modalFecharCaixa')).hide();showSuccess('Caixa fechado.');verificarCaixa()};
    const montarPayload = () => {
		const cliente=+document.getElementById('pdv_cliente').value;
		const desconto=+document.getElementById('pdv_desconto').value||0;
		return { ...(preVendaAtual ? {id: preVendaAtual} : {}), valor_desconto:desconto, apelido_consumidor:document.getElementById('pdv_apelido_consumidor').value.trim(), ...(cliente ? {id_cliente: cliente} : {}), itens: obterItens().map(item => ({ id_produto: (item.produto.produto || item.produto).id, quantidade: item.quantidade, ...(configuracaoPDV.permitir_alterar_preco ? {valor_unitario:item.precoUnitario} : {}) })) };
	};
    document.getElementById('btnFinalizarVenda')?.addEventListener('click', async () => {
        const itens = obterItens();
        if (!itens.length) return showError('Adicione produtos antes de finalizar a venda.');
		if (!emPagamento) return iniciarPagamento();
		const condicao=+selectCondicao.value;if(!condicao)return showError('Selecione a condição de pagamento.');
		const forma=+selectForma.value;if(!forma)return showError('Selecione a forma de pagamento.');
		const payload = {...montarPayload(), id_condicao_pagamento:condicao, id_forma_pagamento: forma};
        const res = await fetch('/api/pdv/vendas', { method: 'POST', headers: { 'Content-Type':'application/json', Authorization:`Bearer ${getToken()}` }, body: JSON.stringify(payload) });
        const data = await res.json();
        if (!res.ok) return showError(data.erro || 'Não foi possível concluir a venda.');
        limparCarrinho(); preVendaAtual=null; voltarParaProdutos(); document.getElementById('pdv_apelido_consumidor').value=''; renderizarTabela(); atualizarTotalDisplay(); showSuccess(data.mensagem);
    });
	async function salvarPreVenda(){ if(!obterItens().length)return showError('Adicione produtos antes de salvar a pré-venda.'); const res=await fetch('/api/pdv/pre-vendas',{method:'POST',headers:headers(),body:JSON.stringify(montarPayload())});const d=await res.json();if(!res.ok)return showError(d.erro||'Não foi possível salvar a pré-venda.');limparCarrinho();preVendaAtual=null;document.getElementById('pdv_apelido_consumidor').value='';renderizarTabela();atualizarTotalDisplay();showSuccess(`Pré-venda #${d.id} salva.`)}
	document.getElementById('btnSalvarPreVenda')?.addEventListener('click',salvarPreVenda);
	document.getElementById('btnCancelarVenda')?.addEventListener('click',async()=>{if(preVendaAtual){const res=await fetch(`/api/pdv/vendas/${preVendaAtual}/cancelar`,{method:'POST',headers:headers()});const d=await res.json();if(!res.ok)return showError(d.erro||'Não foi possível cancelar.');showSuccess(d.mensagem)}limparCarrinho();preVendaAtual=null;voltarParaProdutos();document.getElementById('pdv_apelido_consumidor').value='';renderizarTabela();atualizarTotalDisplay()});
	async function consultarPreVendas(){const res=await fetch('/api/pdv/pre-vendas',{headers:headers()});const lista=await res.json();if(!res.ok)return showError(lista.erro||'Não foi possível consultar pré-vendas.');const body=document.getElementById('preVendasBody');body.innerHTML=lista.length?lista.map(v=>`<tr><td>#${v.id}</td><td>${v.cliente||v.apelido_consumidor||'Consumidor final'}</td><td>R$ ${Number(v.valor_total).toFixed(2).replace('.',',')}</td><td><button class="btn btn-sm btn-primary" data-id="${v.id}">Selecionar</button></td></tr>`).join(''):'<tr><td colspan="4" class="text-muted">Nenhuma pré-venda aberta.</td></tr>';new bootstrap.Modal(document.getElementById('modalPreVendas')).show()}
	document.getElementById('btnConsultarPreVendas')?.addEventListener('click',consultarPreVendas);
	document.getElementById('preVendasBody')?.addEventListener('click',async e=>{const b=e.target.closest('button[data-id]');if(!b)return;const res=await fetch(`/api/pdv/pre-vendas/${b.dataset.id}`,{headers:headers()});const v=await res.json();if(!res.ok)return showError(v.erro||'Não foi possível abrir a pré-venda.');const produtos=await Promise.all(v.itens.map(i=>fetch(`/api/produtos/${i.id_produto}`,{headers:headers()}).then(r=>r.json()).then(p=>({...p,quantidade:i.quantidade,valor_unitario:i.valor_unitario}))));definirItens(produtos);preVendaAtual=v.id;document.getElementById('pdv_cliente').value=v.id_cliente||'';document.getElementById('pdv_apelido_consumidor').value=v.apelido_consumidor||'';bootstrap.Modal.getInstance(document.getElementById('modalPreVendas')).hide();renderizarTabela();atualizarTotalDisplay();showSuccess(`Pré-venda #${v.id} carregada.`)});

	// 1. Inicializar UI limpa
    renderizarTabela();
    atualizarTotalDisplay();
	bloquearPdv(true);
	carregarConfiguracao().then(configurado=>{if(configurado){carregarClientes();verificarCaixa()}}); selectCondicao.addEventListener('change',carregarFormas);

    // 2. Handler para buscar o produto
    async function handleAddProduto() {
        const query = inputProduto.value.trim();
        if (query) {
            await buscarEAdicionarProduto(query);
            inputProduto.value = '';
            inputProduto.focus();
        }
    }

    // 3. Eventos do Input de Busca de Produto
    if (btnAddProduto) {
        btnAddProduto.addEventListener('click', handleAddProduto);
    }

    let debounceTimeout;
    if (inputProduto) {
        inputProduto.addEventListener('input', (e) => {
            clearTimeout(debounceTimeout);
            debounceTimeout = setTimeout(() => {
                const query = inputProduto.value.trim();
                buscarEAdicionarProduto(query);
            }, 400); // 400ms debounce
        });

        inputProduto.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                clearTimeout(debounceTimeout);
                handleAddProduto();
            }
        });
    }

    // Fechar dropdown ao clicar fora
    document.addEventListener('click', (e) => {
        const dropdown = document.getElementById('dropdown_produto');
        if (dropdown && !dropdown.contains(e.target) && e.target !== inputProduto) {
            dropdown.classList.add('d-none');
        }
    });

    // 4. Delegação de Eventos para os botões dentro da Tabela
    if (tbody) {
        // Remover Item
        tbody.addEventListener('click', (e) => {
            const btnRemover = e.target.closest('.btn-remover');
            if (btnRemover) {
                const index = parseInt(btnRemover.getAttribute('data-index'), 10);
                removerItem(index);
                renderizarTabela();
                atualizarTotalDisplay();
            }
        });

        // Alterar Quantidade Manulamente
        tbody.addEventListener('change', (e) => {
            if (e.target.classList.contains('input-qtd')) {
                const index = parseInt(e.target.getAttribute('data-index'), 10);
                const novaQtd = parseInt(e.target.value, 10);
                
                if (novaQtd > 0) {
                    const itens = obterItens();
                    const item = itens[index];
                    
					item.quantidade = novaQtd;
					item.subtotal = novaQtd * item.precoUnitario;
                    
                    renderizarTabela();
                    atualizarTotalDisplay();
                } else {
                    e.target.value = 1; // Reseta se colocar número inválido
				}
			}
			if (e.target.classList.contains('input-preco')) {
				const index = parseInt(e.target.getAttribute('data-index'), 10);
				const novoPreco = Number(e.target.value);
				if (novoPreco < 0 || Number.isNaN(novoPreco)) return;
				const item = obterItens()[index];
				item.precoUnitario = novoPreco;
				item.subtotal = item.quantidade * novoPreco;
				renderizarTabela(); atualizarTotalDisplay();
			}
        });
    }
});
