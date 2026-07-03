import { checkAuth, getToken } from '../utils/auth.js';
import { renderKPIs } from './kpis.js';
import { renderGraficoCategorias } from './graficos.js';

document.addEventListener('DOMContentLoaded', async () => {
    if (!checkAuth()) return;
    const token = getToken();

    try {
        const res = await fetch('/api/dashboard/resumo', {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (res.status === 401) {
            window.location.href = '/';
            return;
        }

        if (!res.ok) {
            console.error('Erro ao buscar dados do dashboard');
            return;
        }

        const data = await res.json();

        renderKPIs(data);
        renderGraficoCategorias(data.despesas_categoria || []);

    } catch (err) {
        console.error('Erro interno:', err);
    }
});
