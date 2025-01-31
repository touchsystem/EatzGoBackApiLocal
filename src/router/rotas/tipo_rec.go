package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasTipoRecebimento = []Rota{
	{
		URI:                "/tipos_recebimento",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarTipoRecebimento,
		RequerAutenticacao: true,
	},
	{
		URI:                "/tipos_recebimento",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarTiposRecebimento,
		RequerAutenticacao: true,
	},
	{
		URI:                "/tipos_recebimento/{id}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarTipoRecebimento,
		RequerAutenticacao: true,
	},
}
