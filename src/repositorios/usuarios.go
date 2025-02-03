package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

// Usuarios representa um repositório de usuarios
type Usuarios struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios cria um repositório de usuários
func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

// Criar insere um usuário no banco de dados
func (repositorio Usuarios) CriarUsuario(usuario modelos.Usuario) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO USUARIOS (NOME, NICK, EMAIL, SENHA, CDEMP, NIVEL) VALUES(?, ?, ?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha, usuario.CDEMP, usuario.Nivel)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil

}

// Buscar traz todos os usuários que atendem um filtro de nome ou nick
func (repositorio Usuarios) BuscarUsuarios(nomeOuNick string, cdEmp string) ([]modelos.Usuario, error) {
	//fmt.Println("0 " + nomeOuNick)
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) // %nomeOuNick%
	//fmt.Println("1 " + nomeOuNick)
	linhas, erro := repositorio.db.Query(
		"SELECT ID, NOME, NICK, EMAIL, CRIADOEM, CDEMP, NIVEL  FROM USUARIOS WHERE (CDEMP = ?) AND (NOME LIKE ? OR NICK LIKE ?)",
		cdEmp, nomeOuNick, nomeOuNick,
	)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
			&usuario.CDEMP,
			&usuario.Nivel,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// BuscarPorID traz um usuário do banco de dados
func (repositorio Usuarios) BuscarPorID(ID uint64) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"SELECT ID, NOME, NICK, EMAIL,NIVEL,  CRIADOEM FROM USUARIOS WHERE ID = ?",
		ID,
	)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario modelos.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.Nivel,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// Atualizar altera as informações de um usuário no banco de dados
func (repositorio Usuarios) AtualizarUsuario(ID uint64, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE USUARIOS SET AGUARDANDO_SYNC = ?, NOME = ?, NICK = ?, EMAIL = ?, NIVEL = ? WHERE ID = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuario.AGUARDANDO_SYNC, usuario.Nome, usuario.Nick, usuario.Email, usuario.Nivel, ID); erro != nil {
		return erro
	}

	return nil
}

// Deletar exclui as informações de um usuário no banco de dados
func (repositorio Usuarios) Deletar(ID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM USUARIOS WHERE ID = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

// BuscarPorEmail busca um usuário por email e retorna o seu id e senha com hash
func (repositorio Usuarios) BuscarPorEmail(email string) (modelos.Usuario, error) {
	linha, erro := repositorio.db.Query("SELECT ID, SENHA,  NIVEL,CDEMP FROM USUARIOS WHERE EMAIL = ?", email)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()

	var usuario modelos.Usuario
	//cdemp pega no json
	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha, &usuario.Nivel, &usuario.CDEMP); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil

}

// BuscarSenha traz a senha de um usuário pelo ID
func (repositorio Usuarios) BuscarSenha(usuarioID uint64) (string, error) {
	linha, erro := repositorio.db.Query("SELECT SENHA FROM USUARIOS WHERE ID = ?", usuarioID)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.Senha); erro != nil {
			return "", erro
		}
	}

	return usuario.Senha, nil
}

// AtualizarSenha altera a senha de um usuário
func (repositorio Usuarios) AtualizarSenha(usuarioID uint64, senha string) error {
	statement, erro := repositorio.db.Prepare("UPDATE USUARIOS SET SENHA = ? WHERE ID = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(senha, usuarioID); erro != nil {
		return erro
	}

	return nil
}

// PesquisarPorNomeOuNick busca usuários com base em parte do nome ou apelido
func (repositorio Usuarios) PesquisarPorNomeOuNick(nomeOuNick string, cdEmp string) ([]modelos.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick)

	linhas, erro := repositorio.db.Query(
		"SELECT ID, NOME, NICK, EMAIL, CRIADOEM, CDEMP, NIVEL FROM USUARIOS WHERE ( NOME LIKE ? OR NICK LIKE ?) and CDEMP = ? ",
		nomeOuNick, nomeOuNick, cdEmp,
	)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario
	for linhas.Next() {
		var usuario modelos.Usuario
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
			&usuario.CDEMP,
			&usuario.Nivel,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

//funções de SYNC

// Função para buscar usuários não sincronizados (ID_HOSTWEB = 0)
func (repositorio Usuarios) BuscarUsuariosNaoSincronizados() ([]modelos.Usuario, error) {
	rows, erro := repositorio.db.Query(`
		SELECT ID, ID_HOSTWEB, AGUARDANDO_SYNC, NOME, NICK, EMAIL, SENHA, CRIADOEM, CDEMP, NIVEL
		FROM USUARIOS WHERE ID_HOSTWEB = 0
	`)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var usuarios []modelos.Usuario
	for rows.Next() {
		var usuario modelos.Usuario
		erro := rows.Scan(
			&usuario.ID,
			&usuario.ID_HOSTWEB,
			&usuario.AGUARDANDO_SYNC,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.Senha,
			&usuario.CriadoEm,
			&usuario.CDEMP,
			&usuario.Nivel,
		)
		if erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// Atualiza o ID_HOSTWEB após sincronização bem-sucedida com a nuvem
func (repositorio Usuarios) AtualizarUsuarioHostWeb(idLocal uint64, idNuvem uint64) error {
	_, erro := repositorio.db.Exec(`
		UPDATE USUARIOS SET ID_HOSTWEB = ? WHERE ID = ?
	`, idNuvem, idLocal)
	return erro
}

// Função para verificar se o usuário já existe localmente
func (repositorio Usuarios) UsuarioExisteNoLocal(idHostWeb uint64) bool {
	var existe int
	erro := repositorio.db.QueryRow("SELECT COUNT(*) FROM USUARIOS WHERE ID_HOSTWEB = ?", idHostWeb).Scan(&existe)
	return erro == nil && existe > 0
}

// Função para verificar se o usuário já existe na nuvem
func (repositorio Usuarios) UsuarioExisteNaNuvem(idHostLocal uint64) bool {
	var existe int
	erro := repositorio.db.QueryRow("SELECT COUNT(*) FROM USUARIOS WHERE ID_HOSTLOCAL = ?", idHostLocal).Scan(&existe)
	return erro == nil && existe > 0
}

// Função para atualizar um usuário existente pelo ID_HOSTLOCAL
func (repositorio Usuarios) AtualizarUsuarioPorHostLocal(idHostLocal uint64, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE USUARIOS SET ID_HOSTWEB = ?, AGUARDANDO_SYNC = ?, NOME = ?, NICK = ?, EMAIL = ?, SENHA = ?, CRIADOEM = ?, CDEMP = ?, NIVEL = ?
		WHERE ID_HOSTLOCAL = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		usuario.ID_HOSTWEB,
		usuario.AGUARDANDO_SYNC,
		usuario.Nome,
		usuario.Nick,
		usuario.Email,
		usuario.Senha,
		usuario.CriadoEm,
		usuario.CDEMP,
		usuario.Nivel,
		idHostLocal,
	)
	return erro
}

// Função para atualizar um usuário existente pelo ID_HOSTWEB
func (repositorio Usuarios) AtualizarUsuarioPorHostWeb(idHostWeb uint64, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE USUARIOS SET AGUARDANDO_SYNC = ?, NOME = ?, NICK = ?, EMAIL = ?, SENHA = ?, CRIADOEM = ?, CDEMP = ?, NIVEL = ?
		WHERE ID_HOSTWEB = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		usuario.AGUARDANDO_SYNC,
		usuario.Nome,
		usuario.Nick,
		usuario.Email,
		usuario.Senha,
		usuario.CriadoEm,
		usuario.CDEMP,
		usuario.Nivel,
		idHostWeb,
	)
	return erro
}

// BuscarUsuariosAguardandoSync busca usuários com AGUARDANDO_SYNC = "S"
func (repositorio Usuarios) BuscarUsuariosAguardandoSync() ([]modelos.Usuario, error) {
	rows, erro := repositorio.db.Query(`
		SELECT ID, ID_HOSTWEB, NOME, NICK, EMAIL, SENHA, CRIADOEM, CDEMP, NIVEL
		FROM USUARIOS WHERE AGUARDANDO_SYNC = 'S'
	`)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var usuarios []modelos.Usuario
	for rows.Next() {
		var usuario modelos.Usuario
		erro := rows.Scan(
			&usuario.ID,
			&usuario.ID_HOSTWEB,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.Senha,
			&usuario.CriadoEm,
			&usuario.CDEMP,
			&usuario.Nivel,
		)
		if erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// DesmarcarAguardandoSyncUsuario desmarca o campo AGUARDANDO_SYNC no usuário local
func (repositorio Usuarios) DesmarcarAguardandoSyncUsuario(usuarioID uint64) error {
	_, erro := repositorio.db.Exec("UPDATE USUARIOS SET AGUARDANDO_SYNC = '' WHERE ID = ?", usuarioID)
	return erro
}

// CriarUsuarioSoLocal insere uma nova entrada de usuário no banco de dados local
func (repositorio Usuarios) CriarUsuarioSoLocal(usuario modelos.Usuario) (uint64, error) {
	// Prepara e executa o SQL de inserção
	statement, erro := repositorio.db.Prepare(`
		INSERT INTO USUARIOS (AGUARDANDO_SYNC, NOME, NICK, EMAIL, SENHA, CRIADOEM, CDEMP, NIVEL)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(
		usuario.AGUARDANDO_SYNC,
		usuario.Nome,
		usuario.Nick,
		usuario.Email,
		usuario.Senha,
		usuario.CriadoEm,
		usuario.CDEMP,
		usuario.Nivel,
	)
	if erro != nil {
		return 0, erro
	}

	idUsuario, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	// Define o ID do usuário com o valor gerado
	usuario.ID = uint64(idUsuario)

	return uint64(idUsuario), nil
}

// BuscarNick busca usuários pelo nick e código da empresa
func (repositorio Usuarios) BuscarNick(nick string, cdEmp string) ([]modelos.Usuario, error) {
	// Executa a consulta com filtro por nick e CDEMP
	linhas, erro := repositorio.db.Query(`
		SELECT ID, NOME, NICK, EMAIL, CRIADOEM, CDEMP, NIVEL
		FROM USUARIOS
		WHERE CDEMP = ? AND TRIM(NICK) = ?
	`, cdEmp, nick)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario

	// Itera sobre os resultados da consulta
	for linhas.Next() {
		var usuario modelos.Usuario

		if erro := linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
			&usuario.CDEMP,
			&usuario.Nivel,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}
