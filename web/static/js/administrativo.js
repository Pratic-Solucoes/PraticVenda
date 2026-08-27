import { getToken } from './utils/auth.js';

const token = getToken();
const erro = document.getElementById('erro');

function mostrarErro(mensagem) {
    erro.textContent = mensagem;
    erro.classList.remove('d-none');
}

function preencherTabela(id, itens, colunas) {
    const corpo = document.getElementById(id);
    corpo.replaceChildren();
    if (!itens.length) {
        const celula = document.createElement('td');
        celula.colSpan = colunas.length;
        celula.className = 'text-muted';
        celula.textContent = 'Nenhum registro encontrado.';
        const linha = document.createElement('tr');
        linha.append(celula);
        corpo.append(linha);
        return;
    }

    itens.forEach(item => {
        const linha = document.createElement('tr');
        colunas.forEach(coluna => {
            const celula = document.createElement('td');
            celula.textContent = coluna(item);
            linha.append(celula);
        });
        corpo.append(linha);
    });
}

async function carregar(url) {
    const resposta = await fetch(url, { headers: { Authorization: `Bearer ${token}` } });
    if (resposta.status === 401 || resposta.status === 403) {
        localStorage.removeItem('token');
        window.location.href = '/';
        throw new Error('Sessão administrativa inválida.');
    }
    if (!resposta.ok) throw new Error('Não foi possível carregar os dados administrativos.');
    return resposta.json();
}

document.getElementById('sairBtn').addEventListener('click', () => {
    localStorage.removeItem('token');
    window.location.href = '/';
});

if (!token) {
    window.location.href = '/';
} else {
    try {
        const [organizacoes, usuarios] = await Promise.all([
            carregar('/api/administrativo/organizacoes'),
            carregar('/api/administrativo/usuarios'),
        ]);
        const nomesOrganizacoes = new Map(organizacoes.map(org => [String(org.id), org.nome_fantasia]));
        preencherTabela('organizacoes', organizacoes, [
            org => org.nome_fantasia,
            org => org.email,
            org => org.schema,
            org => org.ativo ? 'Ativa' : 'Inativa',
        ]);
        preencherTabela('usuarios', usuarios, [
            usuario => usuario.nome,
            usuario => usuario.email,
            usuario => nomesOrganizacoes.get(String(usuario.id_empresa)) || `#${usuario.id_empresa}`,
            usuario => usuario.ativo ? 'Ativo' : 'Inativo',
        ]);
    } catch (error) {
        if (token) mostrarErro(error.message);
    }
}
