package controllers

import (
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"errors"
	"net/http"
	"strconv"
)

// BuscarContaPorMesa busca as vendas de uma mesa específica
func BuscarContaPorMesa(w http.ResponseWriter, r *http.Request) {
	// Extrair parâmetros
	queryParams := r.URL.Query()
	numMesa, err := strconv.Atoi(queryParams.Get("mesa"))
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("Número da mesa inválido"))
		return
	}

	nickUsuario := queryParams.Get("nick")
	if nickUsuario == "" {
		respostas.Erro(w, http.StatusBadRequest, errors.New("O campo 'nick' é obrigatório"))
		return
	}

	// Conectar ao banco
	db, err := banco.Conectar("DB_NOME")
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeContas(db)

	// Buscar informações da mesa
	mesa, err := repositorio.BuscarMesaPorNumero(numMesa)
	if err != nil {
		respostas.Erro(w, http.StatusNotFound, errors.New("Mesa não encontrada"))
		return
	}
	// Buscar vendas associadas à mesa e calcular total
	vendas, totalConta, err := repositorio.BuscarVendasPorMesa(numMesa)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Montar resposta com o total da conta incluído
	conta := modelos.Conta{
		MesaNumero:  mesa.MesaCartao,
		NickUsuario: nickUsuario,
		NomeCliente: mesa.Apelido,
		IDCliente:   &mesa.IDCli,
		Celular:     mesa.Celular,
		QtdPessoas:  mesa.QtdPessoas,
		Vendas:      vendas,
		TotalConta:  totalConta, // Adiciona o total da conta corretamente
	}

	respostas.JSON(w, http.StatusOK, conta)
}
