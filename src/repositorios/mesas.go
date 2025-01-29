package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"errors"
)

type Mesas struct {
	db *sql.DB
}

func NovoRepositorioDeMesas(db *sql.DB) *Mesas {
	return &Mesas{db}
}

func (repositorio Mesas) AtualizarStatusMesa(id uint64, status string) error {
	_, erro := repositorio.db.Exec(`UPDATE MESA SET STATUS = ? WHERE MESA_CARTAO = ?`, status, id)
	return erro
}

// VerificarSeMesaExiste verifica se a mesa já existe pelo número
func (repositorio Mesas) VerificarSeMesaExiste(mesaCartao int) (bool, error) {
	var id uint64
	err := repositorio.db.QueryRow("SELECT ID FROM MESA WHERE MESA_CARTAO = ?", mesaCartao).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CriarMesa insere uma nova mesa no banco de dados
func (repositorio Mesas) CriarMesa(mesa modelos.Mesa) (uint64, error) {
	statement, err := repositorio.db.Prepare(`
		INSERT INTO MESA (MESA_CARTAO, STATUS) 
		VALUES (?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	resultado, err := statement.Exec(
		mesa.MesaCartao,
		mesa.Status,
	)
	if err != nil {
		return 0, err
	}

	id, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

// AtualizarMesa atualiza uma mesa no banco de dados
func (repositorio Mesas) AtualizarMesa(id uint64, mesa modelos.Mesa) error {
	statement, err := repositorio.db.Prepare(`
		UPDATE MESA 
		SET  STATUS = ?, ID_USER = ?, ID_CLI = ?, NICK = ?,  ABERTURA = ?, 
		    QTD_PESSOAS = ?, TURISTA = ?, CELULAR = ?, APELIDO = ? 
		WHERE ID = ?
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(

		mesa.Status,
		mesa.IDUser,
		mesa.IDCli,
		mesa.Nick,

		mesa.Abertura,
		mesa.QtdPessoas,
		mesa.Turista,
		mesa.Celular,
		mesa.Apelido,
		id,
	)
	return err
}

// BuscarMesas busca mesas no banco de dados com filtro opcional por STATUS
func (repositorio Mesas) BuscarMesas(status string) ([]modelos.Mesa, error) {
	// Construir a query base
	query := `SELECT ID, MESA_CARTAO, STATUS, ID_USER, ID_CLI, NICK,  ABERTURA, QTD_PESSOAS, TURISTA, CELULAR, APELIDO FROM MESA`
	var args []interface{}

	// Adicionar filtro de STATUS, se informado
	if status != "" {
		query += " WHERE STATUS = ?"
		args = append(args, status)
	}

	// Executar a consulta
	rows, err := repositorio.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Processar os resultados
	var mesas []modelos.Mesa
	for rows.Next() {
		var mesa modelos.Mesa
		if err := rows.Scan(
			&mesa.ID,
			&mesa.MesaCartao,
			&mesa.Status,
			&mesa.IDUser,
			&mesa.IDCli,
			&mesa.Nick,

			&mesa.Abertura,
			&mesa.QtdPessoas,
			&mesa.Turista,
			&mesa.Celular,
			&mesa.Apelido,
		); err != nil {
			return nil, err
		}
		mesas = append(mesas, mesa)
	}

	return mesas, nil
}

// BuscarMesaPorID busca uma única mesa pelo ID no banco de dados
func (repositorio Mesas) BuscarMesaPorID(mesaID uint64) (modelos.Mesa, error) {
	query := `
		SELECT ID, MESA_CARTAO, STATUS, ID_USER, ID_CLI, NICK,  ABERTURA, QTD_PESSOAS, TURISTA, CELULAR, APELIDO
		FROM MESA
		WHERE ID = ?
	`

	var mesa modelos.Mesa
	err := repositorio.db.QueryRow(query, mesaID).Scan(
		&mesa.ID,
		&mesa.MesaCartao,
		&mesa.Status,
		&mesa.IDUser,
		&mesa.IDCli,
		&mesa.Nick,

		&mesa.Abertura,
		&mesa.QtdPessoas,
		&mesa.Turista,
		&mesa.Celular,
		&mesa.Apelido,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return modelos.Mesa{}, errors.New("mesa não encontrada")
		}
		return modelos.Mesa{}, err
	}

	return mesa, nil
}
