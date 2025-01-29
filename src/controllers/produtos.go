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

// CriarProduto insere um produto no banco de dados
func CriarProduto(w http.ResponseWriter, r *http.Request) {

	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var produto modelos.Produto
	if erro = json.Unmarshal(corpoRequest, &produto); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = produto.Preparar("cadastro"); erro != nil {
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
	//*****Faz checagem do nível do Usuário

	//ExtrairUsuarioNivel busca o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	// Buscar o nível de acesso para a operação CLI-BUSCAR
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("CriarProduto")
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
		produto.ID = 0
		produto.ID_HOSTLOCAL = 0
	}
	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produto.ID, erro = repositorio.CriarProdutos(produto)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, produto)
}

// BuscarProdutos busca todos os produtos salvos no banco com filtro de status e LIKE em outros campos
func BuscarProdutos(w http.ResponseWriter, r *http.Request) {
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

	// Extrair o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarProdutos")
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

	// Repositório e busca de produtos com filtro
	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtos, erro := repositorio.BuscarProdutos(status, filtro)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, produtos)
}

// BuscarProduto busca um produto salvo no banco
func BuscarProduto(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	produtoID, erro := strconv.ParseUint(parametros["produtoId"], 10, 64)
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
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarProduto")
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

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produto, erro := repositorio.BuscarProdutoPorID(produtoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, produto)
}

////*******

// AtualizarProduto altera as informações de um produto no banco
func AtualizarProduto(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID, erro := strconv.ParseUint(parametros["produtoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var produto modelos.Produto
	if erro = json.Unmarshal(corpoRequisicao, &produto); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	sync := strings.Contains(r.URL.Path, "/sync/")
	if sync {
		produto.AGUARDANDO_SYNC = ""

	} else {
		produto.AGUARDANDO_SYNC = "S"

	}

	if erro = produto.Preparar("edicao"); erro != nil {
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
	if produto.AGUARDANDO_SYNC == "S" {
		//ExtrairUsuarioNivel busca o nível do usuário
		nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
		if erro != nil {
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
		// Buscar o nível de acesso para a operação CLI-BUSCAR
		nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("AtualizarProduto")
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
	//*****fim da checagem do nível do usuário

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	if erro = repositorio.AtualizarProdutos(produtoID, produto); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarProduto exclui as informações de um produto no banco

// DeletarProduto exclui as informações de um produto no banco
func DeletarProduto(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	//fmt.Println("aaa")
	produtoID, erro := strconv.ParseUint(parametros["produtoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//	fmt.Println("sss")
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

	// Extrair o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("DeletarProduto")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Criar o objeto do produto com as informações necessárias
	produto := modelos.Produto{
		ID:     produtoID,
		STATUS: "D", // Define o status como "Deletado"
		CODM:   "X",
		DES2:   "X",
	}

	// Validar o produto com o status atualizado
	if erro := produto.Validar(""); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Repositório e exclusão
	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	if erro := repositorio.DeletarProduto(produtoID, produto); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarProdutosPreferidos busca produtos com STATUS = 'C' e OB4 = 'S'
func BuscarPreferidos(w http.ResponseWriter, r *http.Request) {
	// Extrair o código da empresa do token
	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Conectar ao banco de dados
	db, erro := banco.ConectarPorEmpresa(cdEmp)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Extrair nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarPreferidos")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Repositório e busca de produtos com STATUS = 'C' e OB4 = 'S'
	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtos, erro := repositorio.BuscarPreferidos()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Retornar os produtos preferidos
	respostas.JSON(w, http.StatusOK, produtos)
}

// BuscarProdutosPorGrupo busca produtos filtrados por ID de grupo
func BuscarProdutosPorGrupo(w http.ResponseWriter, r *http.Request) {
	// Extrair o ID do grupo dos parâmetros da URL
	parametros := mux.Vars(r)
	grupoID, erro := strconv.ParseUint(parametros["grupoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("ID do grupo inválido"))
		return
	}

	// Extrair o código da empresa do token
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

	// Extrair o nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarProdutosPorGrupo")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Repositório e busca de produtos por grupo
	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtos, erro := repositorio.BuscarProdutosPorGrupo(grupoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, produtos)
}

// BuscarProdutoPorCODM busca um produto específico pelo CODM
func BuscarProdutoPorCODM(w http.ResponseWriter, r *http.Request) {
	codm := r.URL.Query().Get("codm")
	if codm == "" {
		respostas.Erro(w, http.StatusBadRequest, errors.New("O parâmetro CODM é obrigatório"))
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
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarProdutoPorCODM")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produto, erro := repositorio.BuscarProdutoPorCODM(codm)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, produto)
}
