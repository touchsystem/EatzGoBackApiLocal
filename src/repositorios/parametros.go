package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

// Parametros representa um repositório de parâmetros
type Parametros struct {
	db *sql.DB
}

// NovoRepositorioDeParametros cria um repositório de parâmetros
func NovoRepositorioDeParametros(db *sql.DB) *Parametros {
	return &Parametros{db}
}

// BuscarParametros retorna todos os parâmetros
func (repositorio Parametros) BuscarParametros() ([]modelos.Parametro, error) {
	rows, err := repositorio.db.Query("SELECT ID, NOME, STATUS, LIMITE FROM PARAMETROS ORDER BY NOME")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parametros []modelos.Parametro
	for rows.Next() {
		var parametro modelos.Parametro
		if err := rows.Scan(&parametro.ID, &parametro.Nome, &parametro.Status, &parametro.Limite); err != nil {
			return nil, err
		}
		parametros = append(parametros, parametro)
	}

	return parametros, nil
}

// BuscarParametroPorID retorna um parâmetro específico pelo ID
func (repositorio Parametros) BuscarParametroPorID(id uint64) (modelos.Parametro, error) {
	row := repositorio.db.QueryRow("SELECT ID, NOME, STATUS, LIMITE FROM PARAMETROS WHERE ID = ?", id)

	var parametro modelos.Parametro
	if err := row.Scan(&parametro.ID, &parametro.Nome, &parametro.Status, &parametro.Limite); err != nil {
		return modelos.Parametro{}, err
	}

	return parametro, nil
}

// AtualizarParametro atualiza os dados de um parâmetro existente
func (repositorio Parametros) AtualizarParametro(id uint64, parametro modelos.Parametro) error {
	statement, err := repositorio.db.Prepare("UPDATE PARAMETROS SET STATUS = ?, LIMITE = ? WHERE ID = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(parametro.Status, parametro.Limite, id)
	return err
}
