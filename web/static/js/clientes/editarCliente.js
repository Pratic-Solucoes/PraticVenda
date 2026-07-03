import { getToken } from "../utils/auth.js";
import { carregarClientes} from "./listarClientes.js";
import { showError } from '/static/js/utils/showError.js';

export function setupEditarCliente() {
  const formEditar = document.getElementById('formEditarCliente');

  window.abrirModalEditarCliente = async function (id) {
    const token = getToken();
    try {
      const res = await fetch(`/api/clientes/${id}`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (!res.ok) {
        showError("Erro ao buscar cliente.");
        return
      }

      const cliente = await res.json();

      document.getElementById('edit_cliente_id').value = cliente.id;
      document.getElementById('edit_cliente_nome').value = cliente.nome;
      document.getElementById('edit_cliente_tipo').value = cliente.tipo;
      
      document.getElementById('edit_cliente_email').value = cliente.email || "";
      document.getElementById('edit_cliente_telefone').value = cliente.telefone || "";
      document.getElementById('edit_cliente_contribuinte').value = cliente.contribuinte || 9;
      document.getElementById('edit_cliente_is_consumidor_final').checked = cliente.is_consumidor_final || false;
      document.getElementById('edit_cliente_ie').value = cliente.ie || "";

      if (cliente.tipo == 'PF') {
        document.getElementById('edit_cliente_cpf_cnpj').value = cliente.cpf || "";
        document.getElementById('edit_cliente_ie').disabled = true;
      } else {
        document.getElementById('edit_cliente_cpf_cnpj').value = cliente.cnpj || "";
        document.getElementById('edit_cliente_ie').disabled = false;
      }

      const modal = bootstrap.Modal.getOrCreateInstance(document.getElementById('modalEditarCliente'));
      modal.show();

    } catch (err) {
      console.error(err);
      showError("Erro interno ao buscar cliente.");
    }
  };

  const tipoSelect = document.getElementById('edit_cliente_tipo');
  if (tipoSelect) {
      tipoSelect.addEventListener('change', function() {
          const ieInput = document.getElementById('edit_cliente_ie');
          if (this.value === 'PF') {
              ieInput.value = "";
              ieInput.disabled = true;
          } else {
              ieInput.disabled = false;
          }
      });
  }

  if (formEditar) {
    formEditar.addEventListener('submit', async (e) => {
      e.preventDefault();

      const token = getToken();
      const id = document.getElementById('edit_cliente_id').value;
      const nome = document.getElementById('edit_cliente_nome').value;
      const tipo = document.getElementById('edit_cliente_tipo').value;
      const cpf_cnpj = document.getElementById('edit_cliente_cpf_cnpj').value;
      const email = document.getElementById('edit_cliente_email').value;
      const telefone = document.getElementById('edit_cliente_telefone').value;
      const ie = document.getElementById('edit_cliente_ie').value;
      const contribuinte = parseInt(document.getElementById('edit_cliente_contribuinte').value, 10);
      const is_consumidor_final = document.getElementById('edit_cliente_is_consumidor_final').checked;

      const payload = {
        nome, tipo, email, telefone, ie, contribuinte, is_consumidor_final
      };

      if (tipo === 'PF') {
        payload.cpf = cpf_cnpj;
        payload.cnpj = "";
      } else if (tipo === 'PJ') {
        payload.cnpj = cpf_cnpj;
        payload.cpf = "";
      }

      try {
        const res = await fetch(`/api/clientes/${id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(payload)
        });

        if (!res.ok) {
          const data = await res.json();
          showError(data.erro || "Erro ao atualizar cliente.");
          return;
        }

        const modal = bootstrap.Modal.getInstance(document.getElementById('modalEditarCliente'));
        if (modal) modal.hide();

        formEditar.reset();

        const tbody = document.getElementById('tabela_clientes_body');
        if (tbody) {
            carregarClientes();
        } else {
            window.location.reload();
        }

      } catch (err) {
        console.error(err);
        showError("Erro interno ao comunicar com servidor.");
      }
    });
  }
}
