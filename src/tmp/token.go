package autenticacao

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CriarToken retorna um token assinado com as permissões do usuário
func CriarToken(usuarioID uint64, cdEmp string) (string, error) {
	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissoes["usuarioId"] = usuarioID
	permissoes["cdEmp"] = cdEmp // Adicionando o campo CDEMP
	//fmt.Println(usuarioID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
	return token.SignedString([]byte(config.SecretKey))

}

// ValidarToken verifica se o token passado na requisição é valido
func ValidarToken(r *http.Request) error {
	tokenString := extrairToken(r)

	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return erro
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("Token inválido")
}

// ExtrairUsuarioCDEMP retorna o cdEmp que está salvo no token
func ExtrairUsuarioCDEMP(r *http.Request) (string, error) {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)

	if erro != nil {
		return "", erro // Retorna uma string vazia em caso de erro
	}

	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		cdEmp, ok := permissoes["cdEmp"].(string)
		if !ok {
			return "", errors.New("cdEmp não é uma string") // Tratamento de erro se cdEmp não for string
		}

		//	fmt.Println(cdEmp) // Opcional: Se você quiser imprimir cdEmp
		return cdEmp, nil
	}

	return "", errors.New("Token inválido")
}

// ExtrairUsuarioID retorna o usuarioId que está salvo no token
func ExtrairUsuarioID(r *http.Request) (uint64, error) {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)

	if erro != nil {
		return 0, erro
	}

	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioID, erro := strconv.ParseUint(fmt.Sprintf("%.0f", permissoes["usuarioId"]), 10, 64)
		//	cdEmp, ok := permissoes["cdEmp"].(string)
		//	fmt.Println(cdEmp)
		//	fmt.Println(ok)
		if erro != nil {
			return 0, erro
		}

		return usuarioID, nil
	}

	//cdEmp := permissoes["cdemp"].(string)

	//fmt.Println(cdEmp)

	return 0, errors.New("Token inválido")
}

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func retornarChaveDeVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Método de assinatura inesperado! %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}
