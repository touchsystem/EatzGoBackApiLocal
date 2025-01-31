package sincronizacao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// Função de sincronização periódica para usuários
func SincronizarRegistroUsuariosPendentes() {
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
			sincronizarUsuariosDesdeNuvem()
			sincronizarUsuariosPendentes()
		}
	}
}

// Função de sincronização dos usuários pendentes com ID_HOSTWEB igual a 0
func sincronizarUsuariosPendentes() {
	tokenSync := os.Getenv("TOKEN_SYNC")

	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		log.Println("Erro ao conectar ao banco de dados:", erro)
		return
	}

	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.BuscarUsuariosNaoSincronizados()
	if erro != nil {
		log.Println("Erro ao buscar usuários não sincronizados:", erro)
		return
	}
	//fmt.Println(usuarios)
	for _, usuario := range usuarios {
		// Verificar se o usuário já existe na nuvem antes de sincronizar
		if repositorio.UsuarioExisteNaNuvem(usuario.ID) {
			log.Printf("Usuário com ID %d já existe na nuvem. Atualizando...\n", usuario.ID)
			repositorio.AtualizarUsuarioPorHostWeb(usuario.ID, usuario)
			continue
		}
		//fmt.Println("teste", usuario, tokenSync)
		idNuvem, erro := enviarUsuarioParaNuvem(usuario, tokenSync)
		if erro != nil {
			log.Println("Erro. . ao sincronizar usuário com a nuvem:", erro)
			continue
		}

		erro = repositorio.AtualizarUsuarioHostWeb(usuario.ID, idNuvem)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTWEB do usuário localmente:", erro)
		}
	}
}

// Função para enviar o usuário para a API da nuvem e receber o ID gerado
func enviarUsuarioParaNuvem(usuario modelos.Usuario, token string) (uint64, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	time.Sleep(10 * time.Second)
	url := fmt.Sprintf("http://%s:5000/usuario-sync", ipNuvem)

	corpoJSON, _ := json.Marshal(usuario)
	//fmt.Println("📤 JSON Enviado:", string(corpoJSON)) // <-- ADICIONADO

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

	//fmt.Println("🔄 Status da resposta da nuvem:", response.StatusCode) // <-- ADICIONADO

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body) // Lê a resposta da API para entender o erro
		fmt.Println("❌ Erro ao sincronizar usuário, resposta:", string(body))
		return 0, fmt.Errorf("erro ao sincronizar usuário com a nuvem: status %d", response.StatusCode)
	}

	var resultado map[string]uint64
	if erro := json.NewDecoder(response.Body).Decode(&resultado); erro != nil {
		fmt.Println("❌ Erro ao decodificar resposta da nuvem:", erro)
		return 0, erro
	}

	return resultado["id_usuario"], nil
}

// Função para sincronizar usuários desde a nuvem
func sincronizarUsuariosDesdeNuvem() {
	tokenSync := os.Getenv("TOKEN_SYNC")

	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		log.Println("Erro ao conectar ao banco de dados:", erro)
		return
	}
	defer db.Close()

	usuariosNuvem, erro := buscarUsuariosNaoSincronizadosNaNuvem(tokenSync)
	if erro != nil {
		log.Println("Erro ao buscar usuários não sincronizados na nuvem:", erro)
		return
	}

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	for _, usuario := range usuariosNuvem {
		// Verificar se o usuário já existe localmente antes de criar um novo

		if repositorio.UsuarioExisteNoLocal(usuario.ID_HOSTWEB) {
			log.Printf("Usuário com ID_HOSTWEB %d já existe localmente. Atualizando...\n", usuario.ID_HOSTWEB)
			repositorio.AtualizarUsuarioPorHostLocal(usuario.ID_HOSTWEB, usuario)
			continue
		}

		idLocal, erro := repositorio.CriarUsuarioSoLocal(usuario)
		if erro != nil {
			log.Println("Erro ao criar usuário local:", erro)
			continue
		}

		erro = atualizarUsuarioHostLocalNaNuvem(usuario.ID, idLocal, tokenSync)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTLOCAL na nuvem:", erro)
		} else {
			log.Printf("ID_HOSTLOCAL atualizado na nuvem para o usuário ID %d com valor %d\n", usuario.ID, idLocal)
		}

		erro = repositorio.AtualizarUsuarioHostWeb(idLocal, usuario.ID)
		if erro != nil {
			log.Println("Erro ao atualizar ID_HOSTWEB localmente:", erro)
		}
	}
}

// Função para buscar usuários na nuvem que ainda não foram sincronizados localmente
func buscarUsuariosNaoSincronizadosNaNuvem(token string) ([]modelos.Usuario, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/usuario-nao-sincronizados", ipNuvem)

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
		log.Printf("Erro ao buscar usuários na nuvem: Status %d, Resposta: %s", response.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("erro ao buscar usuários na nuvem: status %d", response.StatusCode)
	}

	var usuarios []modelos.Usuario
	if erro := json.NewDecoder(response.Body).Decode(&usuarios); erro != nil {
		log.Println("Erro ao decodificar resposta da nuvem:", erro)
		return nil, erro
	}

	return usuarios, nil
}

// Função para atualizar o ID_HOSTLOCAL do usuário na nuvem
func atualizarUsuarioHostLocalNaNuvem(idNuvem, idLocal uint64, token string) error {
	ipNuvem := os.Getenv("IP_NUVEM")
	url := fmt.Sprintf("http://%s:5000/atualizar-host-local-usuario", ipNuvem)

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

//rotinas de sync

// SincronizarUsuariosPeriodicamente realiza a sincronização periódica de usuários
func SincronizarUsuariosPeriodicamente() {
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

		sincronizarUsuariosParaNuvemPeriodico() // pega usuários aguardando sync local e envia para a nuvem
		sincronizarUsuariosDaNuvemParaLocalPeriodico()
	}
}

// SINCRONIZAÇÃO PERIÓDICA
func sincronizarUsuariosParaNuvemPeriodico() {
	tokenSync := os.Getenv("TOKEN_SYNC")
	ipNuvem := os.Getenv("IP_NUVEM")

	// Conectar ao banco de dados principal
	db, erro := banco.Conectar("DB_NOME")
	if erro != nil {
		log.Println("Erro ao conectar ao banco de dados:", erro)
		return
	}

	defer db.Close()
	// Cria o repositório de usuários e busca usuários com AGUARDANDO_SYNC = "S"
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.BuscarUsuariosAguardandoSync()
	if erro != nil {
		log.Println("Erro ao buscar usuários aguardando sincronização LOCAL:", erro)
		return
	}

	// Itera sobre os usuários aguardando sincronização
	for _, usuario := range usuarios {
		usuario.AGUARDANDO_SYNC = ""

		// Define a URL para atualizar o usuário na nuvem
		url := fmt.Sprintf("http://%s:5000/sync/usuarios/%d", ipNuvem, usuario.ID_HOSTWEB)
		corpoJSON, err := json.Marshal(usuario)
		if err != nil {
			log.Printf("Erro ao sincronizar usuário ID %d para JSON: %v", usuario.ID_HOSTWEB, err)
			continue
		}

		// Cria uma requisição PUT para enviar o usuário à nuvem
		request, err := http.NewRequest("PUT", url, bytes.NewBuffer(corpoJSON))
		if err != nil {
			log.Printf("Erro ao criar solicitação de atualização para usuário ID %d na nuvem: %v", usuario.ID_HOSTWEB, err)
			continue
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+tokenSync)

		// Envia a requisição
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Printf("Erro ao sincronizar usuário ID %d com a nuvem: %v", usuario.ID_HOSTWEB, err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusNoContent {
			log.Printf("Erro ao sincronizar usuário ID %d com a nuvem, status HTTP: %d", usuario.ID_HOSTWEB, response.StatusCode)
			continue
		}

		// Atualiza o usuário localmente para desmarcar AGUARDANDO_SYNC
		err = repositorio.DesmarcarAguardandoSyncUsuario(usuario.ID)
		if err != nil {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC para usuário ID %d: %v", usuario.ID, err)
		} else {
			log.Printf("Usuário ID %d sincronizado com sucesso.", usuario.ID)
		}
	}
}

// SINCRONIZAÇÃO DE USUÁRIOS DA NUVEM PARA O LOCAL
func sincronizarUsuariosDaNuvemParaLocalPeriodico() {
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

	// Solicita usuários aguardando sincronização
	url := fmt.Sprintf("http://%s:5000/usuarios-aguardando-sync", ipNuvem)

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
		log.Printf("Erro ao buscar usuários aguardando sincronização na nuvem: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Erro ao buscar usuários aguardando sincronização na nuvem, status: %d", response.StatusCode)
		return
	}

	var usuarios []modelos.Usuario
	if err = json.NewDecoder(response.Body).Decode(&usuarios); err != nil {
		log.Printf("Erro ao decodificar resposta: %v", err)
		return
	}

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	for _, usuario := range usuarios {
		err := repositorio.AtualizarUsuario(usuario.ID_HOSTLOCAL, usuario)
		if err != nil {
			log.Printf("Erro ao atualizar usuário ID_HOSTLOCAL %d localmente: %v", usuario.ID_HOSTLOCAL, err)
			continue
		}
		usuario.AGUARDANDO_SYNC = ""
		// Solicitação para desmarcar o AGUARDANDO_SYNC na nuvem
		urlDesmarcar := fmt.Sprintf("http://%s:5000/sync/desmarcarUsuario/%d", ipNuvem, usuario.ID)

		requestDesmarcar, err := http.NewRequest("PUT", urlDesmarcar, nil)
		if err != nil {
			log.Printf("Erro ao criar solicitação de desmarcação para usuário ID %d: %v", usuario.ID, err)
			continue
		}
		requestDesmarcar.Header.Set("Content-Type", "application/json")
		requestDesmarcar.Header.Set("Authorization", "Bearer "+tokenSync)

		responseDesmarcar, err := client.Do(requestDesmarcar)
		if err != nil {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC na nuvem para usuário ID %d, erro: %v", usuario.ID, err)
			continue
		}
		defer responseDesmarcar.Body.Close()

		if responseDesmarcar.StatusCode != http.StatusNoContent {
			log.Printf("Erro ao desmarcar AGUARDANDO_SYNC na nuvem para usuário ID %d, status: %d", usuario.ID, responseDesmarcar.StatusCode)
		} else {
			log.Printf("Usuário ID %d sincronizado com sucesso e desmarcado na nuvem.", usuario.ID)
		}
	}
}
