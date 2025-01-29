package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// BuscarNiveis busca todos os níveis de acesso no banco de dados
func BuscarNiveis(w http.ResponseWriter, r *http.Request) {
	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.ConectarPorEmpresa(cdEmp)
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
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarNiveis")
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

	repositorio := repositorios.NovoRepositorioDeNiveis(db)
	niveis, erro := repositorio.BuscarNiveis()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, niveis)
}

// BuscarNivel busca um nível de acesso específico por ID
func BuscarNivel(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	nivelID, erro := strconv.ParseUint(parametros["nivelId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.ConectarPorEmpresa(cdEmp)
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
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarNivel")
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

	repositorio := repositorios.NovoRepositorioDeNiveis(db)
	nivel, erro := repositorio.BuscarNivelPorID(nivelID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, nivel)
}

// AtualizarNivel altera as informações de um nível de acesso no banco de dados
func AtualizarNivel(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	nivelID, erro := strconv.ParseUint(parametros["nivelId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var nivel modelos.NivelAcesso
	if erro = json.Unmarshal(corpoRequest, &nivel); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = nivel.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.ConectarPorEmpresa(cdEmp)
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
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("AtualizarNivel")
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

	repositorio := repositorios.NovoRepositorioDeNiveis(db)
	if erro = repositorio.AtualizarNivel(nivelID, nivel); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
