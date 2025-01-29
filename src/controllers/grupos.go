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

	"github.com/gorilla/mux"
)

// CriarGrupo insere um grupo no banco de dados
func CriarGrupo(w http.ResponseWriter, r *http.Request) {

	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var grupo modelos.Grupo
	if erro = json.Unmarshal(corpoRequest, &grupo); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	//	if erro = grupo.Preparar("cadastro"); erro != nil {
	//		respostas.Erro(w, http.StatusBadRequest, erro)
	//		return
	//	}

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

	//*****Faz checagem do nível do Usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("CriarGrupo")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//*****fim da checagem do nível do usuário

	sync := strings.Contains(r.URL.Path, "sync")
	if sync {
		// Sincronização ativada
	} else {
		grupo.ID = 0
		grupo.ID_HOSTLOCAL = 0
	}
	repositorio := repositorios.NovoRepositorioDeGrupos(db)
	grupo.ID, erro = repositorio.CriarGrupo(grupo)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, grupo)
}

// BuscarGrupos recupera todos os registros de grupos
func BuscarGrupos(w http.ResponseWriter, r *http.Request) {
	// Extrair o código da empresa do token
	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Conectar ao banco de dados com o código da empresa
	db, erro := banco.ConectarPorEmpresa(cdEmp)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()
	// Extrair o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarGrupos")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	// Obter os parâmetros de status e filtro da URL
	status := r.URL.Query().Get("status")
	filtro := r.URL.Query().Get("filtro")

	repositorio := repositorios.NovoRepositorioDeGrupos(db)
	grupos, erro := repositorio.BuscarGrupos(status, filtro)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, grupos)
}

// BuscarGrupo busca um grupo salvo no banco
func BuscarGrupo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	grupoID, erro := strconv.ParseUint(parametros["grupoId"], 10, 64)
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
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarGrupo")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	//*****fim da checagem do nível do usuário

	repositorio := repositorios.NovoRepositorioDeGrupos(db)
	grupo, erro := repositorio.BuscarGrupoPorID(grupoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, grupo)
}

// AtualizarGrupo altera as informações de um grupo no banco
func AtualizarGrupo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	grupoID, erro := strconv.ParseUint(parametros["grupoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var grupo modelos.Grupo
	if erro = json.Unmarshal(corpoRequisicao, &grupo); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	sync := strings.Contains(r.URL.Path, "/sync/")
	if sync {
		grupo.AGUARDANDO_SYNC = ""
	} else {
		grupo.AGUARDANDO_SYNC = "S"
	}

	//	if erro = grupo.Preparar("edicao"); erro != nil {
	//		respostas.Erro(w, http.StatusBadRequest, erro)
	//		return
	//	}

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

	if grupo.AGUARDANDO_SYNC == "S" {
		nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
		if erro != nil {
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
		nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("AtualizarGrupo")
		if erro != nil {
			respostas.Erro(w, http.StatusInternalServerError, erro)
			return
		}
		if nivelUsuario < uint64(nivelRequerido) {
			erro := errors.New("nível de acesso insuficiente")
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}
	}

	repositorio := repositorios.NovoRepositorioDeGrupos(db)
	if erro = repositorio.AtualizarGrupo(grupoID, grupo); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarGrupo exclui as informações de um grupo no banco
func DeletarGrupo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	grupoID, erro := strconv.ParseUint(parametros["grupoId"], 10, 64)
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
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("DeletarGrupo")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDeGrupos(db)
	if erro = repositorio.DeletarGrupo(grupoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
