package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasClientes = []Rota{
	{
		URI:                "/clientes",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarCliente,
		RequerAutenticacao: false,
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
	{
		URI:                "/clientes/{clienteId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarCliente,
		RequerAutenticacao: true,
	},
	{
		URI:                "/clientes/{clienteId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarCliente,
		RequerAutenticacao: true,
	},
	{
		URI:                "/sync/clientes/{clienteId}", // Nova rota para sincronização
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarCliente,
		RequerAutenticacao: true,
	},
}
