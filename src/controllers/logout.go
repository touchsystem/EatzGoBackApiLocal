package controllers

import (
	"api/src/autenticacao"
	"api/src/respostas"
	"net/http"
)

// Deslogar adiciona o token à blacklist e invalida o token JWT
func Deslogar(w http.ResponseWriter, r *http.Request) {
	// Extrair o token da requisição usando autenticacao.ExtrairToken
	tokenString := autenticacao.ExtrairToken(r)

	// Adicionar o token à blacklist
	autenticacao.InvalidateToken(tokenString)

	// Retornar resposta de sucesso
	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Token invalidado com sucesso"})
}
