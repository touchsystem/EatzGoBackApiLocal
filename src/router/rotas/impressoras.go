package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasImpressoras = []Rota{

	// Buscar todas as impressoras
	{
		URI:                "/impressoras",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarImpressoras,
		RequerAutenticacao: true,
	},

	// Buscar impressora por ID
	{
		URI:                "/impressoras/{impressoraId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarImpressoraPorID,
		RequerAutenticacao: true,
	},

	// Atualizar impressora
	{
		URI:                "/impressoras/{impressoraId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarImpressora,
		RequerAutenticacao: true,
	},
}
