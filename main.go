package main

import (
	"api/src/autenticacao"
	"api/src/config"
	"api/src/router"
	"api/src/sincronizacao"
	"os"

	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	config.Carregar()
	r := router.Gerar()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://192.168.0.35:3000", "http://192.168.0.199:3000", "http://146.190.123.175:5000", "http://146.190.123.175:3000", "http://eatzgo.com", "https://eatzgo.com", "http://www.eatzgo.com", "https://wwww.eatzgo.com"}, // Especifica as origens permitidas
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Aplica o middleware CORS nas rotas
	handler := c.Handler(r)

	go sincronizacao.SincronizarRegristroClientesPendentes()
	go sincronizacao.SincronizarRegistroProdutosPendentes()
	go sincronizacao.SincronizarRegistroUsuariosPendentes()
	go sincronizacao.SincronizarRegistroGruposPendentes()

	go sincronizacao.SincronizarClientesPeriodicamente()
	go sincronizacao.SincronizarProdutosPeriodicamente()
	go sincronizacao.SincronizarUsuariosPeriodicamente()
	go sincronizacao.SincronizarGruposPeriodicamente()

	tokenSync := os.Getenv("TOKEN_SYNC")
	cdEmp, erro := autenticacao.ExtrairCdEmpDoTokenString(tokenSync)
	if erro != nil {
		log.Println("Erro ao extrair cdEmp do token:", erro)
		return
	}

	//fmt.Printf("2.4.1227  - API LOCAL Escutando na porta %d\n  -  ", config.Porta, cdEmp)
	fmt.Printf("2.5.0302c  -  API LOCAL Escutando na porta %d\n - Empresa: %s\n", config.Porta, cdEmp)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), handler))
}
