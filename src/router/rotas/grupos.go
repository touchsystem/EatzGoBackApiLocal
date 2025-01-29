package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasGrupos = []Rota{
	{
		URI:                "/grupos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarGrupo,
		RequerAutenticacao: false,
	},
	{
		URI:                "/grupos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarGrupos,
		RequerAutenticacao: true,
	},
	{
		URI:                "/grupos/{grupoId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarGrupo,
		RequerAutenticacao: true,
	},
	{
		URI:                "/grupo-atualizar/{grupoId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarGrupo,
		RequerAutenticacao: true,
	},
	{
		URI:                "/grupo-deletar/{grupoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarGrupo,
		RequerAutenticacao: true,
	},
}
