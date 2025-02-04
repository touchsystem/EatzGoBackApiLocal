package controllers

import (
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// BuscarContaPorMesa busca as vendas de uma mesa específica
func BuscarContaPorMesa(w http.ResponseWriter, r *http.Request) {
	// Extrair parâmetros da URL
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

	// Buscar parâmetro ID=2 (Taxa de Serviço)
	parametroTaxa, err := repositorio.BuscarParametroPorID(2)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Cálculo da taxa de serviço
	var taxaServico float64 = 0
	if parametroTaxa.Status == "S" && parametroTaxa.Limite > 0 {
		taxaServico = (totalConta * parametroTaxa.Limite) / 100
	}

	// Buscar parâmetro ID=12 (Conta Resumida)
	parametroResumida, err := repositorio.BuscarParametroPorID(12)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Se o parâmetro ID=12 for "S", agrupar produtos por CODM
	if parametroResumida.Status == "S" {
		vendas = agruparVendasPorCODM(vendas)
	}

	// Buscar parâmetro ID=22 (Exibir Nome da Empresa)
	parametroEmpresa, err := repositorio.BuscarParametroPorID(22)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Criar a estrutura base da resposta JSON
	conta := map[string]interface{}{
		"mesa_numero":  mesa.MesaCartao,
		"nick_usuario": nickUsuario,
		"nome_cliente": mesa.Apelido,
		"id_cliente":   mesa.IDCli,
		"celular":      mesa.Celular,
		"qtd_pessoas":  mesa.QtdPessoas,
		"taxa_servico": taxaServico,
		"total_conta":  totalConta,
		"vendas":       vendas,
	}

	// Se o parâmetro ID=22 for "S", buscar `nome_empresa` usando `BuscarDatasSistema`
	if parametroEmpresa.Status == "S" {
		nomeEmpresa, err := obterNomeEmpresa(r)
		if err == nil && nomeEmpresa != "" {
			conta["nome_empresa"] = nomeEmpresa
		}
	}

	// Enviar resposta JSON final
	respostas.JSON(w, http.StatusOK, conta)
}

// agruparVendasPorCODM agrupa os produtos com o mesmo CODM, somando as quantidades
func agruparVendasPorCODM(vendas []modelos.ContaVenda) []modelos.ContaVenda {
	agrupadas := make(map[string]modelos.ContaVenda)

	for _, venda := range vendas {
		if vendaExistente, existe := agrupadas[venda.CODM]; existe {
			// Se já existe, somamos a quantidade
			vendaExistente.QTD += venda.QTD
			agrupadas[venda.CODM] = vendaExistente
		} else {
			// Se não existe, adicionamos ao mapa
			agrupadas[venda.CODM] = modelos.ContaVenda{
				CODM: venda.CODM,
				DES2: venda.DES2,
				PV:   venda.PV,
				QTD:  venda.QTD,
			}
		}
	}

	// Converter mapa para slice
	var resultado []modelos.ContaVenda
	for _, venda := range agrupadas {
		resultado = append(resultado, venda)
	}

	return resultado
}

// obterNomeEmpresa chama `BuscarDatasSistema` e retorna apenas o nome da empresa
func obterNomeEmpresa(r *http.Request) (string, error) {
	wTemp := &responseWriterInterceptor{}
	BuscarDatasSistema(wTemp, r) // Chama a função existente

	// Retorna apenas o nome da empresa sem JSON extra
	return wTemp.nomeEmpresa, nil
}

// responseWriterInterceptor intercepta a resposta de `BuscarDatasSistema`
type responseWriterInterceptor struct {
	http.ResponseWriter
	nomeEmpresa string
}

func (rw *responseWriterInterceptor) Header() http.Header {
	return http.Header{}
}

func (rw *responseWriterInterceptor) Write(b []byte) (int, error) {
	// Processar o JSON retornado de `BuscarDatasSistema`
	var resposta map[string]interface{}
	err := json.Unmarshal(b, &resposta)
	if err == nil {
		if nome, ok := resposta["nome_empresa"].(string); ok {
			rw.nomeEmpresa = nome
		}
	}
	return len(b), nil // Evita chamar `WriteHeader` automaticamente
}

func (rw *responseWriterInterceptor) WriteHeader(statusCode int) {
	// Evita múltiplos `WriteHeader`
}
