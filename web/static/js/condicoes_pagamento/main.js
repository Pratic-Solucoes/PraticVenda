import { checkAuth, getToken } from '../utils/auth.js';
import { carregarCondicoesPagamento } from './listarCondicoesPagamento.js';
import { setupEditarCondicaoPagamento } from './editarCondicaoPagamento.js';
import { showError } from '../utils/showError.js';

document.addEventListener('DOMContentLoaded', async () => {
    if (!checkAuth()) return;

    const tbody = document.getElementById('tabela_condicoes_pagamento_body');
    const btnNova = document.getElementById('btnNovaCondicaoPagamento');
    const btnFecharForm = document.getElementById('btnFecharFormCondicaoPagamento');
    const btnCancelar = document.getElementById('btnCancelarCondicaoPagamento');
    const collapseFormCP = document.getElementById('collapseFormCondicaoPagamento');
    const formCPTitulo = document.getElementById('formCondicaoPagamentoTitulo');
    const formInlineCP = document.getElementById('formInlineCondicaoPagamento');
    const editCPId = document.getElementById('edit_cp_id');
    const containerFormasPagamento = document.getElementById('edit_cp_formas_pagamento_container');

    let bsCollapse = null;
    if (collapseFormCP) {
        bsCollapse = new bootstrap.Collapse(collapseFormCP, { toggle: false });
    }

    // Carrega Formas de Pagamento para Checkboxes
    async function carregarFormasPagamentoParaCheckboxes() {
        if (!containerFormasPagamento) return;
        try {
            const res = await fetch('/api/formas-pagamento', {
                headers: { 'Authorization': `Bearer ${getToken()}` }
            });
            if (res.ok) {
                const formas = await res.json();
                containerFormasPagamento.innerHTML = '';
                if (formas && formas.length > 0) {
                    formas.forEach(f => {
                        const div = document.createElement('div');
                        div.className = 'form-check';
                        div.innerHTML = `
                            <input class="form-check-input fp-checkbox" type="checkbox" value="${f.id}" id="fp_check_${f.id}">
                            <label class="form-check-label small" for="fp_check_${f.id}">
                                ${f.descricao}
                            </label>
                        `;
                        containerFormasPagamento.appendChild(div);
                    });
                } else {
                    containerFormasPagamento.innerHTML = '<span class="text-muted small">Nenhuma forma encontrada.</span>';
                }
            } else {
                showError("Não foi possível carregar as formas de pagamento no formulário.");
                containerFormasPagamento.innerHTML = '<span class="text-danger small">Erro ao carregar.</span>';
            }
        } catch (err) {
            console.error("Erro ao carregar formas de pagamento", err);
            containerFormasPagamento.innerHTML = '<span class="text-danger small">Erro de comunicação.</span>';
        }
    }

    // Carrega lista inicial
    if (tbody) {
        carregarCondicoesPagamento();
        carregarFormasPagamentoParaCheckboxes();
    }

    // Abre formulário limpo para nova condição de pagamento
    function abrirFormularioNovo() {
        if (formCPTitulo) {
            formCPTitulo.innerHTML = '<i class="bi bi-tags me-2"></i> Nova Condição de Pagamento';
        }
        if (formInlineCP) formInlineCP.reset();
        if (editCPId) editCPId.value = '';
        if (bsCollapse) bsCollapse.show();
    }

    function fecharFormulario() {
        if (bsCollapse) bsCollapse.hide();
        if (formInlineCP) formInlineCP.reset();
        if (editCPId) editCPId.value = '';
    }

    if (btnNova) btnNova.addEventListener('click', abrirFormularioNovo);
    if (btnFecharForm) btnFecharForm.addEventListener('click', fecharFormulario);
    if (btnCancelar) btnCancelar.addEventListener('click', fecharFormulario);

    // Expõe funções globais utilizadas por editarCondicaoPagamento.js
    window.fecharFormularioInlineCondicaoPagamento = fecharFormulario;
    window.abrirFormularioInlineCondicaoPagamento = function () {
        if (bsCollapse) bsCollapse.show();
    };

    setupEditarCondicaoPagamento();
});
