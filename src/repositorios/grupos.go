package repositorios

import (
	"api/src/modelos"
	"api/src/utils"
	"database/sql"
	"fmt"
	"log"
)

// Grupos representa um repositório de grupos
type Grupos struct {
	db *sql.DB
}

// NovoRepositorioDeGrupos cria um repositório de grupos
func NovoRepositorioDeGrupos(db *sql.DB) *Grupos {
	return &Grupos{db}
}

// CriarGrupo insere um grupo no banco de dados local e sincroniza com a nuvem
func (repositorio Grupos) CriarGrupo(grupo modelos.Grupo) (uint64, error) {
	statement, erro := repositorio.db.Prepare(`
		INSERT INTO GRUPOS (COD_GP, NOME, TIPO, CONTA_CONTABIL, AGUARDANDO_SYNC, ID_HOSTWEB)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	// Insere os valores para cada coluna do grupo
	resultado, erro := statement.Exec(
		grupo.COD_GP, grupo.NOME, grupo.TIPO, grupo.CONTA_CONTABIL, grupo.AGUARDANDO_SYNC, grupo.ID_HOSTLOCAL,
	)
	if erro != nil {
		return 0, erro
	}

	// Obtém o último ID inserido
	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	grupo.ID = uint64(ultimoIDInserido)

	// Sincroniza com a nuvem e captura o ID gerado na nuvem
	idHostWeb, erro := utils.SincronizarGrupoNuvem(grupo)
	if erro != nil {
		log.Printf("Erro ao sincronizar com a nuvem: %v", erro)
		return uint64(ultimoIDInserido), nil
	}

	// Atualiza o campo ID_HOSTWEB com o ID retornado da nuvem
	_, erro = repositorio.db.Exec("UPDATE GRUPOS SET ID_HOSTWEB = ? WHERE ID = ?", idHostWeb, ultimoIDInserido)
	if erro != nil {
		return uint64(ultimoIDInserido), fmt.Errorf("Erro ao atualizar ID_HOSTWEB: %v", erro)
	}

	return uint64(ultimoIDInserido), nil
}

// Buscar traz todos os grupos
func (repositorio Grupos) BuscarGrupo() ([]modelos.Grupo, error) {
	rows, erro := repositorio.db.Query(`
		SELECT ID, COD_GP, NOME, TIPO, CONTA_CONTABIL, AGUARDANDO_SYNC
		FROM GRUPOS
	`)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var grupos []modelos.Grupo
	for rows.Next() {
		var grupo modelos.Grupo
		if erro = rows.Scan(
			&grupo.ID,
			&grupo.COD_GP,
			&grupo.NOME,
			&grupo.TIPO,
			&grupo.CONTA_CONTABIL,
			&grupo.AGUARDANDO_SYNC,
		); erro != nil {
			log.Println(erro)
			return nil, erro
		}
		grupos = append(grupos, grupo)
	}

	return grupos, nil
}

// BuscarPorID busca um grupo no banco de dados pelo ID
func (repositorio Grupos) BuscarGrupoPorID(grupoID uint64) (modelos.Grupo, error) {
	row := repositorio.db.QueryRow(`
		SELECT ID, COD_GP, NOME, TIPO, CONTA_CONTABIL, AGUARDANDO_SYNC, ID_HOSTWEB, ID_HOSTLOCAL
		FROM GRUPOS WHERE ID = ?
	`, grupoID)

	var grupo modelos.Grupo
	if erro := row.Scan(
		&grupo.ID,
		&grupo.COD_GP,
		&grupo.NOME,
		&grupo.TIPO,
		&grupo.CONTA_CONTABIL,
		&grupo.AGUARDANDO_SYNC,
		&grupo.ID_HOSTWEB,
		&grupo.ID_HOSTLOCAL,
	); erro != nil {
		return modelos.Grupo{}, erro
	}

	return grupo, nil
}

// AtualizarGrupo altera um grupo no banco de dados
func (repositorio Grupos) AtualizarGrupo(grupoID uint64, grupo modelos.Grupo) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE GRUPOS SET
			COD_GP = ?, NOME = ?, TIPO = ?, CONTA_CONTABIL = ?, AGUARDANDO_SYNC = ?
		WHERE ID = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		grupo.COD_GP,
		grupo.NOME,
		grupo.TIPO,
		grupo.CONTA_CONTABIL,
		grupo.AGUARDANDO_SYNC,
		grupoID,
	)
	return erro
}

// Deletar exclui um grupo do banco de dados
func (repositorio Grupos) DeletarGrupo(grupoID uint64) error {
	statement, erro := repositorio.db.Prepare("UPDATE GRUPOS SET TIPO = ?, AGUARDANDO_SYNC = ? WHERE ID = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec("D", "S", grupoID)
	return erro
}

//PARTE NOVA SYNC

// BuscarGruposNaoSincronizados busca grupos com ID_HOSTWEB = 0
func (repositorio Grupos) BuscarGruposNaoSincronizados() ([]modelos.Grupo, error) {
	rows, erro := repositorio.db.Query(`
		SELECT ID, ID_HOSTWEB,  AGUARDANDO_SYNC, COD_GP, NOME, TIPO, CONTA_CONTABIL
		FROM GRUPOS WHERE ID_HOSTWEB = 0
	`)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var grupos []modelos.Grupo
	for rows.Next() {
		var grupo modelos.Grupo
		erro := rows.Scan(
			&grupo.ID,
			&grupo.ID_HOSTWEB,
			&grupo.AGUARDANDO_SYNC,
			&grupo.COD_GP,
			&grupo.NOME,
			&grupo.TIPO,
			&grupo.CONTA_CONTABIL,
		)
		if erro != nil {
			return nil, erro
		}
		grupos = append(grupos, grupo)
	}

	return grupos, nil
}

// GrupoExisteNaNuvem verifica se o grupo já existe na nuvem (com base no ID_HOSTLOCAL)
func (repositorio Grupos) GrupoExisteNaNuvem(idHostLocal uint64) bool {
	var existe int
	erro := repositorio.db.QueryRow("SELECT COUNT(*) FROM GRUPOS WHERE ID_HOSTLOCAL = ?", idHostLocal).Scan(&existe)
	return erro == nil && existe > 0
}

// AtualizarGrupoPorHostWeb atualiza um grupo existente com base no ID_HOSTWEB
func (repositorio Grupos) AtualizarGrupoPorHostWeb(idHostWeb uint64, grupo modelos.Grupo) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE GRUPOS SET 
			ID_HOSTLOCAL = ?, AGUARDANDO_SYNC = ?, COD_GP = ?, NOME = ?, TIPO = ?, CONTA_CONTABIL = ?
		WHERE ID_HOSTWEB = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		grupo.ID_HOSTLOCAL,
		grupo.AGUARDANDO_SYNC,
		grupo.COD_GP,
		grupo.NOME,
		grupo.TIPO,
		grupo.CONTA_CONTABIL,
		idHostWeb,
	)
	return erro
}

// AtualizarGrupoHostWeb atualiza o campo ID_HOSTWEB após a sincronização
func (repositorio Grupos) AtualizarGrupoHostWeb(idLocal uint64, idNuvem uint64) error {
	_, erro := repositorio.db.Exec(`
		UPDATE GRUPOS SET ID_HOSTWEB = ? WHERE ID = ?
	`, idNuvem, idLocal)
	return erro
}

// GrupoExisteNoLocal verifica se o grupo já existe localmente (com base no ID_HOSTWEB)
func (repositorio Grupos) GrupoExisteNoLocal(idHostWeb uint64) bool {
	var existe int
	erro := repositorio.db.QueryRow("SELECT COUNT(*) FROM GRUPOS WHERE ID_HOSTWEB = ?", idHostWeb).Scan(&existe)
	return erro == nil && existe > 0
}

// AtualizarGrupoPorHostLocal atualiza um grupo existente com base no ID_HOSTLOCAL
func (repositorio Grupos) AtualizarGrupoPorHostLocal(idHostLocal uint64, grupo modelos.Grupo) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE GRUPOS SET 
			ID_HOSTWEB = ?, AGUARDANDO_SYNC = ?, COD_GP = ?, NOME = ?, TIPO = ?, CONTA_CONTABIL = ?
		WHERE ID_HOSTLOCAL = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		grupo.ID_HOSTWEB,
		grupo.AGUARDANDO_SYNC,
		grupo.COD_GP,
		grupo.NOME,
		grupo.TIPO,
		grupo.CONTA_CONTABIL,
		idHostLocal,
	)
	return erro
}

// CriarGrupoSoLocal insere um novo grupo no banco de dados local
func (repositorio Grupos) CriarGrupoSoLocal(grupo modelos.Grupo) (uint64, error) {
	statement, erro := repositorio.db.Prepare(`
		INSERT INTO GRUPOS (COD_GP, NOME, TIPO, CONTA_CONTABIL) 
		VALUES (?, ?, ?, ?)
	`)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	// Insere os valores do grupo
	resultado, erro := statement.Exec(
		grupo.COD_GP,
		grupo.NOME,
		grupo.TIPO,
		grupo.CONTA_CONTABIL,
	)
	if erro != nil {
		return 0, erro
	}

	// Obtém o último ID inserido
	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// BuscarGruposAguardandoSync busca grupos com AGUARDANDO_SYNC = "S"
func (repositorio Grupos) BuscarGruposAguardandoSync() ([]modelos.Grupo, error) {
	rows, erro := repositorio.db.Query(`
		SELECT ID, ID_HOSTWEB, AGUARDANDO_SYNC, COD_GP, NOME, TIPO, CONTA_CONTABIL
		FROM GRUPOS WHERE AGUARDANDO_SYNC = 'S'
	`)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var grupos []modelos.Grupo
	for rows.Next() {
		var grupo modelos.Grupo
		erro := rows.Scan(
			&grupo.ID,
			&grupo.ID_HOSTWEB,
			&grupo.AGUARDANDO_SYNC,
			&grupo.COD_GP,
			&grupo.NOME,
			&grupo.TIPO,
			&grupo.CONTA_CONTABIL,
		)
		if erro != nil {
			return nil, erro
		}
		grupos = append(grupos, grupo)
	}

	return grupos, nil
}

// DesmarcarAguardandoSyncLocal desmarca o campo AGUARDANDO_SYNC no grupo local
func (repositorio Grupos) DesmarcarAguardandoSyncLocal(grupoID uint64) error {
	_, erro := repositorio.db.Exec("UPDATE GRUPOS SET AGUARDANDO_SYNC = '' WHERE ID = ?", grupoID)
	return erro
}

// BuscarGrupos recupera todos os registros de grupos com filtros de status e filtro LIKE
func (repositorio Grupos) BuscarGrupos(status string, filtro string) ([]modelos.Grupo, error) {
	// Construir a query base
	query := `SELECT ID, ID_HOSTWEB, COD_GP, NOME, TIPO, CONTA_CONTABIL 
              FROM GRUPOS WHERE 1=1` // WHERE 1=1 permite adicionar condições dinamicamente

	// Inicializar os argumentos da query
	var args []interface{}

	// Adicionar condição de status, se informado
	if status != "" {
		query += " AND TIPO = ?"
		args = append(args, status)
	}

	// Adicionar condição de filtro para os campos NOME e COD_GP, se informado
	if filtro != "" {
		filtro = fmt.Sprintf("%%%s%%", filtro) // Adicionar % para o LIKE
		query += " AND (NOME LIKE ? OR COD_GP LIKE ?)"
		args = append(args, filtro, filtro)
	}

	// Adicionar ordem na consulta
	query += " ORDER BY COD_GP"

	// Executar a consulta com filtros aplicados
	rows, erro := repositorio.db.Query(query, args...)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	// Mapear os resultados para o slice de grupos
	var grupos []modelos.Grupo
	for rows.Next() {
		var grupo modelos.Grupo
		if erro := rows.Scan(
			&grupo.ID,
			&grupo.ID_HOSTLOCAL,
			&grupo.COD_GP,
			&grupo.NOME,
			&grupo.TIPO,
			&grupo.CONTA_CONTABIL,
		); erro != nil {
			return nil, erro
		}
		grupos = append(grupos, grupo)
	}

	return grupos, nil
}
