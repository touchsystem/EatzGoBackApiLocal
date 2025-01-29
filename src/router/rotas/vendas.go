package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasVendas = []Rota{

	{
		URI:                "/grupos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarGrupos,
		RequerAutenticacao: true,
	},
	{
		URI:                "/produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProdutos,
		RequerAutenticacao: true,
	},
	{
		URI:                "/produtos/grupo/{grupoId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProdutosPorGrupo,
		RequerAutenticacao: true,
	},
	{
		URI:                "/nick",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarNick,
		RequerAutenticacao: true,
	},
	{
		URI:                "/preferidos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPreferidos,
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
		URI:                "/clientes",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarClientes,
		RequerAutenticacao: true,
	},
	{
		URI:                "/clientes/{clienteId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarCliente,
		RequerAutenticacao: true,
	},
}
