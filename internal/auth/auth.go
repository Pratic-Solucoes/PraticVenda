package auth

import (
	"context"
	"net/http"
	"strings"
)

func Autenticar(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Token ausente ou mal formatado", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		dadosUsuarioRequesicao, err := ValidarTokenJWT(tokenString)
		if err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "usuario_id", dadosUsuarioRequesicao["usuario_id"])
		ctx = context.WithValue(ctx, "usuario_nome", dadosUsuarioRequesicao["nome"])
		ctx = context.WithValue(ctx, "schema", dadosUsuarioRequesicao["schema"])

		r = r.WithContext(ctx)

		next(w, r)
	}
}

// AutenticarAdministrador restringe operações globais aos tokens emitidos
// pelo login administrativo.
func AutenticarAdministrador(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Token ausente ou mal formatado", http.StatusUnauthorized)
			return
		}

		dadosUsuarioRequesicao, err := ValidarTokenJWT(strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil || dadosUsuarioRequesicao["administrador"] != true {
			http.Error(w, "Acesso administrativo não autorizado", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "usuario_id", dadosUsuarioRequesicao["usuario_id"])
		ctx = context.WithValue(ctx, "usuario_nome", dadosUsuarioRequesicao["nome"])
		ctx = context.WithValue(ctx, "schema", "public")
		next(w, r.WithContext(ctx))
	}
}
