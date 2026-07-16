let itens = [];

export function adicionarItem(produto, quantidade = 1) {
    // Verifica se já existe o produto no carrinho para somar a quantidade
    const indexExistente = itens.findIndex(item => item.produto.id === produto.id);
    
    if (indexExistente >= 0) {
        itens[indexExistente].quantidade += quantidade;
        itens[indexExistente].subtotal = itens[indexExistente].quantidade * itens[indexExistente].produto.preco_venda;
    } else {
        itens.push({
            produto: produto,
            quantidade: quantidade,
            subtotal: quantidade * produto.preco_venda
        });
    }
}

export function removerItem(index) {
    if (index >= 0 && index < itens.length) {
        itens.splice(index, 1);
    }
}

export function obterItens() {
    return itens;
}

export function limparCarrinho() {
    itens = [];
}

export function obterTotal() {
    return itens.reduce((total, item) => total + item.subtotal, 0);
}
