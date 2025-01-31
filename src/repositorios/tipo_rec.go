package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type TiposRecebimento struct {
	db *sql.DB
}

func NovoRepositorioDeTiposRecebimento(db *sql.DB) *TiposRecebimento {
	return &TiposRecebimento{db}
}

func (repositorio TiposRecebimento) CriarTipoRecebimento(tipo modelos.TipoRecebimento) (uint64, error) {
	// Valida os dados antes de salvar
	if err := tipo.Validar(); err != nil {
		return 0, err
	}

	stmt, err := repositorio.db.Prepare("INSERT INTO TIPO_REC (NOME, CAMBIO, FT_CONV, STATUS) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(tipo.Nome, tipo.Cambio, tipo.FtConv, "A")
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (repositorio TiposRecebimento) BuscarTiposRecebimento() ([]modelos.TipoRecebimento, error) {
	rows, err := repositorio.db.Query("SELECT ID, NOME, CAMBIO, FT_CONV, STATUS FROM TIPO_REC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tiposRecebimento []modelos.TipoRecebimento
	for rows.Next() {
		var tipo modelos.TipoRecebimento
		if err := rows.Scan(&tipo.ID, &tipo.Nome, &tipo.Cambio, &tipo.FtConv, &tipo.Status); err != nil {
			return nil, err
		}
		tiposRecebimento = append(tiposRecebimento, tipo)
	}

	return tiposRecebimento, nil
}

func (repositorio TiposRecebimento) BuscarTipoRecebimentoPorID(id uint64) (modelos.TipoRecebimento, error) {
	var tipo modelos.TipoRecebimento
	err := repositorio.db.QueryRow("SELECT ID, NOME, CAMBIO, FT_CONV, STATUS FROM TIPO_REC WHERE ID = ?", id).
		Scan(&tipo.ID, &tipo.Nome, &tipo.Cambio, &tipo.FtConv, &tipo.Status)
	if err != nil {
		return modelos.TipoRecebimento{}, err
	}

	return tipo, nil
}

func (repositorio TiposRecebimento) AtualizarTipoRecebimento(id uint64, tipo modelos.TipoRecebimento) error {
	// Valida antes de atualizar
	if err := tipo.Validar(); err != nil {
		return err
	}

	_, err := repositorio.db.Exec(
		"UPDATE TIPO_REC SET NOME = ?, CAMBIO = ?, FT_CONV = ?, STATUS = ? WHERE ID = ?",
		tipo.Nome, tipo.Cambio, tipo.FtConv, tipo.Status, id,
	)
	return err
}
