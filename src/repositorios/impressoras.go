package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

// Impressoras representa o repositório de impressoras
type Impressoras struct {
	db *sql.DB
}

// NovoRepositorioDeImpressoras cria um novo repositório de impressoras
func NovoRepositorioDeImpressoras(db *sql.DB) *Impressoras {
	return &Impressoras{db}
}

// BuscarImpressoras recupera todas as impressoras com filtros opcionais
func (repositorio Impressoras) BuscarImpressoras(status string, filtro string) ([]modelos.Impressora, error) {
	query := `SELECT ID,  COD_IMP, NOME, END_IMP, END_SER  FROM IMPRESSORAS WHERE 1=1`
	var args []interface{}

	// Ordenar os resultados
	query += " ORDER BY COD_IMP"

	rows, erro := repositorio.db.Query(query, args...)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var impressoras []modelos.Impressora
	for rows.Next() {
		var impressora modelos.Impressora
		if erro := rows.Scan(
			&impressora.ID,
			&impressora.COD_IMP,
			&impressora.NOME,
			&impressora.END_IMP,
			&impressora.END_SER,
		); erro != nil {
			return nil, erro
		}
		impressoras = append(impressoras, impressora)
	}

	return impressoras, nil
}

// BuscarImpressoraPorID recupera uma impressora específica pelo ID
func (repositorio Impressoras) BuscarImpressoraPorID(impressoraID int) (modelos.Impressora, error) {
	query := `SELECT ID, COD_IMP, NOME, END_IMP, END_SER  FROM IMPRESSORAS WHERE ID = ?`
	row := repositorio.db.QueryRow(query, impressoraID)

	var impressora modelos.Impressora
	if erro := row.Scan(
		&impressora.ID,
		&impressora.COD_IMP,
		&impressora.NOME,
		&impressora.END_IMP,
		&impressora.END_SER,
	); erro != nil {
		return modelos.Impressora{}, erro
	}
	//fmt.Println(impressoraBusca)
	return impressora, nil
}

// AtualizarImpressora altera uma impressora no banco de dados
func (repositorio Impressoras) AtualizarImpressora(impressoraID uint64, impressora modelos.Impressora) error {
	statement, erro := repositorio.db.Prepare(
		`UPDATE IMPRESSORAS 
		 SET 
		    AGUARDANDO_SYNC = ?, 
		
		    NOME = ?,
			END_IMP = ?,
			END_SER = ?

		 WHERE ID = ?`,
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		impressora.AGUARDANDO_SYNC,

		impressora.NOME,
		impressora.END_IMP,
		impressora.END_SER,

		impressoraID,
	)
	return erro
}

// BuscarImpressoraPorCODIMP busca uma impressora no banco de dados pelo COD_IMP
func (repositorio Impressoras) BuscarImpressoraPorCODIMP(codImp string) (modelos.Impressora, error) {
	query := `SELECT ID, COD_IMP, NOME, END_IMP, END_SER FROM IMPRESSORAS WHERE TRIM(COD_IMP) = ?`
	row := repositorio.db.QueryRow(query, codImp)

	var impressora modelos.Impressora
	if erro := row.Scan(
		&impressora.ID,
		&impressora.COD_IMP,
		&impressora.NOME,
		&impressora.END_IMP,
		&impressora.END_SER,
	); erro != nil {
		return modelos.Impressora{}, erro
	}

	return impressora, nil
}
