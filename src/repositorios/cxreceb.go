package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type CxRecebRepositorio struct {
	db *sql.DB
}

func NovoRepositorioCxReceb(db *sql.DB) *CxRecebRepositorio {
	return &CxRecebRepositorio{db}
}

// CriarCxReceb insere um registro na tabela CX_RECEB e retorna o ID gerado
func (repositorio *CxRecebRepositorio) CriarCxReceb(cxReceb modelos.CxReceb) (int, error) {
	query := `
		INSERT INTO CX_RECEB (ID_HOSTWEB, STATUS, DATA, ID_CLI, ID_USER, TOTAL, TROCO, MESA, NR_PESSOAS)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// No método CriarCxReceb
	dataFormatada := cxReceb.Data.Format("2006-01-02")

	// Executar a consulta
	result, err := repositorio.db.Exec(
		query,
		cxReceb.IDHostWeb,
		cxReceb.Status,
		dataFormatada, // Certifique-se de enviar no formato "YYYY-MM-DD"
		cxReceb.IDCli,
		cxReceb.IDUser,
		cxReceb.Total,
		cxReceb.Troco,
		cxReceb.Mesa,
		cxReceb.NRPessoas,
	)
	if err != nil {
		return 0, err
	}

	// Obter o ID gerado
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// CriarCxRecebTipos insere vários registros na tabela CX_RECEB_TIPO
func (repositorio *CxRecebRepositorio) CriarCxRecebTipos(cxRecebTipos []modelos.CxRecebTipo) error {
	query := `
		INSERT INTO CX_RECEB_TIPO (ID_CX_RECEB, ID_HOSTWEB, ID_TIPO_REC, MOEDA_NAC, MOEDA_EXT)
		VALUES (?, ?, ?, ?, ?)
	`

	// Inserir cada registro
	for _, tipo := range cxRecebTipos {
		_, err := repositorio.db.Exec(query, tipo.IDCxReceb, tipo.IDHostWeb, tipo.IDTipoRec, tipo.MoedaNac, tipo.MoedaExt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repositorio *CxRecebRepositorio) BuscarCxRecebPorID(idCxReceb int) (modelos.CxReceb, []modelos.CxRecebTipo, error) {
	var cxReceb modelos.CxReceb
	query := `SELECT * FROM CX_RECEB WHERE ID = ?`
	erro := repositorio.db.QueryRow(query, idCxReceb).Scan(&cxReceb.ID, &cxReceb.IDHostWeb, &cxReceb.Status, &cxReceb.Data, &cxReceb.IDCli, &cxReceb.IDUser, &cxReceb.Total, &cxReceb.Troco, &cxReceb.Mesa, &cxReceb.NRPessoas, &cxReceb.CriadoEm)
	if erro != nil {
		return cxReceb, nil, erro
	}

	var tipos []modelos.CxRecebTipo
	queryTipos := `SELECT * FROM CX_RECEB_TIPO WHERE ID_CX_RECEB = ?`
	rows, erro := repositorio.db.Query(queryTipos, idCxReceb)
	if erro != nil {
		return cxReceb, nil, erro
	}
	defer rows.Close()

	for rows.Next() {
		var tipo modelos.CxRecebTipo
		if erro := rows.Scan(&tipo.ID, &tipo.IDCxReceb, &tipo.IDHostWeb, &tipo.IDTipoRec, &tipo.MoedaNac, &tipo.MoedaExt); erro != nil {
			return cxReceb, nil, erro
		}
		tipos = append(tipos, tipo)
	}

	return cxReceb, tipos, nil
}
