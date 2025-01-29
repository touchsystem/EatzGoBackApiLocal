package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// BuscarParametro busca um parâmetro pelo ID
func BuscarParametro(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	id, err := strconv.ParseUint(parametros["parametroId"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := banco.Conectar("DB_NOME")
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeParametros(db)
	parametro, err := repositorio.BuscarParametroPorID(id)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, parametro)
}

// BuscarParametros lista todos os parâmetros
// buscar parameteros

func BuscarParametros(w http.ResponseWriter, r *http.Request) {
	db, err := banco.Conectar("DB_NOME")
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeParametros(db)
	parametros, err := repositorio.BuscarParametros()
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, parametros)
}

// AtualizarParametro atualiza um parâmetro existente
func AtualizarParametro(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	id, err := strconv.ParseUint(parametros["parametroId"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	var parametro modelos.Parametro
	if err := json.NewDecoder(r.Body).Decode(&parametro); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	// Validação do modelo
	if err := parametro.Validar(); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := banco.Conectar("DB_NOME")
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	// Extrair nível do usuário e verificar permissões
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("AtualizarParametro")
	if erro != nil || nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}

	// Atualizar o parâmetro no banco
	repositorio := repositorios.NovoRepositorioDeParametros(db)
	if err := repositorio.AtualizarParametro(id, parametro); err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
