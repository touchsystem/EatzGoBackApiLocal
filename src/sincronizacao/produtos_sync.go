package sincronizacao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
)

// Sincronização periódica para produtos pendentes
func SincronizarRegistroProdutosPendentes() {
	tempoSyncStr := os.Getenv("TIME_SYNC")
	tempoSync, erro := time.ParseDuration(tempoSyncStr + "m")
	if erro != nil {
		log.Println("Erro ao ler TIME_SYNC:", erro)
		return
	}

	ticker := time.NewTicker(tempoSync)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sincronizarProdutosDesdeNuvem()
			sincronizarProdutosPendentes()
		}
	}
}

// Sincroniza os produtos pendentes com ID_HOSTWEB igual a 0
func sincronizarProdutosPendentes() {
	tokenSync := os.Getenv("TOKEN_SYNC")

	cdEmp, erro := autenticacao.ExtrairCdEmpDoTokenString(tokenSync)
	if erro != nil {
		log.Println("Erro ao extrair cdEmp do token:", erro)
		return
	}

	db, erro := banco.ConectarPorEmpresa(cdEmp)
	if erro != nil {
		log.Println("Erro ao conectar ao banco local:", erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtos, erro := repositorio.BuscarProdutosNaoSincronizados()
	if erro != nil {
		log.Println("Erro ao buscar produtos não sincronizados:", erro)
		return
	}

	for _, produto := range produtos {
		// Verifica se o produto já existe na nuvem antes de sincronizar
		if repositorio.ProdutoExisteNaNuvem(produto.ID) {
			log.Printf("Produto com ID %d já existe na nuvem. Atualizando...\n", produto.ID)
			repositorio.AtualizarProdutoPorHostWeb(produto.ID, produto)
			continue
		}

		idNuvem, erro := enviarProdutoParaNuvem(produto, tokenSync)
		if erro != nil {
			log.Println("Erro ao sincronizar produto com a nuvem:", erro)
			continue
		}

		erro = repositorio.AtualizarProdutoHostWeb(produto.ID, idNuvem)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTWEB do produto localmente:", erro)
		}
	}
}

// Envia o produto para a API da nuvem e recebe o ID gerado
func enviarProdutoParaNuvem(produto modelos.Produto, token string) (uint64, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/produto-sync", ipNuvem)

	corpoJSON, _ := json.Marshal(produto)
	request, erro := http.NewRequest("POST", url, bytes.NewBuffer(corpoJSON))
	if erro != nil {
		return 0, erro
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, erro := client.Do(request)
	if erro != nil {
		return 0, erro
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("erro ao sincronizar produto com a nuvem: status %d", response.StatusCode)
	}

	var resultado map[string]uint64
	if erro := json.NewDecoder(response.Body).Decode(&resultado); erro != nil {
		return 0, erro
	}

	return resultado["id_produto"], nil
}

// Sincroniza produtos desde a nuvem
func sincronizarProdutosDesdeNuvem() {
	tokenSync := os.Getenv("TOKEN_SYNC")
	cdEmp, erro := autenticacao.ExtrairCdEmpDoTokenString(tokenSync)
	if erro != nil {
		log.Println("Erro ao extrair cdEmp do token:", erro)
		return
	}

	db, erro := banco.ConectarPorEmpresa(cdEmp)
	if erro != nil {
		log.Println("Erro ao conectar ao banco local:", erro)
		return
	}
	defer db.Close()

	produtosNuvem, erro := buscarProdutosNaoSincronizadosNaNuvem(tokenSync)
	if erro != nil {
		log.Println("Erro ao buscar produtos não sincronizados na nuvem:", erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDeProdutos(db)

	for _, produto := range produtosNuvem {
		// Verifica se o produto já existe localmente antes de criar um novo
		if repositorio.ProdutoExisteNoLocal(uint64(produto.ID_HOSTWEB)) {
			log.Printf("Produto com ID_HOSTLOCAL %d já existe localmente. Atualizando...\n", produto.ID_HOSTLOCAL)
			repositorio.AtualizarProdutoPorHostLocal(uint64(produto.ID_HOSTWEB), produto)
			continue
		}

		idLocal, erro := repositorio.CriarProdutosSoLocal(produto)
		if erro != nil {
			log.Println("Erro ao criar produto local:", erro)
			continue
		}

		erro = atualizarProdutoHostLocalNaNuvem(produto.ID, idLocal, tokenSync)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTLOCAL na nuvem:", erro)
		} else {
			log.Printf("ID_HOSTLOCAL atualizado na nuvem para o produto ID %d com valor %d\n", produto.ID, idLocal)
		}

		erro = repositorio.AtualizarProdutoHostWeb(idLocal, produto.ID)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTWEB localmente:", erro)
		}
	}
}

// Busca produtos na nuvem que ainda não foram sincronizados localmente
func buscarProdutosNaoSincronizadosNaNuvem(token string) ([]modelos.Produto, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/produto-nao-sincronizados", ipNuvem)

	request, erro := http.NewRequest("POST", url, nil)
	if erro != nil {
		return nil, erro
	}
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, erro := client.Do(request)
	if erro != nil {
		return nil, erro
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		log.Printf("Erro ao buscar produtos na nuvem: Status %d, Resposta: %s", response.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("erro ao buscar produtos na nuvem: status %d", response.StatusCode)
	}

	var produtos []modelos.Produto
	if erro := json.NewDecoder(response.Body).Decode(&produtos); erro != nil {
		log.Println("Erro ao decodificar resposta da nuvem:", erro)
		return nil, erro
	}

	return produtos, nil
}

// Atualiza o ID_HOSTLOCAL do produto na nuvem
func atualizarProdutoHostLocalNaNuvem(idNuvem, idLocal uint64, token string) error {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/atualizar-host-local-produto", ipNuvem)

	corpoJSON, _ := json.Marshal(map[string]uint64{"id_nuvem": idNuvem, "id_local": idLocal})
	request, erro := http.NewRequest("PUT", url, bytes.NewBuffer(corpoJSON))
	if erro != nil {
		return erro
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, erro := client.Do(request)
	if erro != nil {
		return erro
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		log.Printf("Erro ao atualizar ID_HOSTLOCAL na nuvem: Status %d, Resposta: %s", response.StatusCode, string(bodyBytes))
		return fmt.Errorf("erro ao atualizar ID_HOSTLOCAL na nuvem: status %d", response.StatusCode)
	}

	return nil
}

// Sincronização periódica completa de produtos
func SincronizarProdutosPeriodicamente() {
	tempoSyncStr := os.Getenv("TIME_SYNC")
	tempoSync, erro := time.ParseDuration(tempoSyncStr + "m")
	if erro != nil {
		log.Println("Erro ao ler TIME_SYNC:", erro)
		return
	}

	ticker := time.NewTicker(tempoSync)
	defer ticker.Stop()

	for {
		<-ticker.C

		sincronizarProdutosParaNuvemPeriodico()
		sincronizarProdutosDaNuvemParaLocalPeriodico()
	}
}

// Sincronização periódica: envia produtos pendentes para a nuvem
func sincronizarProdutosParaNuvemPeriodico() {
	tokenSync := os.Getenv("TOKEN_SYNC")
	ipNuvem := os.Getenv("IP_NUVEM")

	cdEmp, erro := autenticacao.ExtrairCdEmpDoTokenString(tokenSync)
	if erro != nil {
		log.Println("Erro ao extrair cdEmp do token:", erro)
		return
	}

	db, erro := banco.ConectarPorEmpresa(cdEmp)
	if erro != nil {
		log.Println("Erro ao conectar ao banco local com cdEmp:", erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtos, erro := repositorio.BuscarProdutosAguardandoSync()
	if erro != nil {
		log.Println("Erro ao buscar produtos aguardando sincronização LOCAL:", erro)
		return
	}

	for _, produto := range produtos {
		produto.AGUARDANDO_SYNC = ""
		url := fmt.Sprintf("http://%s:5000/sync/produtos/%d", ipNuvem, produto.ID_HOSTWEB)

		corpoJSON, err := json.Marshal(produto)
		if err != nil {
			log.Printf("Erro ao serializar produto ID %d para JSON: %v", produto.ID_HOSTWEB, err)
			continue
		}

		request, err := http.NewRequest("PUT", url, bytes.NewBuffer(corpoJSON))
		if err != nil {
			log.Printf("Erro ao criar solicitação de atualização para produto ID %d na nuvem: %v", produto.ID_HOSTWEB, err)
			continue
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+tokenSync)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Printf("Erro B ao sincronizar produto ID %d com a nuvem: %v", produto.ID_HOSTWEB, err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNoContent {
			log.Printf("Erro  C ao sincronizar produto ID %d com a nuvem, status HTTP: %d", produto.ID_HOSTWEB, response.StatusCode)
			continue
		}

		err = repositorio.DesmarcarAguardandoSyncLocal(produto.ID)
		if err != nil {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC para produto ID %d: %v", produto.ID, err)
		}
	}
}

// Sincronização periódica: sincroniza produtos da nuvem para o local
func sincronizarProdutosDaNuvemParaLocalPeriodico() {
	tokenSync := os.Getenv("TOKEN_SYNC")
	ipNuvem := os.Getenv("IP_NUVEM")

	cdEmp, erro := autenticacao.ExtrairCdEmpDoTokenString(tokenSync)
	if erro != nil {
		log.Println("Erro ao extrair cdEmp do token:", erro)
		return
	}

	db, erro := banco.ConectarPorEmpresa(cdEmp)
	if erro != nil {
		log.Println("Erro ao conectar ao banco local com cdEmp:", erro)
		return
	}
	defer db.Close()

	url := fmt.Sprintf("http://%s:5000/produtos-aguardando-sync", ipNuvem)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Erro ao criar solicitação GET: %v", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+tokenSync)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Erro ao buscar produtos aguardando sincronização na nuvem: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Erro  A ao buscar produtos aguardando sincronização na nuvem, status: %d", response.StatusCode)
		return
	}

	var produtos []modelos.Produto
	if err = json.NewDecoder(response.Body).Decode(&produtos); err != nil {
		log.Printf("Erro ao decodificar resposta: %v", err)
		return
	}
	//fmt.Println(produtos)
	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	for _, produto := range produtos {
		err := repositorio.AtualizarProdutos(uint64(produto.ID_HOSTLOCAL), produto)
		if err != nil {
			log.Printf("Erro ao atualizar produto ID_HOSTLOCAL %d localmente: %v", produto.ID_HOSTLOCAL, err)
			continue
		}
		//	fmt.Println(produto.ID, " ", produto.ID_HOSTLOCAL, " ", produto.ID_HOSTWEB)
		urlDesmarcar := fmt.Sprintf("http://%s:5000/sync/desmarcarProduto/%d", ipNuvem, produto.ID)
		//	fmt.Println(produto.ID, " ", produto.ID_HOSTLOCAL, " ", produto.ID_HOSTWEB)
		requestDesmarcar, err := http.NewRequest("PUT", urlDesmarcar, nil)
		if err != nil {
			log.Printf("Erro ao criar solicitação de desmarcação para produto ID %d: %v", produto.ID, err)
			continue
		}
		requestDesmarcar.Header.Set("Content-Type", "application/json")
		requestDesmarcar.Header.Set("Authorization", "Bearer "+tokenSync)

		responseDesmarcar, err := client.Do(requestDesmarcar)
		if err != nil {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC na nuvem para produto ID %d, erro: %v", produto.ID, err)
			continue
		}
		defer responseDesmarcar.Body.Close()

		if responseDesmarcar.StatusCode != http.StatusNoContent {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC na nuvem para produto ID %d, status: %d", produto.ID, responseDesmarcar.StatusCode)
		}
	}
}
