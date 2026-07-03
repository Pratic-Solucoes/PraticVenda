document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');
    const emailInput = document.getElementById('entrada_email');
    const senhaInput = document.getElementById('entrada_senha');

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const email = emailInput.value;
        const senha = senhaInput.value;

        if (!email || !senha) {
            showError("Por favor, preencha email e senha.");
            return;
        }

        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, senha })
            });

            const data = await response.json();

            if (!response.ok) {
                // Em caso de erro, a API retorna mensagem na propriedade correspondente 
                // Exemplo: { "erro": "dados login inválidos" } ou similar
                showError(data.erro || data.mensagem || "Erro ao realizar login");
                return;
            }

            // Sucesso: salvar o token
            if (data.token) {
                localStorage.setItem('token', data.token);
            }

            // Redirecionar para o dashboard
            window.location.href = '/dashboard';

        } catch (error) {
            console.error("Erro no login:", error);
            showError("Erro interno ao comunicar com o servidor.");
        }
    });

    function showError(message) {
        const modalBody = document.getElementById('errorModalBody');
        modalBody.textContent = message;
        const errorModalElement = document.getElementById('errorModal');
        const errorModal = bootstrap.Modal.getOrCreateInstance(errorModalElement);
        errorModal.show();
    }
});
