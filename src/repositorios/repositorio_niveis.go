package repositorios

import "database/sql"

// BuscarNivelPorCodigo busca o nível de acesso pelo código do programa
func (repositorio Niveis) BuscarNivelPorCodigo(codigo string) (int, error) {
	var nivel int
	row := repositorio.db.QueryRow("SELECT nivel FROM nivel_acesso WHERE codigo = ?", codigo)

	if erro := row.Scan(&nivel); erro != nil {
		return 0, erro
	}

	return nivel, nil
}

// BuscarOuCriarNivelPorCodigo busca o nível de acesso pelo código do programa.
// Se não encontrar, cria o nível automaticamente com valor 0.
func (repositorio Niveis) BuscarOuCriarNivelPorCodigo(codigo string) (int, error) {
	var nivel int
	row := repositorio.db.QueryRow("SELECT nivel FROM nivel_acesso WHERE codigo = ?", codigo)

	// Tenta buscar o nível
	if erro := row.Scan(&nivel); erro != nil {
		// Se o nível não for encontrado, cria o nível automaticamente com 0
		if erro == sql.ErrNoRows {
			statement, erro := repositorio.db.Prepare("INSERT INTO nivel_acesso (codigo, nome_acesso, nivel, criadoEm) VALUES (?, ?, ?, NOW())")
			if erro != nil {
				return 0, erro
			}
			defer statement.Close()

			// Define o nome_acesso padrão como "Acesso Automático"
			_, erro = statement.Exec(codigo, codigo, 0)
			if erro != nil {
				return 0, erro
			}

			// Retorna o nível 0 recém-adicionado
			return 0, nil
		}
		// Se o erro for diferente de "não encontrado"
		return 0, erro
	}

	// Retorna o nível encontrado
	return nivel, nil
}
