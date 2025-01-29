package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasCxReceb = []Rota{
	{
		URI:                "/cxreceb",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarCxReceb,
		RequerAutenticacao: true,
	},
	{
		URI:                "/cxreceb/{idCxReceb}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarCxReceb,
		RequerAutenticacao: true,
	},
}
