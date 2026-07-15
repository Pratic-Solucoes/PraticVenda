import { getToken } from "../utils/auth.js";
import { carregarClientes } from "./listarClientes.js";
import { showError, showSuccess } from "../utils/showError.js";
import { validaRespostaRequisicao } from "../utils/resposta.js";
import { fecharModal } from "../utils/fecharModal.js";

let clienteIdAtual = null;
let modalEndereco = null;
let modalTelefone = null;

async function requisicaoAPI(url, options = {}) {
  const token = getToken();
  const res = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...options.headers,
    },
  });

  return await validaRespostaRequisicao(res);
}

async function buscarClienteAPI(id) {
  return requisicaoAPI(`/api/clientes/${id}`, { method: "GET" });
}

async function criarClienteAPI(payload) {
  return requisicaoAPI(`/api/clientes`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

async function atualizarClienteAPI(id, payload) {
  return requisicaoAPI(`/api/clientes/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

async function criarEnderecoAPI(idCliente, payload) {
  return requisicaoAPI(`/api/clientes/${idCliente}/enderecos`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

async function atualizarEnderecoAPI(idCliente, idEndereco, payload) {
  return requisicaoAPI(`/api/clientes/${idCliente}/enderecos/${idEndereco}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

async function criarTelefoneAPI(idCliente, payload) {
  return requisicaoAPI(`/api/clientes/${idCliente}/telefones`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

async function atualizarTelefoneAPI(idCliente, idTelefone, payload) {
  return requisicaoAPI(`/api/clientes/${idCliente}/telefones/${idTelefone}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

function popularFormularioCliente(cliente) {
  document.getElementById("edit_cliente_id").value = cliente.id;
  document.getElementById("edit_cliente_nome").value = cliente.nome || "";
  document.getElementById("edit_cliente_email").value = cliente.email || "";
  document.getElementById("edit_cliente_tipo").value = cliente.tipo || "PF";
  document.getElementById("edit_cliente_ie").value = cliente.ie || "";
  document.getElementById("edit_cliente_contribuinte").value = cliente.contribuinte || 9;
  document.getElementById("edit_cliente_is_consumidor_final").checked =
    cliente.is_consumidor_final || false;

  const cpfCnpjInput = document.getElementById("edit_cliente_cpf_cnpj");
  const ieInput = document.getElementById("edit_cliente_ie");

  if (cliente.tipo === "PJ") {
    cpfCnpjInput.value = cliente.cnpj || "";
    ieInput.disabled = false;
  } else {
    cpfCnpjInput.value = cliente.cpf || "";
    ieInput.value = "";
    ieInput.disabled = true;
  }
}

function renderizarEnderecos(enderecos = []) {
  const container = document.getElementById("containerEnderecosCliente");
  const template = document.getElementById("template-endereco-item");
  if (!container || !template) return;

  container.innerHTML = "";

  if (enderecos.length === 0) {
    container.innerHTML =
      '<p class="text-muted text-center py-3">Nenhum endereço cadastrado.</p>';
    return;
  }

  enderecos.forEach((end) => {
    const clone = template.content.cloneNode(true);
    clone.querySelector(".logradouro").textContent = end.logradouro;
    clone.querySelector(".numero").textContent = end.numero;
    clone.querySelector(".bairro").textContent = end.bairro;
    clone.querySelector(".municipio").textContent = end.municipio;
    clone.querySelector(".uf").textContent = end.uf;
    clone.querySelector(".cep").textContent = end.cep;
    clone
      .querySelector(".btn-editar-endereco")
      .addEventListener("click", () => abrirModalEndereco(end));
    container.appendChild(clone);
  });
}

function renderizarTelefones(telefones = []) {
  const container = document.getElementById("containerTelefonesCliente");
  const template = document.getElementById("template-telefone-item");
  if (!container || !template) return;

  container.innerHTML = "";

  if (telefones.length === 0) {
    container.innerHTML =
      '<p class="text-muted text-center py-3">Nenhum telefone cadastrado.</p>';
    return;
  }

  telefones.forEach((tel) => {
    const clone = template.content.cloneNode(true);
    clone.querySelector(".ddd").textContent = tel.ddd;
    clone.querySelector(".numero").textContent = tel.numero;
    clone
      .querySelector(".btn-editar-telefone")
      .addEventListener("click", () => abrirModalTelefone(tel));
    container.appendChild(clone);
  });
}

async function carregarDadosCliente(id) {
  const cliente = await buscarClienteAPI(id);
  popularFormularioCliente(cliente);
  renderizarEnderecos(cliente.enderecos || []);
  renderizarTelefones(cliente.telefones || []);
}

function abrirModalEndereco(endereco = {}) {
  const form = document.getElementById("formEnderecoCliente");
  form.reset();
  document.getElementById("endereco_id").value = endereco.id || "";
  document.getElementById("end_cep").value = endereco.cep || "";
  document.getElementById("end_logradouro").value = endereco.logradouro || "";
  document.getElementById("end_numero").value = endereco.numero || "";
  document.getElementById("end_bairro").value = endereco.bairro || "";
  document.getElementById("end_municipio").value = endereco.municipio || "";
  document.getElementById("end_uf").value = endereco.uf || "";
  document.getElementById("end_codigo_municipio").value = endereco.codigo_municipio || "";
  modalEndereco.show();
}

function abrirModalTelefone(telefone = {}) {
  const form = document.getElementById("formTelefoneCliente");
  form.reset();
  document.getElementById("telefone_id").value = telefone.id || "";
  document.getElementById("tel_ddd").value = telefone.ddd || "";
  document.getElementById("tel_numero").value = telefone.numero || "";
  modalTelefone.show();
}

function montarPayloadCliente() {
  const tipo = document.getElementById("edit_cliente_tipo").value;
  const cpfCnpj = document.getElementById("edit_cliente_cpf_cnpj").value;

  return {
    nome: document.getElementById("edit_cliente_nome").value,
    email: document.getElementById("edit_cliente_email").value,
    tipo,
    cpf: tipo === "PF" ? cpfCnpj : "",
    cnpj: tipo === "PJ" ? cpfCnpj : "",
    ie: document.getElementById("edit_cliente_ie").value,
    contribuinte: parseInt(document.getElementById("edit_cliente_contribuinte").value, 10),
    is_consumidor_final: document.getElementById("edit_cliente_is_consumidor_final").checked,
  };
}

function montarPayloadEndereco() {
  return {
    cep: document.getElementById("end_cep").value,
    logradouro: document.getElementById("end_logradouro").value,
    numero: document.getElementById("end_numero").value,
    bairro: document.getElementById("end_bairro").value,
    municipio: document.getElementById("end_municipio").value,
    uf: document.getElementById("end_uf").value,
    codigo_municipio: document.getElementById("end_codigo_municipio").value,
  };
}

function montarPayloadTelefone() {
  return {
    ddd: document.getElementById("tel_ddd").value,
    numero: document.getElementById("tel_numero").value,
  };
}

async function salvarDadosCadastrais(e) {
  e.preventDefault();

  try {
    if (clienteIdAtual) {
      await atualizarClienteAPI(clienteIdAtual, montarPayloadCliente());
      showSuccess("Cliente atualizado com sucesso!");
    } else {
      await criarClienteAPI(montarPayloadCliente());
      showSuccess("Cliente cadastrado com sucesso!");
    }
    
    // Sucesso, fechar o formulário e recarregar
    if (window.fecharFormularioInlineCliente) {
      window.fecharFormularioInlineCliente();
    }
    carregarClientes();
  } catch (err) {
    showError(err.message);
  }
}

async function salvarEndereco(e) {
  e.preventDefault();

  const idEndereco = document.getElementById("endereco_id").value;
  const payload = montarPayloadEndereco();

  try {
    if (idEndereco) {
      await atualizarEnderecoAPI(clienteIdAtual, idEndereco, payload);
      showSuccess("Endereço atualizado com sucesso!");
    } else {
      await criarEnderecoAPI(clienteIdAtual, payload);
      showSuccess("Endereço adicionado com sucesso!");
    }

    fecharModal("modalEnderecoCliente");
    await carregarDadosCliente(clienteIdAtual);
  } catch (err) {
    showError(err.message);
  }
}

async function salvarTelefone(e) {
  e.preventDefault();

  const idTelefone = document.getElementById("telefone_id").value;
  const payload = montarPayloadTelefone();

  try {
    if (idTelefone) {
      await atualizarTelefoneAPI(clienteIdAtual, idTelefone, payload);
      showSuccess("Telefone atualizado com sucesso!");
    } else {
      await criarTelefoneAPI(clienteIdAtual, payload);
      showSuccess("Telefone adicionado com sucesso!");
    }

    fecharModal("modalTelefoneCliente");
    await carregarDadosCliente(clienteIdAtual);
  } catch (err) {
    showError(err.message);
  }
}

function resetarAbaInicial() {
  const firstTab = document.getElementById("edit-basico-tab");
  if (firstTab) {
    bootstrap.Tab.getOrCreateInstance(firstTab).show();
  }
}

function configurarTipoPessoa() {
  const tipoSelect = document.getElementById("edit_cliente_tipo");
  if (!tipoSelect) return;

  tipoSelect.addEventListener("change", function () {
    const ieInput = document.getElementById("edit_cliente_ie");
    if (this.value === "PF") {
      ieInput.value = "";
      ieInput.disabled = true;
    } else {
      ieInput.disabled = false;
    }
  });
}

export function setupEditarCliente() {
  // Configurar modais de endereço e telefone
  const modalEndEl = document.getElementById("modalEnderecoCliente");
  const modalTelEl = document.getElementById("modalTelefoneCliente");
  if (modalEndEl) modalEndereco = bootstrap.Modal.getOrCreateInstance(modalEndEl);
  if (modalTelEl) modalTelefone = bootstrap.Modal.getOrCreateInstance(modalTelEl);

  window.abrirModalEditarCliente = async function (id) {
    clienteIdAtual = id;
    
    // Change Title
    const title = document.getElementById("formClienteTitulo");
    if (title) title.innerHTML = '<i class="bi bi-pencil-square me-2"></i> Editando Cliente #' + id;

    resetarAbaInicial();

    try {
      await carregarDadosCliente(id);
      if (window.abrirFormularioInlineCliente) {
        window.abrirFormularioInlineCliente();
      }
    } catch (err) {
      showError(err.message);
    }
  };

  window.resetarClienteIdAtual = function() {
    clienteIdAtual = null;
  };

  configurarTipoPessoa();

  const formEditar = document.getElementById("formEditarCliente");
  if (formEditar) formEditar.addEventListener("submit", salvarDadosCadastrais);

  const formEndereco = document.getElementById("formEnderecoCliente");
  if (formEndereco) formEndereco.addEventListener("submit", salvarEndereco);

  const formTelefone = document.getElementById("formTelefoneCliente");
  if (formTelefone) formTelefone.addEventListener("submit", salvarTelefone);

  const btnNovoEndereco = document.getElementById("btnNovoEndereco");
  if (btnNovoEndereco) {
    btnNovoEndereco.addEventListener("click", () => abrirModalEndereco());
  }

  const btnNovoTelefone = document.getElementById("btnNovoTelefone");
  if (btnNovoTelefone) {
    btnNovoTelefone.addEventListener("click", () => abrirModalTelefone());
  }
}
