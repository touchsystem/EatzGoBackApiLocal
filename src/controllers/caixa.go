package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/repositorios"
	"api/src/respostas"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func BuscarVendasMesa(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	mesa, err := strconv.Atoi(parametros["mesa"])
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	cdEmp, err := autenticacao.ExtrairUsuarioCDEMP(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnauthorized, err)
		return
	}

	db, err := banco.ConectarPorEmpresa(cdEmp)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeCaixa(db)
	vendas, err := repositorio.BuscarVendasMesa(mesa)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, vendas)
}
