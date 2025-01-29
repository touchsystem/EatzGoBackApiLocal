package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotaLogin = []Rota{
	{
		URI:                "/login",
		Metodo:             http.MethodPost,
		Funcao:             controllers.Login,
		RequerAutenticacao: false,
	},
	{
		URI:                "/login-sync",
		Metodo:             http.MethodPost,
		Funcao:             controllers.Login_sync,
		RequerAutenticacao: false,
	},
	{
		URI:                "/me",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarDadosDoToken,
		RequerAutenticacao: false,
	},
}
