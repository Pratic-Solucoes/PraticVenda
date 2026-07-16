import { obterItens, obterTotal } from './carrinho.js';

export function renderizarTabela() {
    const tbody = document.getElementById('tabela_pdv_body');
    const itens = obterItens();

    if (itens.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="7" class="text-muted py-5 text-center fs-5">
                    <i class="bi bi-bag-x fs-1 d-block mb-3 text-secondary"></i>
                    Nenhum produto adicionado à venda.
                </td>
            </tr>
        `;
        return;
    }

    tbody.innerHTML = '';
    itens.forEach((item, index) => {
        const produtoObj = item.produto.produto || item.produto; // Lida com ProdutoCompleto ou Produto
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td class="fw-bold">${index + 1}</td>
            <td>${produtoObj.codigo_barras || produtoObj.id || '-'}</td>
            <td class="text-start">${produtoObj.nome || 'Produto Sem Nome'}</td>
            <td>
                <input type="number" class="form-control form-control-sm text-center input-qtd mx-auto" style="width: 80px;" data-index="${index}" value="${item.quantidade}" min="1">
            </td>
            <td>R$ ${(produtoObj.preco_venda || 0).toFixed(2).replace('.', ',')}</td>
            <td class="fw-bold text-success">R$ ${item.subtotal.toFixed(2).replace('.', ',')}</td>
            <td>
                <button class="btn btn-sm btn-outline-danger btn-remover" data-index="${index}" title="Remover">
                    <i class="bi bi-trash"></i>
                </button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

export function atualizarTotalDisplay() {
    const total = obterTotal();
    const display = document.querySelector('.total-display');
    if (display) {
        display.textContent = `R$ ${total.toFixed(2).replace('.', ',')}`;
    }
}
