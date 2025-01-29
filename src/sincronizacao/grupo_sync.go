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

// Sincronização periódica para grupos pendentes
func SincronizarRegistroGruposPendentes() {
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
			sincronizarGruposDesdeNuvem()
			sincronizarGruposPendentes()
		}
	}
}

// Sincroniza os grupos pendentes com ID_HOSTWEB igual a 0
func sincronizarGruposPendentes() {
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

	repositorio := repositorios.NovoRepositorioDeGrupos(db)
	grupos, erro := repositorio.BuscarGruposNaoSincronizados()
	if erro != nil {
		log.Println("Erro ao buscar grupos não sincronizados:", erro)
		return
	}

	for _, grupo := range grupos {
		// Verifica se o grupo já existe na nuvem antes de sincronizar
		if repositorio.GrupoExisteNaNuvem(grupo.ID) {
			log.Printf("Grupo com ID %d já existe na nuvem. Atualizando...\n", grupo.ID)
			repositorio.AtualizarGrupoPorHostWeb(grupo.ID, grupo)
			continue
		}

		idNuvem, erro := enviarGrupoParaNuvem(grupo, tokenSync)
		if erro != nil {
			log.Println("Erro ao sincronizar grupo com a nuvem:", erro)
			continue
		}

		erro = repositorio.AtualizarGrupoHostWeb(grupo.ID, idNuvem)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTWEB do grupo localmente:", erro)
		}
	}
}

// Envia o grupo para a API da nuvem e recebe o ID gerado
func enviarGrupoParaNuvem(grupo modelos.Grupo, token string) (uint64, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/grupos-sync", ipNuvem)

	corpoJSON, _ := json.Marshal(grupo)
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
		return 0, fmt.Errorf("erro ao sincronizar grupo com a nuvem: status %d", response.StatusCode)
	}

	var resultado map[string]uint64
	if erro := json.NewDecoder(response.Body).Decode(&resultado); erro != nil {
		return 0, erro
	}

	return resultado["id_grupo"], nil
}

// Sincroniza grupos desde a nuvem
func sincronizarGruposDesdeNuvem() {
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

	gruposNuvem, erro := buscarGruposNaoSincronizadosNaNuvem(tokenSync)
	if erro != nil {
		log.Println("Erro ao buscar grupos não sincronizados na nuvem:", erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDeGrupos(db)

	for _, grupo := range gruposNuvem {
		// Verifica se o grupo já existe localmente antes de criar um novo
		if repositorio.GrupoExisteNoLocal(uint64(grupo.ID_HOSTWEB)) {
			log.Printf("Grupo com ID_HOSTLOCAL %d já existe localmente. Atualizando...\n", grupo.ID_HOSTLOCAL)
			repositorio.AtualizarGrupoPorHostLocal(uint64(grupo.ID_HOSTWEB), grupo)
			continue
		}

		idLocal, erro := repositorio.CriarGrupoSoLocal(grupo)
		if erro != nil {
			log.Println("Erro ao criar grupo local:", erro)
			continue
		}

		erro = atualizarGrupoHostLocalNaNuvem(grupo.ID, idLocal, tokenSync)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTLOCAL na nuvem:", erro)
		} else {
			log.Printf("ID_HOSTLOCAL atualizado na nuvem para o grupo ID %d com valor %d\n", grupo.ID, idLocal)
		}

		erro = repositorio.AtualizarGrupoHostWeb(idLocal, grupo.ID)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTWEB localmente:", erro)
		}
	}
}

// Busca grupos na nuvem que ainda não foram sincronizados localmente
func buscarGruposNaoSincronizadosNaNuvem(token string) ([]modelos.Grupo, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/grupos-nao-sincronizados", ipNuvem)

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
		log.Printf("Erro ao buscar grupos na nuvem: Status %d, Resposta: %s", response.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("erro ao buscar grupos na nuvem: status %d", response.StatusCode)
	}

	var grupos []modelos.Grupo
	if erro := json.NewDecoder(response.Body).Decode(&grupos); erro != nil {
		log.Println("Erro ao decodificar resposta da nuvem:", erro)
		return nil, erro
	}

	return grupos, nil
}

// Atualiza o ID_HOSTLOCAL do grupo na nuvem
func atualizarGrupoHostLocalNaNuvem(idNuvem, idLocal uint64, token string) error {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/atualizar-host-local-grupo", ipNuvem)

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

// Sincronização periódica completa de grupos
func SincronizarGruposPeriodicamente() {
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

		sincronizarGruposParaNuvemPeriodico()
		sincronizarGruposDaNuvemParaLocalPeriodico()
	}
}

// Sincronização periódica: envia grupos pendentes para a nuvem
func sincronizarGruposParaNuvemPeriodico() {
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

	repositorio := repositorios.NovoRepositorioDeGrupos(db)
	grupos, erro := repositorio.BuscarGruposAguardandoSync()
	if erro != nil {
		log.Println("Erro ao buscar grupos aguardando sincronização LOCAL:", erro)
		return
	}

	for _, grupo := range grupos {
		grupo.AGUARDANDO_SYNC = ""
		url := fmt.Sprintf("http://%s:5000/sync/grupos/%d", ipNuvem, grupo.ID_HOSTWEB)

		corpoJSON, err := json.Marshal(grupo)
		if err != nil {
			log.Printf("Erro ao serializar grupo ID %d para JSON: %v", grupo.ID_HOSTWEB, err)
			continue
		}

		request, err := http.NewRequest("PUT", url, bytes.NewBuffer(corpoJSON))
		if err != nil {
			log.Printf("Erro ao criar solicitação de atualização para grupo ID %d na nuvem: %v", grupo.ID_HOSTWEB, err)
			continue
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+tokenSync)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Printf("Erro ao sincronizar grupo ID %d com a nuvem: %v", grupo.ID_HOSTWEB, err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNoContent {
			log.Printf("Erro ao sincronizar grupo ID %d com a nuvem, status HTTP: %d", grupo.ID_HOSTWEB, response.StatusCode)
			continue
		}

		err = repositorio.DesmarcarAguardandoSyncLocal(grupo.ID)
		if err != nil {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC para grupo ID %d: %v", grupo.ID, err)
		}
	}
}

// Sincronização periódica: sincroniza grupos da nuvem para o local
func sincronizarGruposDaNuvemParaLocalPeriodico() {
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

	url := fmt.Sprintf("http://%s:5000/grupos-aguardando-sync", ipNuvem)

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
		log.Printf("Erro ao buscar grupos aguardando sincronização na nuvem: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Erro ao buscar grupos aguardando sincronização na nuvem, status: %d", response.StatusCode)
		return
	}

	var grupos []modelos.Grupo
	if err = json.NewDecoder(response.Body).Decode(&grupos); err != nil {
		log.Printf("Erro ao decodificar resposta: %v", err)
		return
	}

	repositorio := repositorios.NovoRepositorioDeGrupos(db)
	for _, grupo := range grupos {
		err := repositorio.AtualizarGrupo(uint64(grupo.ID_HOSTLOCAL), grupo)
		if err != nil {
			log.Printf("Erro ao atualizar grupo ID_HOSTLOCAL %d localmente: %v", grupo.ID_HOSTLOCAL, err)
			continue
		}

		urlDesmarcar := fmt.Sprintf("http://%s:5000/sync/desmarcarGrupo/%d", ipNuvem, grupo.ID)
		requestDesmarcar, err := http.NewRequest("PUT", urlDesmarcar, nil)
		if err != nil {
			log.Printf("Erro ao criar solicitação de desmarcação para grupo ID %d: %v", grupo.ID, err)
			continue
		}
		requestDesmarcar.Header.Set("Content-Type", "application/json")
		requestDesmarcar.Header.Set("Authorization", "Bearer "+tokenSync)

		responseDesmarcar, err := client.Do(requestDesmarcar)
		if err != nil {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC na nuvem para grupo ID %d, erro: %v", grupo.ID, err)
			continue
		}
		defer responseDesmarcar.Body.Close()

		if responseDesmarcar.StatusCode != http.StatusNoContent {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC na nuvem para grupo ID %d, status: %d", grupo.ID, responseDesmarcar.StatusCode)
		}
	}
}
