package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Login é responsável por autenticar um usuário na API
func Login(w http.ResponseWriter, r *http.Request) {
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

	db, erro := banco.Conectar("DB_NOME") // Conectar ao banco principal
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioSalvoNoBanco, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//fmt.Println(usuarioSalvoNoBanco.CDEMP)
	token, erro := autenticacao.CriarToken(usuarioSalvoNoBanco.ID, usuarioSalvoNoBanco.CDEMP, usuarioSalvoNoBanco.Nivel)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuarioID := strconv.FormatUint(usuarioSalvoNoBanco.ID, 10)

	respostas.JSON(w, http.StatusOK, modelos.DadosAutenticacao{ID: usuarioID, Token: token})
}

// BuscarDadosDoToken retorna os dados contidos no token do usuário autenticado, incluindo expiração, email, nome e nick
func BuscarDadosDoToken(w http.ResponseWriter, r *http.Request) {
	// Extrair dados básicos do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//fmt.Println(cdEmp)
	expTime, erro := autenticacao.ExtrairExpiracaoToken(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.Conectar("DB_NOME") // Conectar ao banco principal
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Criar um repositório para buscar o usuário pelo ID
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Retornar os dados como JSON, incluindo a expiração e dados adicionais do banco de dados
	dados := map[string]interface{}{
		"usuarioID": usuario.ID,
		"cdEmp":     cdEmp,
		"nivel":     usuario.Nivel,
		"exp":       expTime,
		"email":     usuario.Email,
		"nome":      usuario.Nome,
		"nick":      usuario.Nick,
	}

	respostas.JSON(w, http.StatusOK, dados)
}
