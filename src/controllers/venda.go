package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/impressao"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Função para criar vendas
// Função para criar vendas
func CriarVendas(w http.ResponseWriter, r *http.Request) {
	// Lê o corpo da requisição
	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	// Estrutura para receber o JSON
	var payload struct {
		Cabecalho struct {
			StatusTPVenda string `json:"status_tp_venda"`
			Mesa          int    `json:"mesa"`
			Celular       string `json:"celular"`
			CPFCliente    string `json:"cpf_cliente"`
			NomeCliente   string `json:"nome_cliente"`
			IDCliente     int    `json:"id_cliente"`
		} `json:"cabecalho"`
		Itens []struct {
			CODM string  `json:"codm"`
			QTD  float64 `json:"qtd"`
			OBS  string  `json:"obs"`
		} `json:"itens"`
	}

	// Decodifica o JSON recebido
	if erro = json.Unmarshal(corpoRequest, &payload); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Extrai o código da empresa do token do usuário
	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Conecta ao banco de dados da empresa
	db, erro := banco.ConectarPorEmpresa(cdEmp)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Valida se a mesa existe e verifica o status
	repositorioMesas := repositorios.NovoRepositorioDeMesas(db)
	mesas, erro := repositorioMesas.BuscarMesas("")
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	var mesaAtual *modelos.Mesa
	for _, mesa := range mesas {
		if mesa.ID == uint64(payload.Cabecalho.Mesa) {
			if mesa.Status != "L" && mesa.Status != "O" {
				respostas.Erro(w, http.StatusBadRequest, errors.New("mesa com status inválido. Apenas mesas 'L' (Livre) ou 'O' (Ocupado) são aceitas"))
				return
			}
			mesaAtual = &mesa
			break
		}
	}

	if mesaAtual == nil {
		respostas.Erro(w, http.StatusBadRequest, errors.New("mesa informada não existe"))
		return
	}

	// Valida se o cliente existe (se `IDCliente` for informado)
	if payload.Cabecalho.IDCliente > 0 {
		repositorioClientes := repositorios.NovoRepositorioDeClientes(db)
		_, erro := repositorioClientes.BuscarPorID(uint64(payload.Cabecalho.IDCliente))
		if erro != nil {
			respostas.Erro(w, http.StatusBadRequest, errors.New("cliente informado não existe"))
			return
		}
	}

	// DATA DO SISTEMA
	repositorioEmpresas := repositorios.NovoRepositorioDeEmpresas(db)
	dataSistema, erro := repositorioEmpresas.BuscarDataSistema()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Inicializa os repositórios
	repositorioProdutos := repositorios.NovoRepositorioDeProdutos(db)
	repositorioVendas := repositorios.NovoRepositorioDeVendas(db)
	repositorioParametros := repositorios.NovoRepositorioDeParametros(db)

	// Gera o ID do usuário e o nick do token
	idUsuarioUint64, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	idUsuario := int(idUsuarioUint64)

	nickUsuario, erro := autenticacao.ExtrairUsuarioNick(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Gera a chave para todas as vendas
	chave := fmt.Sprintf("%d_%d_%s", idUsuario, payload.Cabecalho.Mesa, time.Now().Format("2006-01-02 15:04:05"))

	// Verifica o valor do parâmetro com ID 23
	parametro23, erro := repositorioParametros.BuscarParametroPorID(23)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("Erro ao buscar o parâmetro ID 23"))
		return
	}

	// Verifica o valor do parâmetro com ID 24
	parametro24, erro := repositorioParametros.BuscarParametroPorID(24)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("Erro ao buscar o parâmetro ID 24"))
		return
	}

	// Verifica o valor do parâmetro com ID 25
	parametro25, erro := repositorioParametros.BuscarParametroPorID(25)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("Erro ao buscar o parâmetro ID 25"))
		return
	}
	parametro14, erro := repositorioParametros.BuscarParametroPorID(14)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, errors.New("Erro ao buscar o parâmetro ID 14"))
		return
	}

	//	fmt.Println(strings.TrimSpace(parametro23.Status), "23")

	//	fmt.Println(strings.TrimSpace(parametro24.Status), "24")
	//	fmt.Println(strings.TrimSpace(parametro25.Status), "25")

	imprimirVenda := strings.ToUpper(strings.TrimSpace(parametro23.Status)) == "S"
	// Determina se a impressao é agrupada
	//	agrupamentoImpressoras := strings.ToUpper(strings.TrimSpace(parametro24.Status)) == "N"

	// Processa os itens da venda
	var vendas []modelos.Venda
	for _, item := range payload.Itens {
		// Busca o produto pelo CODM
		produto, erro := repositorioProdutos.BuscarProdutoPorCODM(item.CODM)
		if erro != nil {
			respostas.Erro(w, http.StatusNotFound, errors.New("Produto não encontrado com o CODM fornecido"))
			return
		}

		// Cria a estrutura da venda
		venda := modelos.Venda{
			CODM:            item.CODM,
			STATUS:          "A",
			IMPRESSORA:      produto.IMPRE,
			PV:              produto.PV,
			STATUS_TP_VENDA: payload.Cabecalho.StatusTPVenda,
			MESA:            payload.Cabecalho.Mesa,
			CELULAR:         payload.Cabecalho.Celular,
			CPF_CLIENTE:     payload.Cabecalho.CPFCliente,
			NOME_CLIENTE:    payload.Cabecalho.NomeCliente,
			ID_CLIENTE:      payload.Cabecalho.IDCliente,
			QTD:             item.QTD,
			OBS:             item.OBS,
			DATA:            dataSistema,
			ID_USER:         idUsuario,
			NICK:            nickUsuario,
			CHAVE:           chave,
		}

		// Altera o status da mesa para 'O' (Ocupado), se estiver 'L' (Livre)
		if mesaAtual.Status == "L" {
			erro = repositorioMesas.AtualizarStatusMesa(uint64(venda.MESA), "O")
			if erro != nil {
				respostas.Erro(w, http.StatusInternalServerError, errors.New("não foi possível atualizar o status da mesa"))
				return
			}
		}

		// Insere a venda no banco de dados
		venda.ID, erro = repositorioVendas.CriarVenda(venda)
		if erro != nil {
			respostas.Erro(w, http.StatusInternalServerError, erro)
			return
		}

		vendas = append(vendas, venda)
	}

	// // começar a imprimir TICKET
	if imprimirVenda {

		// Verifica o tipo de impressão com base nos parâmetros
		if strings.ToUpper(strings.TrimSpace(parametro25.Status)) == "S" {
			// Impressão em fichas individuais (prioridade para parametro25)
			for _, venda := range vendas {

				// Declara iDvenda fora do if
				var iDvenda string

				if strings.ToUpper(strings.TrimSpace(parametro14.Status)) == "N" {
					iDvenda = "" // Atribui uma string vazia se parametro14.Status for "S"
				} else {
					iDvenda = strconv.Itoa(venda.ID) // Converte venda.ID para string
				}

				produto, _ := repositorioProdutos.BuscarProdutoPorCODM(venda.CODM)
				for i := 0; i < int(venda.QTD); i++ {
					fichaTexto := impressao.SetFontSize(2, 2) + fmt.Sprintf("%d  ", 1) + impressao.SetFontSize(1, 1) + fmt.Sprintf("%s %s \n", strings.TrimSpace(produto.DES2), iDvenda)
					if venda.OBS != "" {
						fichaTexto += fmt.Sprintf("Obs: %s\n", venda.OBS)
					}
					fichaTexto += "------------------------------------------\n"

					// Conversão de venda.MESA para string usando strconv.Itoa
					err := imprimirDetalhesDaVenda(&venda, fichaTexto, time.Now().Format("2006-01-02 15:04:05"), r, strconv.Itoa(venda.MESA))
					if err != nil {
						log.Printf("Erro ao imprimir ficha para venda %s: %v", venda.CODM, err)
					}
				}
			}
		} else if strings.ToUpper(strings.TrimSpace(parametro24.Status)) == "N" {
			// Impressão agrupada PARAMETRO 24 = N
			agrupadas := make(map[string]map[string]*modelos.Venda)

			for _, venda := range vendas {

				if agrupadas[venda.IMPRESSORA] == nil {
					agrupadas[venda.IMPRESSORA] = make(map[string]*modelos.Venda)
				}

				if itemExistente, existe := agrupadas[venda.IMPRESSORA][venda.CODM]; existe {
					itemExistente.QTD += venda.QTD
					if venda.OBS != "" {
						itemExistente.OBS += " | " + venda.OBS
					}
				} else {
					agrupadas[venda.IMPRESSORA][venda.CODM] = &modelos.Venda{
						CODM:            venda.CODM,
						STATUS:          venda.STATUS,
						IMPRESSORA:      venda.IMPRESSORA,
						PV:              venda.PV,
						STATUS_TP_VENDA: venda.STATUS_TP_VENDA,
						MESA:            venda.MESA,
						CELULAR:         venda.CELULAR,
						CPF_CLIENTE:     venda.CPF_CLIENTE,
						NOME_CLIENTE:    venda.NOME_CLIENTE,
						ID_CLIENTE:      venda.ID_CLIENTE,
						QTD:             venda.QTD,
						OBS:             venda.OBS,
						DATA:            venda.DATA,
						ID_USER:         venda.ID_USER,
						NICK:            venda.NICK,
						CHAVE:           venda.CHAVE,
					}
				}
			}

			// Imprime as vendas agrupadas
			for impressora, itens := range agrupadas {
				textoImpressao := impressao.SetCharacterSetCP850()

				for _, venda := range itens {
					quantidadeTexto := ""
					if math.Mod(venda.QTD, 1) == 0 {
						quantidadeTexto = strconv.Itoa(int(venda.QTD))
					} else {
						quantidadeTexto = fmt.Sprintf("%.3f", venda.QTD)
					}

					produto, _ := repositorioProdutos.BuscarProdutoPorCODM(venda.CODM)
					textoImpressao += impressao.SetFontSize(2, 2) + quantidadeTexto + "  " + impressao.SetFontSize(1, 1) + produto.DES2 + "\n"
					if venda.OBS != "" {
						textoImpressao += fmt.Sprintf("Obs: %s\n", venda.OBS)
					}
				}

				for _, vendaExemplo := range itens {
					// Conversão de vendaExemplo.MESA para string usando strconv.Itoa
					err := imprimirDetalhesDaVenda(vendaExemplo, textoImpressao, dataSistema, r, strconv.Itoa(vendaExemplo.MESA))
					if err != nil {
						log.Printf("Erro ao imprimir para impressora %s: %v", impressora, err)
					}
					break
				}
			}
			//FIM PARAMETRO 24 = N

		} else if strings.ToUpper(strings.TrimSpace(parametro24.Status)) == "S" {
			// PARAMETRO 24 = S

			// Percorre as vendas para impressão individual
			for _, venda := range vendas {
				textoImpressao := impressao.SetCharacterSetCP850()

				// Declara iDvenda fora do if
				var iDvenda string

				if strings.ToUpper(strings.TrimSpace(parametro14.Status)) == "N" {
					iDvenda = "" // Atribui uma string vazia se parametro14.Status for "S"
				} else {
					iDvenda = strconv.Itoa(venda.ID) // Converte venda.ID para string
				}
				// Define a quantidade como texto
				quantidadeTexto := ""
				if math.Mod(venda.QTD, 1) == 0 {
					quantidadeTexto = strconv.Itoa(int(venda.QTD))
				} else {
					quantidadeTexto = fmt.Sprintf("%.3f", venda.QTD)
				}

				// Busca detalhes do produto
				produto, _ := repositorioProdutos.BuscarProdutoPorCODM(venda.CODM)
				textoImpressao += impressao.SetFontSize(2, 2) + quantidadeTexto + "  " + impressao.SetFontSize(1, 1) + strings.TrimSpace(produto.DES2) + " " + iDvenda + "\n"
				if venda.OBS != "" {
					textoImpressao += fmt.Sprintf("Obs: %s\n", venda.OBS)
				}

				// Adiciona separador para cada item
				textoImpressao += "------------------------------------------\n"

				// Chamada da função de impressão com o número da mesa convertido para string
				err := imprimirDetalhesDaVenda(&venda, textoImpressao, dataSistema, r, strconv.Itoa(venda.MESA))
				if err != nil {
					log.Printf("Erro ao imprimir para impressora %s: %v", venda.IMPRESSORA, err)
				}
			}
			// FIM PARAMETRO 24 = S
		}

	}

	// Retorna o array de vendas criadas como resposta
	respostas.JSON(w, http.StatusCreated, vendas)
}

// Função para imprimir detalhes da venda
func imprimirDetalhesDaVenda(venda *modelos.Venda, textoImpressao string, dataSistema string, r *http.Request, mesa string) error {
	if venda == nil || venda.IMPRESSORA == "" {
		return fmt.Errorf("venda inválida ou impressora não informada")
	}

	// Extrai o código da empresa do token do usuário
	cdEmp, erro := autenticacao.ExtrairUsuarioCDEMP(r)
	if erro != nil {
		return fmt.Errorf("erro ao extrair código da empresa: %v", erro)
	}

	// Conecta ao banco de dados da empresa
	db, erro := banco.ConectarPorEmpresa(cdEmp)
	if erro != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", erro)
	}
	defer db.Close()

	// Inicializar o repositório de impressoras
	repositorioImpressoras := repositorios.NovoRepositorioDeImpressoras(db)

	// Buscar o endereço da impressora pelo código
	impressora, erro := repositorioImpressoras.BuscarImpressoraPorCODIMP(venda.IMPRESSORA)
	if erro != nil {
		return fmt.Errorf("erro ao buscar impressora para código %s: %v", venda.IMPRESSORA, erro)
	}

	// Cabeçalho e detalhes da venda para impressão
	texto := impressao.SetCharacterSetCP850() // Configura o conjunto de caracteres
	texto += "------------------------------------------\n"
	texto += impressao.SetBold(true) + "VENDA - DETALHES\n" + impressao.SetBold(false)
	texto += fmt.Sprintf("Ponto Produção: %s %s\n", venda.IMPRESSORA, strings.TrimSpace(impressora.END_IMP))
	texto += fmt.Sprintf("Data/Hora: %s\n", dataSistema)
	texto += fmt.Sprintf("Mesa: %s\n", mesa) // Inclui o número da mesa
	texto += "------------------------------------------\n"
	texto += textoImpressao + "\n" // Adiciona os detalhes já formatados do texto de impressão
	texto += "------------------------------------------\n"
	texto += impressao.CutPaper()
	texto += impressao.ResetPrinter()

	// Configuração do endereço completo da impressora
	enderecoImpressora := fmt.Sprintf("\\\\%s\\%s", strings.TrimSpace(impressora.END_SER), strings.TrimSpace(impressora.END_IMP))

	// Envia o texto para a impressora
	erro = impressao.PrintToPrinter(enderecoImpressora, texto)
	if erro != nil {
		return fmt.Errorf("erro ao enviar para a impressora no endereço %s: %v", enderecoImpressora, erro)
	}

	//	log.Printf("Impressão enviada com sucesso para a impressora no endereço: %s", enderecoImpressora)
	return nil
}

// BuscarVendaPorID
func BuscarVendaPorID(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	id, erro := strconv.Atoi(parametros["vendaId"])
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

	repositorio := repositorios.NovoRepositorioDeVendas(db)
	venda, erro := repositorio.BuscarPorID(id)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, venda)
}
func AtualizarVenda(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	id, erro := strconv.Atoi(parametros["vendaId"])
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var venda modelos.Venda
	if erro = json.Unmarshal(corpoRequest, &venda); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = venda.Preparar(); erro != nil {
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

	repositorio := repositorios.NovoRepositorioDeVendas(db)
	if erro = repositorio.Atualizar(id, venda); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

func DeletarVenda(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	id, erro := strconv.Atoi(parametros["vendaId"])
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

	repositorio := repositorios.NovoRepositorioDeVendas(db)
	if erro = repositorio.Deletar(id); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
