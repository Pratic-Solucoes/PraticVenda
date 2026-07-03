export function getToken() {
    return localStorage.getItem('token');
}

export function checkAuth() {
    const token = getToken();
    if (!token) {
        window.location.href = '/';
        return false;
    }
    return token;
}
