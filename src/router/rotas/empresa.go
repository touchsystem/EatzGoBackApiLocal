package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasEmpresa = []Rota{
	{
		URI:                "/datas-sistema",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarDatasSistema,
		RequerAutenticacao: true,
	},
	{
		URI:                "/datas-sistema",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AlterarDataSistema,
		RequerAutenticacao: true,
	},
}
