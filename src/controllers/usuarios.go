package controllers

import (
	"api/src/autenticacao"
	"api/src/banco" // Import para acessar a função GetStringConexaoBanco
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CriarUsuario insere um usuário no banco de dados principal
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	tokenSync := os.Getenv("TOKEN_SYNC")
	cdEmp, erro := autenticacao.ExtrairCdEmpDoTokenString(tokenSync)
	if erro != nil {
		log.Println("Erro ao extrair cdEmp do token:", erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Verificar se o usuário existe na nuvem
	url := fmt.Sprintf("http://%s:5000/usuarios/verificar-existencia", os.Getenv("IP_NUVEM"))
	payload, _ := json.Marshal(map[string]string{"email": usuario.Email, "nick": usuario.Nick})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+tokenSync)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		respostas.Erro(w, http.StatusBadGateway, errors.New("Erro ao verificar usuário na nuvem"))
		return
	}
	defer resp.Body.Close()

	var resultado struct {
		Existe bool `json:"existe"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil || resultado.Existe {
		respostas.Erro(w, http.StatusConflict, errors.New("Usuário já existe na nuvem"))
		return
	}

	// Criptografar a senha
	hash, err := seguranca.Hash(usuario.Senha)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("Erro ao criptografar a senha"))
		return
	}
	usuario.Senha = string(hash)

	// Prosseguir com o cadastro local
	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario.CDEMP = cdEmp
	usuario.ID, erro = repositorio.CriarUsuario(usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, usuario)
}

// BuscarUsuarios busca todos os usuários no banco de dados principal
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	// Conectar ao banco de dados principal
	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	//ExtrairUsuarioNivel busca o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	// Buscar o nível de acesso para a operação CLI-BUSCAR
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarUsuarios")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	// Verificar se o nível do usuário é suficiente
	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	// *****fim da checagem do nível do usuário
	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//*****fim da checagem do nível do usuário

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.BuscarUsuarios(nomeOuNick, cdEmp)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// BuscarUsuario busca um usuário específico no banco de dados principal
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Conectar ao banco de dados principal
	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	//ExtrairUsuarioNivel busca o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	// Buscar o nível de acesso para a operação CLI-BUSCAR
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarUsuario")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	// Verificar se o nível do usuário é suficiente
	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	//*****fim da checagem do nível do usuário

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuario)
}

// AtualizarUsuario altera as informações de um usuário no banco de dados principal
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	sync := strings.Contains(r.URL.Path, "/sync/")
	if sync {
		usuario.AGUARDANDO_SYNC = ""
	} else {
		usuario.AGUARDANDO_SYNC = "S"
	}
	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Conectar ao banco de dados principal
	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	if usuario.AGUARDANDO_SYNC == "S" {
		//ExtrairUsuarioNivel busca o nível do usuário
		nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
		if erro != nil {
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
		// Buscar o nível de acesso para a operação CLI-BUSCAR
		nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("AtualizarUsuario")
		if erro != nil {
			respostas.Erro(w, http.StatusInternalServerError, erro)
			return
		}
		// Verificar se o nível do usuário é suficiente
		if nivelUsuario < uint64(nivelRequerido) {
			erro := errors.New("nível de acesso insuficiente")
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}

		//*****fim da checagem do nível do usuário
	}
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.AtualizarUsuario(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarUsuario exclui as informações de um usuário no banco de dados principal
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível deletar um usuário que não seja o seu"))
		return
	}

	// Conectar ao banco de dados principal
	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	//ExtrairUsuarioNivel busca o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	// Buscar o nível de acesso para a operação CLI-BUSCAR
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("DeletarUsuario")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	// Verificar se o nível do usuário é suficiente
	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	//*****fim da checagem do nível do usuário

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// AtualizarSenha permite alterar a senha de um usuário
func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioIDNoToken != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar a senha de um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Conectar ao banco de dados principal
	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	senhaSalvaNoBanco, erro := repositorio.BuscarSenha(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("A senha atual não condiz com a que está salva no banco"))
		return
	}

	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.AtualizarSenha(usuarioID, string(senhaComHash)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarNick busca usuário por nick
func BuscarNick(w http.ResponseWriter, r *http.Request) {
	// Decodificar o JSON do corpo da requisição
	var body struct {
		Nick string `json:"nick"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("JSON inválido"))
		return
	}

	// Garantir que o campo "nick" foi informado
	if body.Nick == "" {
		respostas.Erro(w, http.StatusBadRequest, errors.New("O campo 'nick' é obrigatório"))
		return
	}

	// Extrair o código da empresa do token
	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Conectar ao banco de dados principal
	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Buscar usuário pelo nick
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	fmt.Println(body.Nick)
	usuarios, erro := repositorio.BuscarNick(body.Nick, cdEmp)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Retornar o resultado
	respostas.JSON(w, http.StatusOK, usuarios)
}
