package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// SincronizarCambioNuvem envia os dados do câmbio para o servidor na nuvem e retorna o ID gerado nas nuvens
func SincronizarCambioNuvem(cambio interface{}) (int, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	tokenSync := os.Getenv("TOKEN_SYNC")

	if ipNuvem == "" || tokenSync == "" {
		return 0, fmt.Errorf("Configurações de IP_NUVEM ou TOKEN_SYNC ausentes")
	}

	url := fmt.Sprintf("http://%s:5000/cambio-sync", ipNuvem)

	jsonData, err := json.Marshal(cambio)
	if err != nil {
		return 0, fmt.Errorf("Erro ao converter dados para JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("Erro ao criar solicitação de sincronização: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenSync)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Erro ao enviar solicitação de sincronização: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Falha na sincronização: código de status %d", resp.StatusCode)
	}

	// Lê a resposta e extrai o ID_CAMBIO da resposta JSON
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Erro ao ler resposta: %v", err)
	}

	var resposta struct {
		ID_CAMBIO int `json:"id_cambio"`
	}
	err = json.Unmarshal(body, &resposta)
	if err != nil {
		return 0, fmt.Errorf("Erro ao parsear ID_CAMBIO da resposta: %v", err)
	}

	return resposta.ID_CAMBIO, nil
}

// SincronizarClienteNuvem envia os dados do cliente para o servidor na nuvem e retorna o ID gerado nas nuvens
func SincronizarClienteNuvem(cliente interface{}) (uint64, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	tokenSync := os.Getenv("TOKEN_SYNC")

	if ipNuvem == "" || tokenSync == "" {
		return 0, fmt.Errorf("Configurações de IP_NUVEM ou TOKEN_SYNC ausentes")
	}

	url := fmt.Sprintf("http://%s:5000/cliente-sync", ipNuvem)

	jsonData, err := json.Marshal(cliente)
	if err != nil {
		return 0, fmt.Errorf("Erro ao converter dados para JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("Erro ao criar solicitação de sincronização: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenSync)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Erro ao enviar solicitação de sincronização: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Falha na sincronização: código de status %d", resp.StatusCode)
	}

	// Lê a resposta e extrai o ID_CLIENTE da resposta JSON
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Erro ao ler resposta: %v", err)
	}

	var resposta struct {
		ID_CLIENTE uint64 `json:"id_cliente"`
	}
	err = json.Unmarshal(body, &resposta)
	if err != nil {
		return 0, fmt.Errorf("Erro ao parsear ID do cliente da resposta: %v", err)
	}

	return resposta.ID_CLIENTE, nil
}

// SincronizarProdutoNuvem envia os dados do produto para o servidor na nuvem e retorna o ID gerado na nuvem
func SincronizarProdutoNuvem(produto interface{}) (uint64, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	tokenSync := os.Getenv("TOKEN_SYNC")

	if ipNuvem == "" || tokenSync == "" {
		return 0, fmt.Errorf("Configurações de IP_NUVEM ou TOKEN_SYNC ausentes")
	}

	url := fmt.Sprintf("http://%s:5000/produto-sync", ipNuvem)

	// Converte os dados do produto para JSON
	jsonData, err := json.Marshal(produto)
	if err != nil {
		return 0, fmt.Errorf("Erro ao converter dados do produto para JSON: %v", err)
	}

	// Cria a solicitação HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("Erro ao criar solicitação de sincronização: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenSync)

	// Envia a solicitação e trata a resposta
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Erro ao enviar solicitação de sincronização: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Falha na sincronização: código de status %d", resp.StatusCode)
	}

	// Lê a resposta e extrai o ID_PRODUTO da resposta JSON
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Erro ao ler resposta: %v", err)
	}

	var resposta struct {
		ID_PRODUTO uint64 `json:"id_produto"`
	}
	err = json.Unmarshal(body, &resposta)
	if err != nil {
		return 0, fmt.Errorf("Erro ao parsear ID do produto da resposta: %v", err)
	}

	return resposta.ID_PRODUTO, nil
}

// SincronizarGrupoNuvem envia os dados do grupo para o servidor na nuvem e retorna o ID gerado na nuvem
func SincronizarGrupoNuvem(grupo interface{}) (uint64, error) {
	ipNuvem := os.Getenv("IP_NUVEM")
	tokenSync := os.Getenv("TOKEN_SYNC")

	if ipNuvem == "" || tokenSync == "" {
		return 0, fmt.Errorf("Configurações de IP_NUVEM ou TOKEN_SYNC ausentes")
	}

	url := fmt.Sprintf("http://%s:5000/grupo-sync", ipNuvem)

	// Converte os dados do grupo para JSON
	jsonData, err := json.Marshal(grupo)
	if err != nil {
		return 0, fmt.Errorf("Erro ao converter dados do grupo para JSON: %v", err)
	}

	// Cria a solicitação HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("Erro ao criar solicitação de sincronização: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenSync)

	// Envia a solicitação e trata a resposta
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Erro ao enviar solicitação de sincronização: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Falha na sincronização: código de status %d", resp.StatusCode)
	}

	// Lê a resposta e extrai o ID_GRUPO da resposta JSON
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Erro ao ler resposta: %v", err)
	}

	var resposta struct {
		ID_GRUPO uint64 `json:"id_grupo"`
	}
	err = json.Unmarshal(body, &resposta)
	if err != nil {
		return 0, fmt.Errorf("Erro ao parsear ID do grupo da resposta: %v", err)
	}

	return resposta.ID_GRUPO, nil
}
