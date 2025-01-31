package rotas

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// Rota representa todas as rotas da API
type Rota struct {
	URI                string
	Metodo             string
	Funcao             func(http.ResponseWriter, *http.Request)
	RequerAutenticacao bool
}

// Configurar coloca todas as rotas dentro do router
func Configurar(r *mux.Router) *mux.Router {
	rotas := rotasUsuarios
	rotas = append(rotas, rotaLogin...)
	rotas = append(rotas, rotasClientes...)
	rotas = append(rotas, rotasProdutos...)
	rotas = append(rotas, rotasNiveis...)
	rotas = append(rotas, rotasLogout)
	rotas = append(rotas, rotasGrupos...)
	rotas = append(rotas, rotasCxReceb...)
	rotas = append(rotas, rotasImpressoras...)
	rotas = append(rotas, rotasVendas...)
	rotas = append(rotas, rotasParametros...)
	rotas = append(rotas, rotasMesas...)
	rotas = append(rotas, rotasEmpresa...)
	rotas = append(rotas, rotasVendas_add...)
	rotas = append(rotas, rotasTipoRecebimento...)
	rotas = append(rotas, rotasCaixa...)

	for _, rota := range rotas {

		if rota.RequerAutenticacao {
			r.HandleFunc(rota.URI,
				middlewares.Logger(middlewares.Autenticar(rota.Funcao)),
			).Methods(rota.Metodo)
		} else {
			r.HandleFunc(rota.URI, middlewares.Logger(rota.Funcao)).Methods(rota.Metodo)
		}

	}

	return r
}
