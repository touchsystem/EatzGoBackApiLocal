package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type Vendas struct {
	db *sql.DB
}

func NovoRepositorioDeVendas(db *sql.DB) *Vendas {
	return &Vendas{db}
}
func (repositorio Vendas) BuscarPorChave(chave string) ([]modelos.Venda, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT 
			ID, CODM, STATUS, IMPRESSORA, PV, STATUS_TP_VENDA, MESA, CELULAR, 
			CPF_CLIENTE, NOME_CLIENTE, ID_CLIENTE, QTD, OBS, DATA, ID_USER, NICK, CHAVE
		FROM VENDA WHERE CHAVE = ? ORDER BY IMPRESSORA`, chave)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var vendas []modelos.Venda
	for linhas.Next() {
		var venda modelos.Venda
		if erro = linhas.Scan(
			&venda.ID, &venda.CODM, &venda.STATUS, &venda.IMPRESSORA, &venda.PV, &venda.STATUS_TP_VENDA,
			&venda.MESA, &venda.CELULAR, &venda.CPF_CLIENTE, &venda.NOME_CLIENTE, &venda.ID_CLIENTE,
			&venda.QTD, &venda.OBS, &venda.DATA, &venda.ID_USER, &venda.NICK, &venda.CHAVE,
		); erro != nil {
			return nil, erro
		}
		vendas = append(vendas, venda)
	}

	return vendas, nil
}

func (repositorio Vendas) CriarVenda(venda modelos.Venda) (int, error) {

	statement, erro := repositorio.db.Prepare(`
		INSERT INTO VENDA (
			CODM, STATUS, IMPRESSORA, MONITOR_PRD, STATUS_TP_VENDA, STATUS_TP_HARD,
			MESA, CELULAR, CPF_CLIENTE, NOME_CLIENTE, ID_CLIENTE, PV, PV_PROM, QTD, 
			ID_USER, NICK, DATA, OBS, OBS2, OBS3, STATUS_PGTO, DATA_IFOOD, STATUS_IFOOD, 
			ID_IFOOD, COMPLEM_CODM, CHAVE
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(
		venda.CODM, venda.STATUS, venda.IMPRESSORA, venda.MONITOR_PRD, venda.STATUS_TP_VENDA, venda.STATUS_TP_HARD,
		venda.MESA, venda.CELULAR, venda.CPF_CLIENTE, venda.NOME_CLIENTE, venda.ID_CLIENTE, venda.PV, venda.PV_PROM, venda.QTD,
		venda.ID_USER, venda.NICK, venda.DATA, venda.OBS, venda.OBS2, venda.OBS3, venda.STATUS_PGTO,
		venda.DATA_IFOOD, venda.STATUS_IFOOD, venda.ID_IFOOD, venda.COMPLEM_CODM, venda.CHAVE,
	)
	if erro != nil {
		return 0, erro
	}

	idInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return int(idInserido), nil
}

func (repositorio Vendas) BuscarPorID(id int) (modelos.Venda, error) {
	linha := repositorio.db.QueryRow(`
		SELECT * FROM VENDA WHERE ID = ?
	`, id)

	var venda modelos.Venda
	if erro := linha.Scan(
		&venda.ID, &venda.CODM, &venda.STATUS, &venda.IMPRESSORA, &venda.MONITOR_PRD, &venda.STATUS_TP_VENDA, &venda.STATUS_TP_HARD,
		&venda.MESA, &venda.CELULAR, &venda.CPF_CLIENTE, &venda.NOME_CLIENTE, &venda.ID_CLIENTE, &venda.PV, &venda.PV_PROM, &venda.QTD,
		&venda.ID_USER, &venda.NICK, &venda.DATA, &venda.OBS, &venda.OBS2, &venda.OBS3, &venda.STATUS_PGTO, &venda.DATA_IFOOD, &venda.STATUS_IFOOD,
		&venda.ID_IFOOD, &venda.COMPLEM_CODM, &venda.CRIADOEM,
	); erro != nil {
		return modelos.Venda{}, erro
	}

	return venda, nil
}

func (repositorio Vendas) Atualizar(id int, venda modelos.Venda) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE VENDA SET 
			CODM = ?, STATUS = ?, IMPRESSORA = ?, MONITOR_PRD = ?, STATUS_TP_VENDA = ?, STATUS_TP_HARD = ?,
			MESA = ?, CELULAR = ?, CPF_CLIENTE = ?, NOME_CLIENTE = ?, ID_CLIENTE = ?, PV = ?, PV_PROM = ?, QTD = ?,
			ID_USER = ?, NICK = ?, DATA = ?, OBS = ?, OBS2 = ?, OBS3 = ?, STATUS_PGTO = ?, DATA_IFOOD = ?, STATUS_IFOOD = ?,
			ID_IFOOD = ?, COMPLEM_CODM = ?
		WHERE ID = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		venda.CODM, venda.STATUS, venda.IMPRESSORA, venda.MONITOR_PRD, venda.STATUS_TP_VENDA, venda.STATUS_TP_HARD,
		venda.MESA, venda.CELULAR, venda.CPF_CLIENTE, venda.NOME_CLIENTE, venda.ID_CLIENTE, venda.PV, venda.PV_PROM, venda.QTD,
		venda.ID_USER, venda.NICK, venda.DATA, venda.OBS, venda.OBS2, venda.OBS3, venda.STATUS_PGTO, venda.DATA_IFOOD, venda.STATUS_IFOOD,
		venda.ID_IFOOD, venda.COMPLEM_CODM, id,
	)
	return erro
}

func (repositorio Vendas) Deletar(id int) error {
	statement, erro := repositorio.db.Prepare(`DELETE FROM VENDA WHERE ID = ?`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(id)
	return erro
}
