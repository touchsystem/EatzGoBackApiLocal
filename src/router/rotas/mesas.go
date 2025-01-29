package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasMesas = []Rota{
	{
		URI:                "/mesas",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarMesas,
		RequerAutenticacao: true,
	},
	{
		URI:                "/mesas/{mesaId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarMesa,
		RequerAutenticacao: true,
	},
	{
		URI:                "/mesas",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarMesas,
		RequerAutenticacao: true,
	},
	{
		URI:                "/mesas/{mesaId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarMesa,
		RequerAutenticacao: true,
	},
}
