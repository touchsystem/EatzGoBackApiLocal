package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CriarMesas(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		NumeroInicial int `json:"numero_inicial"`
		NumeroFinal   int `json:"numero_final"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	if payload.NumeroInicial > payload.NumeroFinal {
		respostas.Erro(w, http.StatusBadRequest, errors.New("o número inicial deve ser menor ou igual ao número final"))
		return
	}

	// Extrair código da empresa do token
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

	// Extrair nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("CriarMesas")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDeMesas(db)

	// Adicionar mesas
	for i := payload.NumeroInicial; i <= payload.NumeroFinal; i++ {
		existe, err := repositorio.VerificarSeMesaExiste(i)
		if err != nil {
			respostas.Erro(w, http.StatusInternalServerError, err)
			return
		}

		if existe {
			continue // Se já existe, pula para a próxima
		}

		mesa := modelos.Mesa{
			MesaCartao: i,
			Status:     "L", // Status "Livre" como padrão
		}

		// Validação do modelo Mesa
		if err := mesa.Validar(); err != nil {
			respostas.Erro(w, http.StatusBadRequest, fmt.Errorf("erro ao validar a mesa %d: %v", i, err))
			return
		}

		_, err = repositorio.CriarMesa(mesa)
		if err != nil {
			respostas.Erro(w, http.StatusInternalServerError, fmt.Errorf("erro ao criar a mesa %d: %v", i, err))
			return
		}
	}

	respostas.JSON(w, http.StatusCreated, map[string]string{"mensagem": "Mesas criadas com sucesso"})
}
func AtualizarMesa(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	mesaID, err := strconv.ParseUint(parametros["mesaId"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("ID da mesa inválido"))
		return
	}

	var mesa modelos.Mesa
	if err := json.NewDecoder(r.Body).Decode(&mesa); err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	// Validar modelo Mesa
	if err := mesa.Validar(); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	// Extrair código da empresa do token
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

	// Extrair nível do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("AtualizarMesa")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDeMesas(db)

	// Verificar se a mesa existe
	existe, err := repositorio.VerificarSeMesaExiste(mesa.MesaCartao)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if !existe {
		respostas.Erro(w, http.StatusNotFound, errors.New("mesa não encontrada"))
		return
	}

	// Atualizar mesa no banco
	if err := repositorio.AtualizarMesa(mesaID, mesa); err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarMesas busca mesas no banco de dados, com filtro opcional por STATUS
func BuscarMesas(w http.ResponseWriter, r *http.Request) {
	// Extrair o código da empresa do token
	cdEmp, err := autenticacao.ExtrairUsuarioCDEMP(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnauthorized, err)
		return
	}

	// Conectar ao banco de dados
	db, err := banco.ConectarPorEmpresa(cdEmp)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
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
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("BuscarMesas")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if nivelUsuario < uint64(nivelRequerido) {
		erro := errors.New("nível de acesso insuficiente")
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Extrair o filtro de STATUS da query string
	status := r.URL.Query().Get("status")

	// Validar STATUS, se informado
	if status != "" {
		statusPermitidos := map[string]bool{
			"L": true, // Livre
			"O": true, // Ocupado
			"F": true, // Fechamento
			"I": true, // Inativa
		}
		if _, permitido := statusPermitidos[status]; !permitido {
			respostas.Erro(w, http.StatusBadRequest, errors.New("status inválido. Valores permitidos: 'L', 'O', 'F', 'I'"))
			return
		}
	}

	// Buscar mesas no repositório
	repositorio := repositorios.NovoRepositorioDeMesas(db)
	mesas, err := repositorio.BuscarMesas(status)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Retornar as mesas
	respostas.JSON(w, http.StatusOK, mesas)
}

// BuscarMesa busca uma única mesa no banco de dados pelo ID
func BuscarMesa(w http.ResponseWriter, r *http.Request) {
	// Extrair os parâmetros da URL
	parametros := mux.Vars(r)
	mesaID, err := strconv.ParseUint(parametros["mesaId"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("ID da mesa inválido"))
		return
	}

	// Extrair o código da empresa do token
	cdEmp, err := autenticacao.ExtrairUsuarioCDEMP(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnauthorized, err)
		return
	}

	// Conectar ao banco de dados
	db, err := banco.ConectarPorEmpresa(cdEmp)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	// Buscar mesa no repositório
	repositorio := repositorios.NovoRepositorioDeMesas(db)
	mesa, err := repositorio.BuscarMesaPorID(mesaID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respostas.Erro(w, http.StatusNotFound, errors.New("mesa não encontrada"))
			return
		}
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Retornar a mesa encontrada
	respostas.JSON(w, http.StatusOK, mesa)
}
