import { getToken } from "../utils/auth.js";
import { showError } from "../utils/showError.js";
import { validaRespostaRequisicao } from "../utils/resposta.js";

async function buscarClientesAPI(busca) {
  const token = getToken();

  let url = "/api/clientes";
  if (busca) {
    url += `?busca=${encodeURIComponent(busca)}`;
  }

  const res = await fetch(url, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return await validaRespostaRequisicao(res);
}

function formatarDocumento(cliente) {
  if (cliente.cpf) return cliente.cpf;
  if (cliente.cnpj) return cliente.cnpj;
  return "-";
}

function formatarTelefone(cliente) {
  if (cliente.telefones && cliente.telefones.length > 0) {
    const tel = cliente.telefones[0];
    return `(${tel.ddd}) ${tel.numero}`;
  }
  return "-";
}

function renderTabela(clientes) {
  const tbody = document.getElementById("tabela_clientes_body");
  if (!tbody) return;

  if (!clientes || clientes.length === 0) {
    tbody.innerHTML = `<tr><td colspan="7" class="text-muted py-5 text-center">Nenhum cliente encontrado.</td></tr>`;
    return;
  }

  tbody.innerHTML = "";
  clientes.forEach((c) => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
            <td>#${c.id}</td>
            <td class="text-start fw-bold">${c.nome}</td>
            <td>${c.email || "-"}</td>
            <td>${formatarTelefone(c)}</td>
            <td><span class="badge bg-secondary">${c.tipo || "-"}</span></td>
            <td>${formatarDocumento(c)}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary" title="Editar" onclick="window.abrirModalEditarCliente(${c.id})"><i class="bi bi-pencil"></i></button>
            </td>
        `;
    tbody.appendChild(tr);
  });
}

export async function carregarClientes(busca = "") {
  const tbody = document.getElementById("tabela_clientes_body");
  if (!tbody) return;

  tbody.innerHTML = `<tr><td colspan="7" class="text-muted py-5 text-center">Carregando...</td></tr>`;

  try {
    const dados = await buscarClientesAPI(busca);
    renderTabela(dados);
  } catch (err) {
    showError(err.message || "Erro ao carregar clientes. Tente novamente mais tarde.");
  }
}
