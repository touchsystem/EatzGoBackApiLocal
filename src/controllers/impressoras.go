// AtualizarImpressora altera as informações de uma impressora no banco de dados

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

// AtualizarImpressora altera as informações de uma impressora no banco de dados
func AtualizarImpressora(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	impressoraID, erro := strconv.ParseUint(parametros["impressoraId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Ler o corpo da requisição
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var impressora modelos.Impressora
	if erro = json.Unmarshal(corpoRequisicao, &impressora); erro != nil {
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
	//inicio nivel usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("AtualizarImpressora")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}
	//fim nível de usuário

	repositorio := repositorios.NovoRepositorioDeImpressoras(db)
	if erro = repositorio.AtualizarImpressora(impressoraID, impressora); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarImpressoras recupera todas as impressoras com filtros opcionais
func BuscarImpressoras(w http.ResponseWriter, r *http.Request) {
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

	//inicio nivel usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarImpressoras")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}
	//fim nível de usuário

	status := r.URL.Query().Get("status")
	filtro := r.URL.Query().Get("filtro")

	repositorio := repositorios.NovoRepositorioDeImpressoras(db)
	impressoras, erro := repositorio.BuscarImpressoras(status, filtro)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, impressoras)
}

// BuscarImpressoraPorID recupera uma impressora específica pelo ID
func BuscarImpressoraPorID(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	impressoraID, erro := strconv.Atoi(parametros["impressoraId"])
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

	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarImpressoraPorID")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeImpressoras(db)
	impressoraBusca, erro := repositorio.BuscarImpressoraPorID(impressoraID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, impressoraBusca)
}

// BuscarImpressoraPorCODIMP recupera uma impressora específica pelo COD_IMP
func BuscarImpressoraPorCODIMP(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	codImp := parametros["codImp"] // Obtém o COD_IMP da URL

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

	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarImpressoraPorCODIMP")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeImpressoras(db)
	impressoraBusca, erro := repositorio.BuscarImpressoraPorCODIMP(codImp)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, impressoraBusca)
}
