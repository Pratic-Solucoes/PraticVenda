import { getToken } from "../utils/auth.js";
import { showError } from "../utils/showError.js";
import { validaRespostaRequisicao } from "../utils/resposta.js";
import { fecharModal } from "../utils/fecharModal.js";
import { carregarClientes } from "./listarClientes.js";

async function criarClienteAPI(payload) {
  const token = getToken();

  const res = await fetch("/api/clientes", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(payload),
  });

  return await validaRespostaRequisicao(res);
}

function montarPayloadNovoCliente() {
  const nome = document.getElementById("cliente_nome").value;
  const tipo = document.getElementById("cliente_tipo").value;
  const cpf_cnpj = document.getElementById("cliente_cpf_cnpj").value;
  const email = document.getElementById("cliente_email").value;

  const payload = { nome, tipo, email };

  if (tipo === "PF") {
    payload.cpf = cpf_cnpj;
  } else if (tipo === "PJ") {
    payload.cnpj = cpf_cnpj;
  }

  return payload;
}

async function handleSubmitNovoCliente(e) {
  e.preventDefault();

  try {
    await criarClienteAPI(montarPayloadNovoCliente());

    fecharModal("modalCliente");
    document.getElementById("formNovoCliente").reset();

    if (document.getElementById("tabela_clientes_body")) {
      carregarClientes();
    } else {
      window.location.reload();
    }
  } catch (err) {
    showError(err.message || "Erro interno ao comunicar com servidor.");
  }
}

export function setupCriarCliente() {
  const formNovo = document.getElementById("formNovoCliente");
  if (!formNovo) return;

  formNovo.addEventListener("submit", handleSubmitNovoCliente);
}
