package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasCaixa = []Rota{
	{
		URI:                "/caixa/vendas/{mesa}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarVendasMesa,
		RequerAutenticacao: true,
	},
}
