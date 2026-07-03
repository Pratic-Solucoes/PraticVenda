export function renderKPIs(data) {
    const formatCurrency = (value) => new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(value);

    const elSemana = document.getElementById('kpi-semana');
    const elVencido = document.getElementById('kpi-vencido');
    const elVendaDia = document.getElementById('kpi-venda-dia');
    const elVendaMes = document.getElementById('kpi-venda-mes');

    if (elSemana) elSemana.textContent = formatCurrency(data.total_semana);
    if (elVencido) elVencido.textContent = formatCurrency(data.total_vencido);
    if (elVendaDia) elVendaDia.textContent = formatCurrency(data.venda_dia);
    if (elVendaMes) elVendaMes.textContent = formatCurrency(data.venda_mes);
}
