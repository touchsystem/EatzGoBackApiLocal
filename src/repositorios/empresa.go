package repositorios

import (
	"database/sql"
)

// Empresa representa o repositório de empresas
type Empresa struct {
	db *sql.DB
}

// NovoRepositorioDeEmpresas cria um novo repositório de empresas
func NovoRepositorioDeEmpresas(db *sql.DB) *Empresa {
	return &Empresa{db}
}

// BuscarDataSistema busca a data do sistema na tabela EMPRESA
func (repositorio Empresa) BuscarDataSistema() (string, error) {
	var dataSistema string

	query := "SELECT DATA_SISTEMA FROM EMPRESA LIMIT 1"
	linha := repositorio.db.QueryRow(query)

	if erro := linha.Scan(&dataSistema); erro != nil {
		return "", erro
	}

	return dataSistema, nil
}

func (repositorio Empresa) AlterarDataSistema(dataSistema string) error {
	query := "UPDATE EMPRESA SET DATA_SISTEMA = ? WHERE ID_EMPRESA = 1" // Supondo que a empresa tenha ID fixo
	statement, erro := repositorio.db.Prepare(query)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(dataSistema)
	if erro != nil {
		return erro
	}

	return nil
}
