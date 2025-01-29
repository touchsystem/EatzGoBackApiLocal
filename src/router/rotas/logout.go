package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasLogout = Rota{
	URI:                "/logout",
	Metodo:             http.MethodPost,
	Funcao:             controllers.Deslogar,
	RequerAutenticacao: false,
}
