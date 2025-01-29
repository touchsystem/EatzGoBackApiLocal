package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasVendas_add = []Rota{
	{
		URI:                "/vendas",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarVendas,
		RequerAutenticacao: true,
	},
	{
		URI:                "/vendas/{vendaId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarVendaPorID,
		RequerAutenticacao: true,
	},
	{
		URI:                "/vendas/{vendaId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarVenda,
		RequerAutenticacao: true,
	},
	{
		URI:                "/vendas/{vendaId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarVenda,
		RequerAutenticacao: true,
	},
}
