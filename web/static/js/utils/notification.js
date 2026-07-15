/**
 * notification.js — Sistema unificado de notificações via Toast
 * 
 * Funções exportadas:
 *   showError(message)   → toast vermelho (erro)
 *   showSuccess(message) → toast verde (sucesso)
 *   showWarning(message) → toast amarelo (aviso)
 *   showInfo(message)    → toast azul (informação)
 */

const DURACAO_MS = 5000;

const CONFIG_TIPO = {
    error: {
        icone: 'bi-x-circle-fill',
        titulo: 'Erro',
        classeToast: 'toast-notif--error',
    },
    success: {
        icone: 'bi-check-circle-fill',
        titulo: 'Sucesso',
        classeToast: 'toast-notif--success',
    },
    warning: {
        icone: 'bi-exclamation-triangle-fill',
        titulo: 'Atenção',
        classeToast: 'toast-notif--warning',
    },
    info: {
        icone: 'bi-info-circle-fill',
        titulo: 'Informação',
        classeToast: 'toast-notif--info',
    },
};

function obterOuCriarContainer() {
    let container = document.getElementById('toast-notif-container');
    if (!container) {
        container = document.createElement('div');
        container.id = 'toast-notif-container';
        document.body.appendChild(container);
    }
    return container;
}

function showNotification(message, tipo = 'error') {
    const config = CONFIG_TIPO[tipo] || CONFIG_TIPO.error;
    const container = obterOuCriarContainer();

    const toastEl = document.createElement('div');
    toastEl.className = `toast-notif ${config.classeToast}`;
    toastEl.setAttribute('role', 'alert');
    toastEl.setAttribute('aria-live', 'assertive');

    toastEl.innerHTML = `
        <div class="toast-notif__header">
            <i class="bi ${config.icone} toast-notif__icon"></i>
            <span class="toast-notif__titulo">${config.titulo}</span>
            <button class="toast-notif__fechar" aria-label="Fechar">
                <i class="bi bi-x-lg"></i>
            </button>
        </div>
        <div class="toast-notif__body">${message}</div>
        <div class="toast-notif__progress"></div>
    `;

    // Fechar ao clicar no X
    toastEl.querySelector('.toast-notif__fechar').addEventListener('click', () => {
        fecharToast(toastEl);
    });

    container.appendChild(toastEl);

    // Disparar animação de entrada no próximo frame
    requestAnimationFrame(() => {
        requestAnimationFrame(() => {
            toastEl.classList.add('toast-notif--visivel');
        });
    });

    // Fechar automaticamente após DURACAO_MS
    const timer = setTimeout(() => fecharToast(toastEl), DURACAO_MS);

    // Pausar temporizador ao passar o mouse
    toastEl.addEventListener('mouseenter', () => clearTimeout(timer));
    toastEl.addEventListener('mouseleave', () => {
        setTimeout(() => fecharToast(toastEl), 1500);
    });
}

function fecharToast(toastEl) {
    toastEl.classList.remove('toast-notif--visivel');
    toastEl.classList.add('toast-notif--saindo');
    toastEl.addEventListener('transitionend', () => toastEl.remove(), { once: true });
}

// ──────────────────────────────────────────────
// Exportações públicas
// ──────────────────────────────────────────────

export function showError(message) {
    showNotification(message, 'error');
}

export function showSuccess(message) {
    showNotification(message, 'success');
}

export function showWarning(message) {
    showNotification(message, 'warning');
}

export function showInfo(message) {
    showNotification(message, 'info');
}
