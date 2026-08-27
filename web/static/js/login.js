document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');
    const emailInput = document.getElementById('entrada_email');
    const senhaInput = document.getElementById('entrada_senha');
    const btnToggleSenha = document.getElementById('toggleSenhaBtn');
    const iconToggleSenha = document.getElementById('toggleSenhaIcon');
    const feedback = document.getElementById('loginFeedback');

    function mostrarErro(mensagem) {
        feedback.textContent = mensagem;
        feedback.classList.remove('d-none');
    }

    function ocultarErro() {
        feedback.textContent = '';
        feedback.classList.add('d-none');
    }

    if (btnToggleSenha) {
        btnToggleSenha.addEventListener('click', () => {
            if (senhaInput.type === 'password') {
                senhaInput.type = 'text';
                iconToggleSenha.classList.remove('bi-eye');
                iconToggleSenha.classList.add('bi-eye-slash');
            } else {
                senhaInput.type = 'password';
                iconToggleSenha.classList.remove('bi-eye-slash');
                iconToggleSenha.classList.add('bi-eye');
            }
        });
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const loginAdministrativo = e.submitter?.dataset.login === 'administrativo';

        const email = emailInput.value;
        const senha = senhaInput.value;

        if (!email || !senha) {
            mostrarErro('Informe seu e-mail e sua senha para continuar.');
            return;
        }

        ocultarErro();

        try {
            const response = await fetch(loginAdministrativo ? '/api/login/administrativo' : '/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, senha })
            });

            const tipoConteudo = response.headers.get('content-type') || '';
            const data = tipoConteudo.includes('application/json')
                ? await response.json()
                : await response.text();

            if (!response.ok) {
                const mensagemApi = typeof data === 'string'
                    ? data
                    : data.erro || data.mensagem;
                mostrarErro(mensagemApi || (loginAdministrativo
                    ? 'Não foi possível entrar como administrador. Verifique e-mail e senha.'
                    : 'E-mail ou senha inválidos. Se esta for uma conta administrativa, use “Entrar como admin”.'));
                return;
            }

            // Sucesso: salvar o token
            if (data.token) {
                localStorage.setItem('token', data.token);
            }

            window.location.href = loginAdministrativo ? '/administrativo' : '/dashboard';

        } catch (error) {
            console.error("Erro no login:", error);
            mostrarErro('Não foi possível comunicar com o servidor. Tente novamente em instantes.');
        }
    });
});
