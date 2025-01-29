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

// Função de sincronização periódica para clientes
func SincronizarRegristroClientesPendentes() {
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
			sincronizarClientesDesdeNuvem() //ERRO ESTA AQUI
			sincronizarClientesPendentes()
		}
	}
}

// Função de sincronização dos clientes pendentes com ID_HOSTWEB igual a 0
func sincronizarClientesPendentes() {
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

	repositorio := repositorios.NovoRepositorioDeClientes(db)
	clientes, erro := repositorio.BuscarClientesNaoSincronizados()
	if erro != nil {
		log.Println("Erro ao buscar clientes não sincronizados:", erro)
		return
	}

	for _, cliente := range clientes {
		//	fmt.Println(cliente.STATUS)
		// Verificar se o cliente já existe na nuvem antes de sincronizar
		if repositorio.ClienteExisteNaNuvem(cliente.ID) {
			log.Printf("Cliente com ID %d já existe na nuvem. Atualizando...\n", cliente.ID)
			repositorio.AtualizarClientePorHostWeb(cliente.ID, cliente)
			continue
		}

		idNuvem, erro := enviarClienteParaNuvem(cliente, tokenSync)
		if erro != nil {
			log.Println("Erro ao sincronizar cliente com a nuvem:", erro)
			continue
		}

		erro = repositorio.AtualizarClienteHostWeb(cliente.ID, idNuvem)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTWEB do cliente localmente:", erro)
		}
	}
}

// Função para enviar o cliente para a API da nuvem e receber o ID gerado
func enviarClienteParaNuvem(cliente modelos.Cliente, token string) (uint64, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	time.Sleep(10 * time.Second)
	url := fmt.Sprintf("http://%s:5000/cliente-sync", ipNuvem)

	corpoJSON, _ := json.Marshal(cliente)
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
		return 0, fmt.Errorf("erro ao sincronizar cliente com a nuvem: status %d", response.StatusCode)
	}

	var resultado map[string]uint64
	if erro := json.NewDecoder(response.Body).Decode(&resultado); erro != nil {
		return 0, erro
	}

	return resultado["id_cliente"], nil
}

// Função para sincronizar clientes desde a nuvem
func sincronizarClientesDesdeNuvem() {
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
	//fmt.Println("cheque aqui clientes")
	//fmt.Println(tokenSync)
	clientesNuvem, erro := buscarClientesNaoSincronizadosNaNuvem(tokenSync)
	if erro != nil {
		log.Println("Erro ao buscar clientes não sincronizados na nuvem:", erro)
		return
	}
	//fmt.Println(clientesNuvem)

	repositorio := repositorios.NovoRepositorioDeClientes(db)

	for _, cliente := range clientesNuvem {
		// Verificar se o cliente já existe localmente antes de criar um novo

		if repositorio.ClienteExisteNoLocal(cliente.ID_HOSTWEB) {
			log.Printf("Cliente com ID_HOSTWEB %d já existe localmente. Atualizando...\n", cliente.ID_HOSTWEB)
			repositorio.AtualizarClientePorHostLocal(cliente.ID_HOSTWEB, cliente)
			continue
		}

		idLocal, erro := repositorio.CriarClienteSoLocal(cliente)
		if erro != nil {
			log.Println("Erro ao criar cliente local:", erro)
			continue
		}

		erro = atualizarClienteHostLocalNaNuvem(cliente.ID, idLocal, tokenSync)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTLOCAL na nuvem:", erro)
		} else {
			log.Printf("ID_HOSTLOCAL atualizado na nuvem para o cliente ID %d com valor %d\n", cliente.ID, idLocal)
		}

		erro = repositorio.AtualizarClienteHostWeb(idLocal, cliente.ID)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTWEB localmente:", erro)
		}
	}
}

// Função para buscar clientes na nuvem que ainda não foram sincronizados localmente
func buscarClientesNaoSincronizadosNaNuvem(token string) ([]modelos.Cliente, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/cliente-nao-sincronizados", ipNuvem)

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
		log.Printf("Erro ao buscar clientes na nuvem: Status %d, Resposta: %s", response.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("erro ao buscar clientes na nuvem: status %d", response.StatusCode)
	}

	var clientes []modelos.Cliente
	if erro := json.NewDecoder(response.Body).Decode(&clientes); erro != nil {
		log.Println("Erro ao decodificar resposta da nuvem:", erro)
		return nil, erro
	}

	return clientes, nil
}

// Função para atualizar o ID_HOSTLOCAL do cliente na nuvem
func atualizarClienteHostLocalNaNuvem(idNuvem, idLocal uint64, token string) error {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/atualizar-host-local-cliente", ipNuvem)

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

func SincronizarClientesPeriodicamente() {
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

		sincronizarClientesParaNuvemPeriodico() //pega cliente aguardando sync local e vai envia para nuvens
		sincronizarClientesDaNuvemParaLocalPeriodico()

	}
}

// SINCRONIZACAO PERIODICA
func sincronizarClientesParaNuvemPeriodico() {
	tokenSync := os.Getenv("TOKEN_SYNC")
	ipNuvem := os.Getenv("IP_NUVEM")

	// Extrai o código da empresa (cdEmp) do token para a conexão com o banco local
	cdEmp, erro := autenticacao.ExtrairCdEmpDoTokenString(tokenSync)
	if erro != nil {
		log.Println("Erro ao extrair cdEmp do token:", erro)
		return
	}

	// Conecta ao banco de dados local da empresa especificada pelo cdEmp
	db, erro := banco.ConectarPorEmpresa(cdEmp)
	if erro != nil {
		log.Println("Erro ao conectar ao banco local com cdEmp:", erro)
		return
	}
	defer db.Close()

	// Cria o repositório de clientes e busca clientes com AGUARDANDO_SYNC = "S"
	repositorio := repositorios.NovoRepositorioDeClientes(db)
	clientes, erro := repositorio.BuscarClientesAguardandoSync()
	if erro != nil {
		log.Println("Erro ao buscar clientes aguardando sincronização LOCAL:", erro)
		return
	}

	// Itera sobre os clientes aguardando sincronização
	for _, cliente := range clientes {
		cliente.AGUARDANDO_SYNC = ""

		// Define a URL para atualizar o cliente na nuvem
		url := fmt.Sprintf("http://%s:5000/sync/clientes/%d", ipNuvem, cliente.ID_HOSTWEB)
		//	log.Printf("Sincronizando cliente %d com a nuvem na URL: %s", cliente.ID_HOSTWEB, url)
		//	fmt.Println(cliente.ID_HOSTWEB)
		// Serializa o cliente para JSON
		corpoJSON, err := json.Marshal(cliente)
		if err != nil {
			log.Printf("Erro ao serializar cliente ID %d para JSON: %v", cliente.ID_HOSTWEB, err)
			continue
		}

		// Cria uma requisição PUT para enviar o cliente à nuvem
		request, err := http.NewRequest("PUT", url, bytes.NewBuffer(corpoJSON))
		if err != nil {
			log.Printf("Erro ao criar solicitação de atualização para cliente ID %d na nuvem: %v", cliente.ID_HOSTWEB, err)
			continue
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+tokenSync)

		// Envia a requisição
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Printf("Erro ao sincronizar cliente ID %d com a nuvem: %v", cliente.ID_HOSTWEB, err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNoContent {
			log.Printf("Erro ao sincronizar cliente ID %d com a nuvem, status HTTP: %d", cliente.ID_HOSTWEB, response.StatusCode)
			continue
		}

		// Atualiza o cliente localmente para desmarcar AGUARDANDO_SYNC
		err = repositorio.DesmarcarAguardandoSyncLocal(cliente.ID)
		if err != nil {
			//	log.Printf("Erro ao desmarcar AGUARDANDO_SYNC para cliente ID %d: %v", cliente.ID, err)
		} else {
			//log.Printf("Cliente ID %d sincronizado com sucesso.", cliente.ID)
		}
	}
}

// SINCRONIZACAO DE CLIENTES DA NUVEM PARA LOCAL
func sincronizarClientesDaNuvemParaLocalPeriodico() {
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

	// Solicita clientes aguardando sincronização
	url := fmt.Sprintf("http://%s:5000/clientes-aguardando-sync", ipNuvem)
	//	log.Printf("Buscando clientes aguardando sync na URL: %s", url)

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
		log.Printf("Erro ao buscar clientes aguardando sincronização na nuvem: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Erro ao buscar clientes aguardando sincronização na nuvem, status: %d", response.StatusCode)
		return
	}

	var clientes []modelos.Cliente
	if err = json.NewDecoder(response.Body).Decode(&clientes); err != nil {
		log.Printf("Erro ao decodificar resposta: %v", err)
		return
	}

	repositorio := repositorios.NovoRepositorioDeClientes(db)
	for _, cliente := range clientes {

		err := repositorio.AtualizarCliente(cliente.ID_HOSTLOCAL, cliente)
		if err != nil {
			log.Printf("Erro ao atualizar cliente ID_HOSTLOCAL %d localmente: %v", cliente.ID_HOSTLOCAL, err)
			continue
		}

		// Solicitação para desmarcar o AGUARDANDO_SYNC na nuvem
		urlDesmarcar := fmt.Sprintf("http://%s:5000/sync/desmarcarCliente/%d", ipNuvem, cliente.ID)
		//log.Printf("Desmarcando AGUARDANDO_SYNC na URL: %s", urlDesmarcar)

		requestDesmarcar, err := http.NewRequest("PUT", urlDesmarcar, nil)
		if err != nil {
			log.Printf("Erro ao criar solicitação de desmarcação para cliente ID %d: %v", cliente.ID, err)
			continue
		}
		requestDesmarcar.Header.Set("Content-Type", "application/json")
		requestDesmarcar.Header.Set("Authorization", "Bearer "+tokenSync)

		responseDesmarcar, err := client.Do(requestDesmarcar)
		if err != nil {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC na nuvem para cliente ID %d, erro: %v", cliente.ID, err)
			continue
		}
		defer responseDesmarcar.Body.Close()

		if responseDesmarcar.StatusCode != http.StatusNoContent {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC na nuvem para cliente ID %d, status: %d", cliente.ID, responseDesmarcar.StatusCode)
		} else {
			//	log.Printf("Cliente ID %d sincronizado com sucesso e desmarcado na nuvem.", cliente.ID)
		}
	}
}
