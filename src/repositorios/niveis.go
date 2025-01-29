package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

// Niveis representa um repositório de níveis de acesso
type Niveis struct {
	db *sql.DB
}

// NovoRepositorioDeNiveis cria um novo repositório de níveis de acesso
func NovoRepositorioDeNiveis(db *sql.DB) *Niveis {
	return &Niveis{db}
}

// Criar insere um novo nível de acesso no banco de dados
func (repositorio Niveis) CriarNivel(nivel modelos.NivelAcesso) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO NIVEL_ACESSO (CODIGO, NOME_ACESSO, NIVEL) VALUES (?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(nivel.Codigo, nivel.NomeAcesso, nivel.Nivel)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// Buscar traz todos os níveis de acesso
func (repositorio Niveis) BuscarNiveis() ([]modelos.NivelAcesso, error) {
	rows, erro := repositorio.db.Query("SELECT ID, CODIGO, NOME_ACESSO, NIVEL FROM NIVEL_ACESSO")
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var niveis []modelos.NivelAcesso
	for rows.Next() {
		var nivel modelos.NivelAcesso

		// Scan para capturar os dados do banco de dados
		if erro = rows.Scan(&nivel.ID, &nivel.Codigo, &nivel.NomeAcesso, &nivel.Nivel); erro != nil {
			return nil, erro
		}

		niveis = append(niveis, nivel)
	}

	return niveis, nil
}

// BuscarPorID traz um nível de acesso específico pelo ID
func (repositorio Niveis) BuscarNivelPorID(nivelID uint64) (modelos.NivelAcesso, error) {
	row := repositorio.db.QueryRow("SELECT ID, CODIGO, NOME_ACESSO, NIVEL FROM NIVEL_ACESSO WHERE ID = ?", nivelID)

	var nivel modelos.NivelAcesso
	if erro := row.Scan(&nivel.ID, &nivel.Codigo, &nivel.NomeAcesso, &nivel.Nivel); erro != nil {
		return modelos.NivelAcesso{}, erro
	}

	return nivel, nil
}

// Atualizar altera um nível de acesso no banco de dados
func (repositorio Niveis) AtualizarNivel(nivelID uint64, nivel modelos.NivelAcesso) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE NIVEL_ACESSO SET  NOME_ACESSO = ?, NIVEL = ? WHERE ID = ?",
	)
	//	fmt.Println(nivel.NomeAcesso, nivel.Nivel, nivelID)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(nivel.NomeAcesso, nivel.Nivel, nivelID)
	return erro
}
