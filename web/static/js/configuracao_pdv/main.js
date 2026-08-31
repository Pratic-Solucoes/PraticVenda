import { checkAuth, getToken } from '../utils/auth.js';
import { showError, showSuccess } from '../utils/showError.js';

const headers = () => ({ Authorization: `Bearer ${getToken()}` });
const option = (item, texto) => `<option value="${item.id}">${texto(item)}</option>`;
const checks = (container, itens, prefixo, texto, selecionados = []) => {
    container.innerHTML = itens.length ? itens.map(item => `<div class="form-check"><input class="form-check-input" type="checkbox" value="${item.id}" id="${prefixo}_${item.id}" ${selecionados.includes(Number(item.id)) ? 'checked' : ''}><label class="form-check-label" for="${prefixo}_${item.id}">${texto(item)}</label></div>`).join('') : '<span class="text-muted">Nenhum item cadastrado.</span>';
};
const selecionados = container => [...container.querySelectorAll('input:checked')].map(input => Number(input.value));

document.addEventListener('DOMContentLoaded', async () => {
    if (!checkAuth()) return;
    const estoque = document.getElementById('pdv_estoque');
    const categoria = document.getElementById('pdv_categoria_credito');
    const formas = document.getElementById('pdv_formas_pagamento');
    const condicoes = document.getElementById('pdv_condicoes_pagamento');
    let configuracao = null;
    try {
        const [estoques, categorias, listaFormas, listaCondicoes, respostaConfiguracao] = await Promise.all([
            fetch('/api/estoques', { headers: headers() }).then(r => r.json()),
            fetch('/api/contas-receber/categorias', { headers: headers() }).then(r => r.json()),
            fetch('/api/formas-pagamento', { headers: headers() }).then(r => r.json()),
            fetch('/api/condicoes-pagamento', { headers: headers() }).then(r => r.json()),
            fetch('/api/configuracoes/pdv', { headers: headers() }),
        ]);
        if (respostaConfiguracao.ok) configuracao = await respostaConfiguracao.json();
        estoque.innerHTML = '<option value="">Selecione...</option>' + estoques.map(item => option(item, x => x.nome)).join('');
        categoria.innerHTML = '<option value="">Selecione...</option>' + categorias.map(item => option(item, x => x.nome)).join('');
        checks(formas, listaFormas, 'forma', x => `${x.descricao} (${x.tipo})`, configuracao?.formas_pagamento || []);
        checks(condicoes, listaCondicoes, 'condicao', x => `${x.descricao} — ${x.qtd_parcelas}x`, configuracao?.condicoes_pagamento || []);
        if (configuracao) {
            estoque.value = configuracao.id_estoque_padrao;
            categoria.value = configuracao.id_categoria_credito;
            document.getElementById('pdv_exigir_cliente_prazo').checked = configuracao.exigir_cliente_prazo;
            document.getElementById('pdv_permitir_desconto').checked = configuracao.permitir_desconto_manual;
            document.getElementById('pdv_permitir_alterar_preco').checked = configuracao.permitir_alterar_preco;
			document.getElementById('pdv_gerar_financeiro_recebido').checked = configuracao.gerar_financeiro_recebido;
        } else { document.getElementById('pdv_exigir_cliente_prazo').checked = true; }
    } catch (error) { console.error(error); showError('Não foi possível carregar as configurações de PDV.'); }
    document.getElementById('formConfiguracaoPDV').addEventListener('submit', async event => {
        event.preventDefault();
        const payload = { id_estoque_padrao: Number(estoque.value), id_categoria_credito: Number(categoria.value), formas_pagamento: selecionados(formas), condicoes_pagamento: selecionados(condicoes), exigir_cliente_prazo: document.getElementById('pdv_exigir_cliente_prazo').checked, permitir_desconto_manual: document.getElementById('pdv_permitir_desconto').checked, permitir_alterar_preco: document.getElementById('pdv_permitir_alterar_preco').checked, gerar_financeiro_recebido: document.getElementById('pdv_gerar_financeiro_recebido').checked };
        const response = await fetch('/api/configuracoes/pdv', { method: 'PUT', headers: { ...headers(), 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
        if (!response.ok) { const data = await response.json(); showError(data.erro || 'Não foi possível salvar as configurações.'); return; }
        showSuccess('Configurações de PDV salvas com sucesso.');
    });
});
