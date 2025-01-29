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

// CriarCliente insere uma nova entrada de cliente no banco de dados
func CriarCliente(w http.ResponseWriter, r *http.Request) {
	var cliente modelos.Cliente
	erro := json.NewDecoder(r.Body).Decode(&cliente)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

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

	//ExtrairUsuarioNivel busca o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	// Buscar o nível de acesso para a operação CLI-BUSCAR
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("CriarCliente")
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
	sync := strings.Contains(r.URL.Path, "sync")
	if sync {
		//	fmt.Println("verdade")

	} else {
		//	fmt.Println("false")
		cliente.ID = 0
		cliente.ID_HOSTLOCAL = 0
	}
	repositorio := repositorios.NovoRepositorioDeClientes(db)
	idCliente, erro := repositorio.CriarCliente(cliente)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	cliente.ID = idCliente
	respostas.JSON(w, http.StatusCreated, cliente)
}

// BuscarClientes busca todos os clientes salvos no banco
// BuscarClientes busca todos os clientes salvos no banco
func BuscarClientes(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")     // Captura o valor combinado de busca (nome, CPF, celular)
	status := r.URL.Query().Get("status") // Captura o valor do status, se fornecido

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

	// Faz checagem do nível do Usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarClientes")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Repositório e busca
	repositorio := repositorios.NovoRepositorioDeClientes(db)
	clientes, erro := repositorio.BuscarClientes(nome, status)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, clientes)
}

// BuscarCliente busca um cliente no banco de dados pelo ID
func BuscarCliente(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	clienteID, erro := strconv.ParseUint(parametros["clienteId"], 10, 64)
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

	// Checagem do nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarCliente")
	if erro != nil || nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeClientes(db)
	cliente, erro := repositorio.BuscarPorID(clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, cliente)
}

// AtualizarCliente altera as informações de um cliente no banco
func AtualizarCliente(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	//fmt.Println("entrei a")
	clienteID, erro := strconv.ParseUint(parametros["clienteId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var cliente modelos.Cliente
	if erro = json.Unmarshal(corpoRequisicao, &cliente); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	sync := strings.Contains(r.URL.Path, "/sync/")
	if sync {
		cliente.AGUARDANDO_SYNC = ""
	} else {
		cliente.AGUARDANDO_SYNC = "S"
	}

	if erro = cliente.Preparar("edicao"); erro != nil {
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

	if cliente.AGUARDANDO_SYNC == "S" {
		//ExtrairUsuarioNivel busca o nível do usuário
		nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
		if erro != nil {
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
		// Buscar o nível de acesso para a operação CLI-BUSCAR
		nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("AtualizarCliente")
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

	//	fmt.Println("enttrei 00")
	repositorio := repositorios.NovoRepositorioDeClientes(db)
	if erro = repositorio.AtualizarCliente(clienteID, cliente); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarCliente exclui um cliente no banco
func DeletarCliente(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	clienteID, erro := strconv.ParseUint(parametros["clienteId"], 10, 64)
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

	// Checagem do nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("DeletarCliente")
	if erro != nil || nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeClientes(db)
	if erro = repositorio.DeletarCliente(clienteID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
