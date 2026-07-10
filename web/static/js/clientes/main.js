import { checkAuth } from "../utils/auth.js";
import { carregarClientes } from "./listarClientes.js";
import { setupEditarCliente } from "./editarCliente.js";

document.addEventListener("DOMContentLoaded", () => {
  if (!checkAuth()) return;

  const tbody = document.getElementById("tabela_clientes_body");
  const formFiltro = document.getElementById("formFiltroClientes");
  const inputBusca = document.getElementById("filtro_busca");

  // Elements for Inline Form
  const btnNovoCliente = document.getElementById("btnNovoCliente");
  const btnFecharFormCliente = document.getElementById("btnFecharFormCliente");
  const btnCancelarCliente = document.getElementById("btnCancelarCliente");
  const collapseFormCliente = document.getElementById("collapseFormCliente");
  const formClienteTitulo = document.getElementById("formClienteTitulo");
  const formEditarCliente = document.getElementById("formEditarCliente");
  const editClienteId = document.getElementById("edit_cliente_id");

  let bsCollapse = null;
  if (collapseFormCliente) {
    bsCollapse = new bootstrap.Collapse(collapseFormCliente, { toggle: false });
  }

  // Carregar clientes iniciais
  if (tbody) {
    carregarClientes();
  }

  // Filtro
  if (formFiltro) {
    formFiltro.addEventListener("submit", (e) => {
      e.preventDefault();
      carregarClientes(inputBusca?.value.trim() || "");
    });
  }

  // Handle Form toggle
  function abrirFormularioNovo() {
    formClienteTitulo.innerHTML = '<i class="bi bi-truck me-2"></i> Novo Cliente';
    formEditarCliente.reset();
    editClienteId.value = "";
    
    // Hide tabs except Dados Básicos for new creation
    const tabEnderecos = document.getElementById("edit-enderecos-tab");
    const tabTelefones = document.getElementById("edit-telefones-tab");
    if (tabEnderecos) tabEnderecos.parentElement.style.display = "none";
    if (tabTelefones) tabTelefones.parentElement.style.display = "none";

    const basicoTab = document.getElementById("edit-basico-tab");
    if (basicoTab) bootstrap.Tab.getOrCreateInstance(basicoTab).show();

    if (bsCollapse) bsCollapse.show();
  }

  function fecharFormulario() {
    if (bsCollapse) bsCollapse.hide();
    formEditarCliente.reset();
    editClienteId.value = "";
  }

  if (btnNovoCliente) {
    btnNovoCliente.addEventListener("click", abrirFormularioNovo);
  }

  if (btnFecharFormCliente) {
    btnFecharFormCliente.addEventListener("click", fecharFormulario);
  }

  if (btnCancelarCliente) {
    btnCancelarCliente.addEventListener("click", fecharFormulario);
  }

  // Make fecharFormulario global for other files to use if needed
  window.fecharFormularioInlineCliente = fecharFormulario;
  window.abrirFormularioInlineCliente = function() {
    if (bsCollapse) bsCollapse.show();
    
    // Show all tabs when editing
    const tabEnderecos = document.getElementById("edit-enderecos-tab");
    const tabTelefones = document.getElementById("edit-telefones-tab");
    if (tabEnderecos) tabEnderecos.parentElement.style.display = "block";
    if (tabTelefones) tabTelefones.parentElement.style.display = "block";
  };

  // Inicializa a lógica dos controllers
  setupEditarCliente();
});
