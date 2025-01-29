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
	"time"
)

// BuscarDatasSistema retorna a data do sistema, o primeiro dia do mês e o último dia do mês
// BuscarDatasSistema retorna a data do sistema, o primeiro dia do mês, o último dia do mês, o cdEmp e o nome da empresa
func BuscarDatasSistema(w http.ResponseWriter, r *http.Request) {
	// Extrair o código da empresa do token

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
	//fmt.Println("enrei  0")
	// Criar repositório e buscar a data do sistema
	repositorio := repositorios.NovoRepositorioDeEmpresas(db)
	dataSistema, erro := repositorio.BuscarDataSistema()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	//	fmt.Println("enrei ")
	// Calcular o primeiro e o último dia do mês
	dataFormatada, _ := time.Parse("2006-01-02", dataSistema)
	primeiroDia := time.Date(dataFormatada.Year(), dataFormatada.Month(), 1, 0, 0, 0, 0, dataFormatada.Location())
	ultimoDia := primeiroDia.AddDate(0, 1, -1)
	//fmt.Println("enrei 1")
	// Conectar ao banco de dados local para buscar o nome da empresa
	dbLocal, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer dbLocal.Close()

	var nomeEmpresa string
	//	fmt.Println(cdEmp)
	query := "SELECT nome FROM EMPRESA"
	erro = dbLocal.QueryRow(query).Scan(&nomeEmpresa)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("Erro ao buscar o nome da empresa no banco local"))
		return
	}
	//fmt.Println("enrei 2")
	// Retornar as informações no JSON
	respostas.JSON(w, http.StatusOK, map[string]string{
		"data_sistema": dataSistema,
		"primeiro_dia": primeiroDia.Format("2006-01-02"),
		"ultimo_dia":   ultimoDia.Format("2006-01-02"),
		"cd_emp":       (cdEmp),
		"nome_empresa": nomeEmpresa,
	})
}

// AlterarDataSistema muda a data do sistema
func AlterarDataSistema(w http.ResponseWriter, r *http.Request) {
	var empresa modelos.Empresa

	// Decodificar o JSON recebido no corpo da requisição
	erro := json.NewDecoder(r.Body).Decode(&empresa)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Validar se a data foi enviada
	if empresa.DATA_SISTEMA == "" {
		respostas.Erro(w, http.StatusBadRequest, errors.New("Data do sistema não fornecida"))
		return
	}

	// Validar o formato da data
	novaData, erro := time.Parse("2006-01-02", empresa.DATA_SISTEMA)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("Formato de data inválido. Use YYYY-MM-DD"))
		return
	}

	// Extrair dados do usuário e empresa
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
	// Verificar nível de acesso do usuário
	nivelUsuario, erro := autenticacao.ExtrairUsuarioNivel(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	repositorioNiveis := repositorios.NovoRepositorioDeNiveis(db)
	nivelRequerido, erro := repositorioNiveis.BuscarOuCriarNivelPorCodigo("AlterarDataSistema")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if nivelUsuario < uint64(nivelRequerido) {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("nível de acesso insuficiente"))
		return
	}

	repositorio := repositorios.NovoRepositorioDeEmpresas(db)

	// Obter a data atual do sistema
	dataAtualStr, erro := repositorio.BuscarDataSistema()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	dataAtual, erro := time.Parse("2006-01-02", dataAtualStr)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("Erro ao interpretar a data do sistema"))
		return
	}

	// Validar que a nova data não avance mais de um dia
	if novaData.After(dataAtual.AddDate(0, 0, 1)) {
		respostas.Erro(w, http.StatusBadRequest, errors.New("A nova data do sistema não pode avançar mais de um dia em relação à data atual do sistema"))
		return
	}

	// Atualizar a data do sistema no banco
	erro = repositorio.AlterarDataSistema(empresa.DATA_SISTEMA)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, "Data do sistema alterada com sucesso!")
}
