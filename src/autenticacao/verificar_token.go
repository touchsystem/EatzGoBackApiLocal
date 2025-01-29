package autenticacao

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// VerificarToken verifica se o token é válido e se não está na blacklist
func VerificarToken(r *http.Request) error {
	// Chamada da função ExtrairToken, que deve estar no mesmo pacote
	tokenString := ExtrairToken(r)

	// Verificar se o token está na blacklist
	if VerificarBlacklist(tokenString) {
		return errors.New("token invalidado")
	}

	token, err := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if err != nil {
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return errors.New("token expirado")
		}
		return errors.New("token inválido")
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("token inválido")
}
