package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasNiveis = []Rota{

	{
		URI:                "/niveis",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarNiveis,
		RequerAutenticacao: true,
	},
	{
		URI:                "/niveis/{nivelId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarNivel,
		RequerAutenticacao: true,
	},
	{
		URI:                "/niveis/{nivelId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarNivel,
		RequerAutenticacao: true,
	},
}
