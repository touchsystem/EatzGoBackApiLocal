package autenticacao

import (
	"api/src/banco"
	"api/src/config"
	"api/src/repositorios"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CriarToken retorna um token assinado com as permissões do usuário
func CriarToken(usuarioID uint64, cdEmp string, nivel uint64) (string, error) {
	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissoes["usuarioId"] = usuarioID
	permissoes["cdEmp"] = cdEmp // Adicionando o campo CDEMP
	permissoes["nivel"] = nivel
	//	fmt.Println(nivel)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
	return token.SignedString([]byte(config.SecretKey))

}

// CriarToken retorna um token assinado com as permissões do usuário
func CriarToken2(usuarioID uint64, cdEmp string, nivel uint64) (string, error) {
	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 150000).Unix()
	permissoes["usuarioId"] = usuarioID
	permissoes["cdEmp"] = cdEmp // Adicionando o campo CDEMP
	permissoes["-1"] = nivel
	//fmt.Println(nivel)
	//fmt.Println(usuarioID)
	//fmt.Println(cdEmp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
	return token.SignedString([]byte(config.SecretKey))

}

// ValidarToken verifica se o token passado na requisição é valido
func ValidarToken(r *http.Request) error {
	tokenString := ExtrairToken(r)

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
	tokenString := ExtrairToken(r)
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

// ExtrairUsuarioNivel retorna o cdEmp que está salvo no token
func ExtrairUsuarioNivel(r *http.Request) (uint64, error) {
	tokenString := ExtrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)

	if erro != nil {
		return 0, erro
	}

	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verifica se "usuarioID" está presente e se é um número
		nivelFloat, ok := permissoes["nivel"].(float64)
		if !ok {
			return 0, errors.New("nivel não é um número válido")
		}
		//	fmt.Println(nivelFloat)
		//	fmt.Println(r)
		// Converte o float64 para uint64
		nivel := uint64(nivelFloat)
		return nivel, nil
	}

	return 0, errors.New("Token inválido EXTRAINDO Nível")
}

// ExtrairUsuarioID retorna o usuarioId que está salvo no token
func ExtrairUsuarioID(r *http.Request) (uint64, error) {
	tokenString := ExtrairToken(r)
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

func ExtrairToken(r *http.Request) string {
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

// ExtrairCdEmpDoTokenString extrai `cdEmp` do token de sincronização fornecido como string
func ExtrairCdEmpDoTokenString(tokenString string) (string, error) {
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return "", erro
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if cdEmp, ok := claims["cdEmp"].(string); ok {
			return cdEmp, nil
		}
		return "", errors.New("cdEmp não encontrado no token")
	}

	return "", errors.New("token inválido")
}

func ExtrairExpiracaoToken(r *http.Request) (time.Time, error) {
	tokenString := ExtrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)

	if erro != nil {
		return time.Time{}, erro
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return time.Time{}, errors.New("campo exp não é um número válido")
		}

		// Converter exp (Unix timestamp) para time.Time
		expTime := time.Unix(int64(exp), 0)
		return expTime, nil
	}

	return time.Time{}, errors.New("token inválido")
}

// ExtrairUsuarioNick obtém o Nick do usuário autenticado
func ExtrairUsuarioNick(r *http.Request) (string, error) {
	// Extrai o ID do usuário do token
	idUsuario, erro := ExtrairUsuarioID(r)
	if erro != nil {
		return "", errors.New("Erro ao extrair o ID do usuário: " + erro.Error())
	}

	// Conecta ao banco de dados principal
	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		return "", errors.New("Erro ao conectar ao banco de dados: " + erro.Error())
	}
	defer db.Close()

	// Cria um repositório de usuários
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	// Busca o usuário pelo ID
	usuario, erro := repositorio.BuscarPorID(idUsuario)
	if erro != nil {
		return "", errors.New("Erro ao buscar o usuário: " + erro.Error())
	}

	// Verifica se o nick está preenchido
	if usuario.Nick == "" {
		return "", errors.New("Nick do usuário não encontrado")
	}

	return usuario.Nick, nil
}
