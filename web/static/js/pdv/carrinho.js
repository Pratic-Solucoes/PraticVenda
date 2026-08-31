let itens = [];

export function adicionarItem(produto, quantidade = 1) {
	produto = produto.produto || produto;
    // Verifica se já existe o produto no carrinho para somar a quantidade
    const indexExistente = itens.findIndex(item => item.produto.id === produto.id);
    
    if (indexExistente >= 0) {
        itens[indexExistente].quantidade += quantidade;
		itens[indexExistente].subtotal = itens[indexExistente].quantidade * itens[indexExistente].precoUnitario;
    } else {
        itens.push({
            produto: produto,
            quantidade: quantidade,
			precoUnitario: produto.preco_venda,
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

export function definirItens(novosItens) {
    itens = novosItens.map(item => ({ produto: item.produto || item, quantidade: item.quantidade, precoUnitario: item.precoUnitario ?? item.valor_unitario, subtotal: item.quantidade * (item.precoUnitario ?? item.valor_unitario) }));
}

export function obterTotal() {
    return itens.reduce((total, item) => total + item.subtotal, 0);
}
