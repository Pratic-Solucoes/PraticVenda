package auth

import (
	"context"
	"gestao/pkg/token"
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

		dadosUsuarioRequesicao, err := token.ValidarTokenJWT(tokenString)
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
