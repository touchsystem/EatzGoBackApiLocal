package autenticacao

var blacklistTokens = make(map[string]bool)

// InvalidateToken adiciona o token à blacklist
func InvalidateToken(tokenString string) {
	blacklistTokens[tokenString] = true
}

// VerificarBlacklist verifica se o token está na blacklist
func VerificarBlacklist(tokenString string) bool {
	_, exists := blacklistTokens[tokenString]
	return exists
}
