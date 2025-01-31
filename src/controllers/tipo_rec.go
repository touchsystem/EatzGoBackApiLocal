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

func CriarTipoRecebimento(w http.ResponseWriter, r *http.Request) {
	var tipoRecebimento modelos.TipoRecebimento
	if err := json.NewDecoder(r.Body).Decode(&tipoRecebimento); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	// Valida os dados antes de prosseguir
	if err := tipoRecebimento.Validar(); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	cdEmp, err := autenticacao.ExtrairUsuarioCDEMP(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnauthorized, err)
		return
	}

	db, err := banco.ConectarPorEmpresa(cdEmp)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("CriarTipoRecebimento")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeTiposRecebimento(db)
	id, err := repositorio.CriarTipoRecebimento(tipoRecebimento)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	tipoRecebimento.ID = id
	respostas.JSON(w, http.StatusCreated, tipoRecebimento)
}

func AtualizarTipoRecebimento(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	id, err := strconv.ParseUint(parametros["id"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("ID inválido"))
		return
	}

	var tipoRecebimento modelos.TipoRecebimento
	if err := json.NewDecoder(r.Body).Decode(&tipoRecebimento); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	// Valida antes de atualizar
	if err := tipoRecebimento.Validar(); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	cdEmp, err := autenticacao.ExtrairUsuarioCDEMP(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnauthorized, err)
		return
	}

	db, err := banco.ConectarPorEmpresa(cdEmp)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeTiposRecebimento(db)
	if err := repositorio.AtualizarTipoRecebimento(id, tipoRecebimento); err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

func BuscarTiposRecebimento(w http.ResponseWriter, r *http.Request) {
	cdEmp, err := autenticacao.ExtrairUsuarioCDEMP(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnauthorized, err)
		return
	}

	db, err := banco.ConectarPorEmpresa(cdEmp)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeTiposRecebimento(db)
	tiposPagamento, err := repositorio.BuscarTiposRecebimento()
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, tiposPagamento)
}
