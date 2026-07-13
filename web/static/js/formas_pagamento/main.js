import { checkAuth } from '../utils/auth.js';
import { carregarFormasPagamento } from './listarFormasPagamento.js';
import { setupCriarFormaPagamento } from './criarFormaPagamento.js';
import { setupEditarFormaPagamento } from './editarFormaPagamento.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    const tbody                    = document.getElementById('tabela_formas_pagamento_body');
    const btnNova                  = document.getElementById('btnNovaFormaPagamento');
    const btnFecharForm            = document.getElementById('btnFecharFormFormaPagamento');
    const btnCancelar              = document.getElementById('btnCancelarFormaPagamento');
    const collapseFormFP           = document.getElementById('collapseFormFormaPagamento');
    const formFPTitulo             = document.getElementById('formFormaPagamentoTitulo');
    const formInlineFP             = document.getElementById('formInlineFormaPagamento');
    const editFPId                 = document.getElementById('edit_fp_id');

    let bsCollapse = null;
    if (collapseFormFP) {
        bsCollapse = new bootstrap.Collapse(collapseFormFP, { toggle: false });
    }

    // Carrega lista inicial
    if (tbody) {
        carregarFormasPagamento();
    }

    // Abre formulário limpo para nova forma de pagamento
    function abrirFormularioNovo() {
        if (formFPTitulo) {
            formFPTitulo.innerHTML = '<i class="bi bi-credit-card me-2"></i> Nova Forma de Pagamento';
        }
        if (formInlineFP)  formInlineFP.reset();
        if (editFPId)      editFPId.value = '';
        if (bsCollapse)    bsCollapse.show();
    }

    function fecharFormulario() {
        if (bsCollapse) bsCollapse.hide();
        if (formInlineFP) formInlineFP.reset();
        if (editFPId)     editFPId.value = '';
    }

    if (btnNova)       btnNova.addEventListener('click', abrirFormularioNovo);
    if (btnFecharForm) btnFecharForm.addEventListener('click', fecharFormulario);
    if (btnCancelar)   btnCancelar.addEventListener('click', fecharFormulario);

    // Expõe funções globais utilizadas por editarFormaPagamento.js
    window.fecharFormularioInlineFormaPagamento = fecharFormulario;
    window.abrirFormularioInlineFormaPagamento  = function() {
        if (bsCollapse) bsCollapse.show();
    };

    setupEditarFormaPagamento();
    setupCriarFormaPagamento();
});
