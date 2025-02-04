package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasConta = []Rota{

	{
		URI:                "/conta",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarContaPorMesa,
		RequerAutenticacao: true,
	},
}
