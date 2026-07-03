package token

import (
	"fmt"
	"gestao/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GerarTokenJWT é uma função que gera um token JWT para um usuário autenticado, usando o ID e nome do usuário como payload
func GerarTokenJWT(usuarioID int, nome, schema string) (string, error) {
	var secretKey = config.GetToken("JWT_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usuario_id": usuarioID,
		"schema":     schema,
		"nome":       nome,
		"exp":        time.Now().Add(time.Hour * 4).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidarTokenJWT é uma função que valida um token JWT recebido,
// verificando sua assinatura e extraindo as claims (informações) contidas no token
func ValidarTokenJWT(tokenString string) (jwt.MapClaims, error) {
	var secretKey = config.GetToken("JWT_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to assert claims")
	}

	return claims, nil
}
