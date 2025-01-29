package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasParametros = []Rota{

	{
		URI:                "/parametros",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarParametros,
		RequerAutenticacao: true,
	},
	{
		URI:                "/parametros/{parametroId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarParametro,
		RequerAutenticacao: true,
	},
	{
		URI:                "/parametros/{parametroId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarParametro,
		RequerAutenticacao: true,
	},
}
