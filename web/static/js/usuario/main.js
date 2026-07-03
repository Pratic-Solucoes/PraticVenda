import { checkAuth } from '../utils/auth.js';
import { carregarDadosUsuario } from './carregarDadoUsuario.js';
import { alterarSenha } from './alterarSenha.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) return;

    carregarDadosUsuario();
    alterarSenha();
});
