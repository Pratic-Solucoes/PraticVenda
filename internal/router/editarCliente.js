import { getToken } from '../utils/auth.js';
import { showError } from '../utils/showError.js';

// Assume carregarClientes is defined elsewhere and reloads the main list
import { carregarClientes } from './listarClientes.js';

let clienteIdAtual;
let modalEditar;
let modalEndereco;
let modalTelefone;

async function fetchAPI(url, options) {
    const token = getToken();
    options.headers = {
        ...options.headers,
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    };
    const res = await fetch(url, options);
    if (!res.ok) {
        const data = await res.json();
        throw new Error(data.erro || `Erro na requisição: ${res.statusText}`);
    }
    if (res.status !== 204 && res.headers.get('Content-Type')?.includes('application/json')) {
        return res.json();
    }
    return null;
}

function popularDadosCadastrais(cliente) {
    document.getElementById('edit_cliente_nome').value = cliente.nome || '';
    document.getElementById('edit_cliente_email').value = cliente.email || '';
    document.getElementById('edit_cliente_tipo').value = cliente.tipo || 'PF';
    document.getElementById('edit_cliente_ie').value = cliente.ie || '';
    document.getElementById('edit_cliente_contribuinte').value = cliente.contribuinte || 9;
    document.getElementById('edit_cliente_is_consumidor_final').checked = cliente.is_consumidor_final || false;

    const cpfCnpjInput = document.getElementById('edit_cliente_cpf_cnpj');
    if (cliente.tipo === 'PJ') {
        cpfCnpjInput.value = cliente.cnpj || '';
    } else {
        cpfCnpjInput.value = cliente.cpf || '';
    }
}

function renderizarEnderecos(enderecos = []) {
    const container = document.getElementById('containerEnderecosCliente');
    const template = document.getElementById('template-endereco-item');
    container.innerHTML = '';

    if (enderecos.length === 0) {
        container.innerHTML = '<p class="text-muted text-center">Nenhum endereço cadastrado.</p>';
        return;
    }

    enderecos.forEach(end => {
        const clone = template.content.cloneNode(true);
        clone.querySelector('.logradouro').textContent = end.logradouro;
        clone.querySelector('.numero').textContent = end.numero;
        clone.querySelector('.bairro').textContent = end.bairro;
        clone.querySelector('.municipio').textContent = end.municipio;
        clone.querySelector('.uf').textContent = end.uf;
        clone.querySelector('.cep').textContent = end.cep;
        clone.querySelector('.btn-editar-endereco').addEventListener('click', () => abrirModalEndereco(end));
        container.appendChild(clone);
    });
}

function renderizarTelefones(telefones = []) {
    const container = document.getElementById('containerTelefonesCliente');
    const template = document.getElementById('template-telefone-item');
    container.innerHTML = '';

    if (telefones.length === 0) {
        container.innerHTML = '<p class="text-muted text-center">Nenhum telefone cadastrado.</p>';
        return;
    }

    telefones.forEach(tel => {
        const clone = template.content.cloneNode(true);
        clone.querySelector('.ddd').textContent = tel.ddd;
        clone.querySelector('.numero').textContent = tel.numero;
        clone.querySelector('.btn-editar-telefone').addEventListener('click', () => abrirModalTelefone(tel));
        container.appendChild(clone);
    });
}

async function carregarDadosCliente(id) {
    try {
        const cliente = await fetchAPI(`/api/clientes/${id}`, { method: 'GET' });
        popularDadosCadastrais(cliente);
        renderizarEnderecos(cliente.enderecos);
        renderizarTelefones(cliente.telefones); // Assumindo que o backend retorna telefones
    } catch (err) {
        showError(err.message);
        modalEditar.hide();
    }
}

function abrirModalEndereco(endereco = {}) {
    const form = document.getElementById('formEnderecoCliente');
    form.reset();
    document.getElementById('endereco_id').value = endereco.id || '';
    document.getElementById('end_cep').value = endereco.cep || '';
    document.getElementById('end_logradouro').value = endereco.logradouro || '';
    document.getElementById('end_numero').value = endereco.numero || '';
    document.getElementById('end_bairro').value = endereco.bairro || '';
    document.getElementById('end_municipio').value = endereco.municipio || '';
    document.getElementById('end_uf').value = endereco.uf || '';
    document.getElementById('end_codigo_municipio').value = endereco.codigo_municipio || '';
    modalEndereco.show();
}

function abrirModalTelefone(telefone = {}) {
    const form = document.getElementById('formTelefoneCliente');
    form.reset();
    document.getElementById('telefone_id').value = telefone.id || '';
    document.getElementById('tel_ddd').value = telefone.ddd || '';
    document.getElementById('tel_numero').value = telefone.numero || '';
    modalTelefone.show();
}

async function salvarDadosCadastrais(e) {
    e.preventDefault();
    const tipo = document.getElementById('edit_cliente_tipo').value;
    const cpfCnpj = document.getElementById('edit_cliente_cpf_cnpj').value;

    const payload = {
        nome: document.getElementById('edit_cliente_nome').value,
        email: document.getElementById('edit_cliente_email').value,
        tipo: tipo,
        cpf: tipo === 'PF' ? cpfCnpj : '',
        cnpj: tipo === 'PJ' ? cpfCnpj : '',
        ie: document.getElementById('edit_cliente_ie').value,
        contribuinte: parseInt(document.getElementById('edit_cliente_contribuinte').value),
        is_consumidor_final: document.getElementById('edit_cliente_is_consumidor_final').checked,
    };

    try {
        await fetchAPI(`/api/clientes/${clienteIdAtual}`, {
            method: 'PUT',
            body: JSON.stringify(payload)
        });
        // Apenas fecha o modal principal se a aba principal for salva.
        // O ideal é apenas mostrar uma mensagem de sucesso.
        alert('Dados do cliente atualizados com sucesso!');
        carregarClientes(); // Recarrega a lista principal
    } catch (err) {
        showError(err.message);
    }
}

async function salvarEndereco(e) {
    e.preventDefault();
    const idEndereco = document.getElementById('endereco_id').value;
    const payload = {
        cep: document.getElementById('end_cep').value,
        logradouro: document.getElementById('end_logradouro').value,
        numero: document.getElementById('end_numero').value,
        bairro: document.getElementById('end_bairro').value,
        municipio: document.getElementById('end_municipio').value,
        uf: document.getElementById('end_uf').value,
        codigo_municipio: document.getElementById('end_codigo_municipio').value,
    };

    const isEdit = !!idEndereco;
    const url = isEdit ? `/api/clientes/${clienteIdAtual}/enderecos/${idEndereco}` : `/api/clientes/${clienteIdAtual}/enderecos`;
    const method = isEdit ? 'PUT' : 'POST';

    try {
        await fetchAPI(url, { method, body: JSON.stringify(payload) });
        modalEndereco.hide();
        carregarDadosCliente(clienteIdAtual); // Recarrega os dados na aba
    } catch (err) {
        showError(err.message);
    }
}

async function salvarTelefone(e) {
    e.preventDefault();
    const idTelefone = document.getElementById('telefone_id').value;
    const payload = {
        ddd: document.getElementById('tel_ddd').value,
        numero: document.getElementById('tel_numero').value,
    };

    const isEdit = !!idTelefone;
    const url = isEdit ? `/api/clientes/${clienteIdAtual}/telefones/${idTelefone}` : `/api/clientes/${clienteIdAtual}/telefones`;
    const method = isEdit ? 'PUT' : 'POST';

    try {
        await fetchAPI(url, { method, body: JSON.stringify(payload) });
        modalTelefone.hide();
        carregarDadosCliente(clienteIdAtual); // Recarrega os dados na aba
    } catch (err) {
        showError(err.message);
    }
}


export function setupEditarCliente() {
    const modalEl = document.getElementById('modalEditarCliente');
    if (!modalEl) return;

    modalEditar = new bootstrap.Modal(modalEl);
    modalEndereco = new bootstrap.Modal(document.getElementById('modalEnderecoCliente'));
    modalTelefone = new bootstrap.Modal(document.getElementById('modalTelefoneCliente'));

    // Função global para abrir o modal de edição
    window.abrirModalEditarCliente = (id) => {
        clienteIdAtual = id;
        document.getElementById('edit_cliente_id').value = id;
        
        // Reseta para a primeira aba ao abrir
        const firstTab = new bootstrap.Tab(document.getElementById('dados-cadastrais-tab'));
        firstTab.show();

        carregarDadosCliente(id);
        modalEditar.show();
    };

    // Listeners para os formulários de salvar
    document.getElementById('formEditarCliente').addEventListener('submit', salvarDadosCadastrais);
    document.getElementById('formEnderecoCliente').addEventListener('submit', salvarEndereco);
    document.getElementById('formTelefoneCliente').addEventListener('submit', salvarTelefone);

    // Listeners para os botões de "Novo"
    document.getElementById('btnNovoEndereco').addEventListener('click', () => abrirModalEndereco());
    document.getElementById('btnNovoTelefone').addEventListener('click', () => abrirModalTelefone());

    // Lógica para alternar campo CPF/CNPJ
    document.getElementById('edit_cliente_tipo').addEventListener('change', (e) => {
        const label = document.querySelector('label[for="edit_cliente_cpf_cnpj"]');
        if (e.target.value === 'PJ') {
            label.textContent = 'CNPJ';
        } else {
            label.textContent = 'CPF';
        }
    });
}

// Adiciona um listener para garantir que o DOM está carregado
document.addEventListener('DOMContentLoaded', () => {
    // Tenta configurar o listener. Se o modal não estiver na página, não faz nada.
    // Isso é útil para SPA onde o conteúdo é carregado dinamicamente.
    if (document.getElementById('modalEditarCliente')) {
        setupEditarCliente();
    }
});

// O arquivo principal (ex: /js/clientes/main.js) deve importar e chamar setupEditarCliente()
// Exemplo de main.js:
// import { setupEditarCliente } from './editarCliente.js';
// import { carregarClientes } from './listarClientes.js';
//
// document.addEventListener('DOMContentLoaded', () => {
//     carregarClientes();
//     setupEditarCliente();
// });