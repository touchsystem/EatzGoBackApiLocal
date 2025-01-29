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
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func CriarCxReceb(w http.ResponseWriter, r *http.Request) {
	// Lendo o corpo da requisição
	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	// Estrutura para decodificar JSON
	var payload struct {
		CxReceb      modelos.CxReceb       `json:"cx_receb"`
		CxRecebTipos []modelos.CxRecebTipo `json:"cx_receb_tipos"`
	}

	// Decodificando JSON
	if erro = json.Unmarshal(corpoRequest, &payload); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Parsing manual do campo "data" no formato "YYYY-MM-DD"
	parsedDate, erro := time.Parse("2006-01-02", payload.CxReceb.Data.Format("2006-01-02"))
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("data no formato inválido, use YYYY-MM-DD"))
		return
	}
	payload.CxReceb.Data = parsedDate

	// Extraindo o CDEMP e conectando ao banco
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

	// Verificando o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("CriarCxReceb")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}

	// Sincronização e criação de dados
	sync := strings.Contains(r.URL.Path, "sync")
	if sync {
		// Lógica de sincronização pode ser aplicada aqui
	} else {
		payload.CxReceb.IDHostWeb = 0
		for i := range payload.CxRecebTipos {
			payload.CxRecebTipos[i].IDHostWeb = 0
		}
	}

	// Criando CX_RECEB e associando os tipos
	repositorio := repositorios.NovoRepositorioCxReceb(db)
	idCxReceb, erro := repositorio.CriarCxReceb(payload.CxReceb)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Atualizando os tipos com o ID gerado
	for i := range payload.CxRecebTipos {
		payload.CxRecebTipos[i].IDCxReceb = idCxReceb
	}

	erro = repositorio.CriarCxRecebTipos(payload.CxRecebTipos)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, map[string]interface{}{
		"id_cx_receb": idCxReceb,
	})
}

func BuscarCxReceb(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	idCxReceb, erro := strconv.Atoi(parametros["idCxReceb"])
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

	// Verificando o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarCxReceb")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}

	repositorio := repositorios.NovoRepositorioCxReceb(db)
	cxReceb, cxRecebTipos, erro := repositorio.BuscarCxRecebPorID(idCxReceb)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, map[string]interface{}{
		"cx_receb":       cxReceb,
		"cx_receb_tipos": cxRecebTipos,
	})
}
