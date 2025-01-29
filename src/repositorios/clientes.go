package repositorios

import (
	"api/src/modelos"
	"api/src/utils"
	"database/sql"
	"fmt"
	"log"
)

// Clientes representa um repositório de clientes
type Clientes struct {
	db *sql.DB
}

// NovoRepositorioDeClientes cria um repositório de clientes
func NovoRepositorioDeClientes(db *sql.DB) *Clientes {
	return &Clientes{db}
}

// CriarCliente insere uma nova entrada de cliente no banco de dados local e sincroniza com a nuvem
func (repositorio Clientes) CriarCliente(cliente modelos.Cliente) (uint64, error) {
	// Verifica se o status está vazio e define como "A"
	if cliente.STATUS == "" {
		cliente.STATUS = "A"
	}
	// Prepara e executa o SQL de inserção
	statement, erro := repositorio.db.Prepare(`
		INSERT INTO CLIENTE (AGUARDANDO_SYNC, CNPJ_CPF, NOME, FANTASIA, 
		ENDERE, CIDADE, BAIRRO, CEP, UF, TELE1, CEL1, CONTATO, EMAIL, PLACA, NUMERO, CIDIBGE, COMPLE, STATUS, ID_HOSTWEB)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(
		cliente.AGUARDANDO_SYNC,
		cliente.CNPJ_CPF,
		cliente.NOME,
		cliente.FANTASIA,
		cliente.ENDERE,
		cliente.CIDADE,
		cliente.BAIRRO,
		cliente.CEP,
		cliente.UF,
		cliente.TELE1,
		cliente.CEL1,
		cliente.CONTATO,
		cliente.EMAIL,
		cliente.PLACA,
		cliente.NUMERO,
		cliente.CIDIBGE,
		cliente.COMPLE,
		cliente.STATUS,
		cliente.ID,
	)
	if erro != nil {
		return 0, erro
	}

	idCliente, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	// Define o ID do cliente com o valor gerado
	cliente.ID = uint64(idCliente)

	// Sincroniza com a nuvem e captura o ID gerado na nuvem
	idHostWeb, erro := utils.SincronizarClienteNuvem(cliente)
	if erro != nil {
		log.Printf("Erro ao sincronizar com a nuvem: %v", erro)
		return uint64(idCliente), nil
	}

	// Atualiza o campo ID_HOSTWEB com o ID retornado da nuvem
	_, erro = repositorio.db.Exec("UPDATE CLIENTE SET ID_HOSTWEB = ? WHERE ID = ?", idHostWeb, idCliente)
	if erro != nil {
		return uint64(idCliente), fmt.Errorf("Erro ao atualizar ID_HOSTWEB: %v", erro)
	}

	return uint64(idCliente), nil
}

// CriarClienteSoLocal insere uma nova entrada de cliente no banco de dados local
func (repositorio Clientes) CriarClienteSoLocal(cliente modelos.Cliente) (uint64, error) {
	// Verifica se o status está vazio e define como "A"
	if cliente.STATUS == "" {
		cliente.STATUS = "A"
	}
	// Prepara e executa o SQL de inserção
	statement, erro := repositorio.db.Prepare(`
		INSERT INTO CLIENTE (AGUARDANDO_SYNC, CNPJ_CPF, NOME, FANTASIA, ENDERE, CIDADE, BAIRRO, CEP, UF, TELE1, CEL1, CONTATO, EMAIL, PLACA, NUMERO, CIDIBGE, COMPLE, STATUS)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?)
	`)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(
		cliente.AGUARDANDO_SYNC,
		cliente.CNPJ_CPF,
		cliente.NOME,
		cliente.FANTASIA,
		cliente.ENDERE,
		cliente.CIDADE,
		cliente.BAIRRO,
		cliente.CEP,
		cliente.UF,
		cliente.TELE1,
		cliente.CEL1,
		cliente.CONTATO,
		cliente.EMAIL,
		cliente.PLACA,
		cliente.NUMERO,
		cliente.CIDIBGE,
		cliente.COMPLE,
		cliente.STATUS,
	)
	if erro != nil {
		return 0, erro
	}

	idCliente, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	// Define o ID do cliente com o valor gerado
	cliente.ID = uint64(idCliente)

	return uint64(idCliente), nil
}

// BuscarClientes traz todos os clientes que atendem a um filtro de nome, status, CPF e celular
func (repositorio Clientes) BuscarClientes(nome string, status string) ([]modelos.Cliente, error) {
	// Formata o filtro para uso com LIKE
	nome = fmt.Sprintf("%%%s%%", nome)

	// Cria a query SQL base
	query := `
		SELECT ID, CNPJ_CPF, NOME, FANTASIA, ENDERE, CIDADE, BAIRRO, CEP, UF,
		TELE1, CEL1, CONTATO, EMAIL, PLACA, NUMERO, CIDIBGE, COMPLE, STATUS
		FROM CLIENTE
		WHERE (NOME LIKE ? OR CNPJ_CPF LIKE ? OR CEL1 LIKE ?)
	`

	// Inicializa os argumentos da consulta
	var args []interface{}
	args = append(args, nome, nome, nome)

	// Adiciona o filtro de status, se necessário
	if status != "" {
		query += " AND STATUS = ?"
		args = append(args, status)
	}

	// Executa a consulta SQL
	rows, erro := repositorio.db.Query(query, args...)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	// Preenche o slice de clientes com os resultados
	var clientes []modelos.Cliente
	for rows.Next() {
		var cliente modelos.Cliente
		if erro = rows.Scan(
			&cliente.ID,
			&cliente.CNPJ_CPF,
			&cliente.NOME,
			&cliente.FANTASIA,
			&cliente.ENDERE,
			&cliente.CIDADE,
			&cliente.BAIRRO,
			&cliente.CEP,
			&cliente.UF,
			&cliente.TELE1,
			&cliente.CEL1,
			&cliente.CONTATO,
			&cliente.EMAIL,
			&cliente.PLACA,
			&cliente.NUMERO,
			&cliente.CIDIBGE,
			&cliente.COMPLE,
			&cliente.STATUS,
		); erro != nil {
			return nil, erro
		}
		clientes = append(clientes, cliente)
	}

	return clientes, nil
}

// BuscarPorID busca um cliente no banco de dados pelo ID
func (repositorio Clientes) BuscarPorID(clienteID uint64) (modelos.Cliente, error) {
	row := repositorio.db.QueryRow(`SELECT ID, CNPJ_CPF, NOME, FANTASIA, ENDERE, CIDADE, BAIRRO, CEP, UF, TELE1, CEL1, CONTATO, EMAIL,  PLACA, NUMERO, CIDIBGE, COMPLE,STATUS FROM CLIENTE WHERE ID = ?`, clienteID)

	var cliente modelos.Cliente
	if erro := row.Scan(
		&cliente.ID,
		&cliente.CNPJ_CPF,
		&cliente.NOME,
		&cliente.FANTASIA,
		&cliente.ENDERE,
		&cliente.CIDADE,
		&cliente.BAIRRO,
		&cliente.CEP,
		&cliente.UF,
		&cliente.TELE1,
		&cliente.CEL1,
		&cliente.CONTATO,
		&cliente.EMAIL,
		&cliente.PLACA,
		&cliente.NUMERO,
		&cliente.CIDIBGE,
		&cliente.COMPLE,
		&cliente.STATUS,
	); erro != nil {
		return modelos.Cliente{}, erro
	}

	return cliente, nil
}

// Atualizar altera um cliente no banco de dados
func (repositorio Clientes) AtualizarCliente(clienteID uint64, cliente modelos.Cliente) error {
	// Verifica se o status está vazio e define como "A"
	if cliente.STATUS == "" {
		cliente.STATUS = "A"
	}
	//fmt.Println(cliente.AGUARDANDO_SYNC, "entrei")
	statement, erro := repositorio.db.Prepare(
		"UPDATE CLIENTE SET AGUARDANDO_SYNC = ?, CNPJ_CPF = ?, NOME = ?, FANTASIA = ?, ENDERE = ?, CIDADE = ?, BAIRRO = ?, CEP = ?, UF = ?, TELE1 = ?, CEL1 = ?, CONTATO = ?, EMAIL = ?, PLACA = ?, NUMERO = ?, CIDIBGE = ?, COMPLE = ?,  STATUS=? WHERE ID = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(cliente.AGUARDANDO_SYNC, cliente.CNPJ_CPF,
		cliente.NOME, cliente.FANTASIA, cliente.ENDERE, cliente.CIDADE,
		cliente.BAIRRO, cliente.CEP, cliente.UF, cliente.TELE1,
		cliente.CEL1, cliente.CONTATO, cliente.EMAIL, cliente.PLACA,
		cliente.NUMERO, cliente.CIDIBGE, cliente.COMPLE,
		cliente.STATUS, clienteID)
	return erro
}

// Deletar exclui um cliente no banco de dados
func (repositorio Clientes) DeletarCliente(clienteID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"update CLIENTE set AGUARDANDO_SYNC = ?, STATUS = ? where ID = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec("S", "D", clienteID)
	return erro
}

/// ROTINAS DE SYSNC

// Função para buscar clientes não sincronizados (ID_HOSTWEB = 0)
func (repositorio Clientes) BuscarClientesNaoSincronizados() ([]modelos.Cliente, error) {
	rows, erro := repositorio.db.Query(`
		SELECT ID, ID_HOSTWEB, AGUARDANDO_SYNC, CNPJ_CPF, NOME, FANTASIA, ENDERE, CIDADE, BAIRRO, CEP, UF, TELE1, CEL1, CONTATO, EMAIL, PLACA, NUMERO, CIDIBGE, COMPLE,STATUS
		FROM CLIENTE WHERE ID_HOSTWEB = 0
	`)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var clientes []modelos.Cliente
	for rows.Next() {
		var cliente modelos.Cliente
		erro := rows.Scan(
			&cliente.ID,
			&cliente.ID_HOSTWEB,
			&cliente.AGUARDANDO_SYNC,
			&cliente.CNPJ_CPF,
			&cliente.NOME,
			&cliente.FANTASIA,
			&cliente.ENDERE,
			&cliente.CIDADE,
			&cliente.BAIRRO,
			&cliente.CEP,
			&cliente.UF,
			&cliente.TELE1,
			&cliente.CEL1,
			&cliente.CONTATO,
			&cliente.EMAIL,
			&cliente.PLACA,
			&cliente.NUMERO,
			&cliente.CIDIBGE,
			&cliente.COMPLE,
			&cliente.STATUS,
		)
		if erro != nil {
			return nil, erro
		}
		clientes = append(clientes, cliente)
	}

	return clientes, nil
}

// Atualiza o ID_HOSTWEB após sincronização bem-sucedida com a nuvem
func (repositorio Clientes) AtualizarClienteHostWeb(idLocal uint64, idNuvem uint64) error {
	_, erro := repositorio.db.Exec(`
		UPDATE CLIENTE SET ID_HOSTWEB = ? WHERE ID = ?
	`, idNuvem, idLocal)
	return erro
}

// Função para verificar se o cliente já existe localmente
func (repositorio Clientes) ClienteExisteNoLocal(idHostWeb uint64) bool {
	var existe int
	erro := repositorio.db.QueryRow("SELECT COUNT(*) FROM CLIENTE WHERE ID_HOSTWEB = ?", idHostWeb).Scan(&existe)
	return erro == nil && existe > 0
}

// Função para verificar se o cliente já existe na nuvem
func (repositorio Clientes) ClienteExisteNaNuvem(idHostLocal uint64) bool {
	var existe int
	erro := repositorio.db.QueryRow("SELECT COUNT(*) FROM CLIENTE WHERE ID_HOSTLOCAL = ?", idHostLocal).Scan(&existe)
	return erro == nil && existe > 0
}

// Função para atualizar um cliente existente pelo ID_HOSTLOCAL
func (repositorio Clientes) AtualizarClientePorHostLocal(idHostLocal uint64, cliente modelos.Cliente) error {
	// Verifica se o status está vazio e define como "A"
	if cliente.STATUS == "" {
		cliente.STATUS = "A"
	}
	statement, erro := repositorio.db.Prepare(`
		UPDATE CLIENTE SET ID_HOSTWEB = ?, AGUARDANDO_SYNC = ?, CNPJ_CPF = ?, NOME = ?, FANTASIA = ?, ENDERE = ?, CIDADE = ?, BAIRRO = ?, CEP = ?, UF = ?, TELE1 = ?, CEL1 = ?, CONTATO = ?, EMAIL = ?, PLACA = ?, NUMERO = ?, CIDIBGE = ?, COMPLE = ?, STATUS = ?, 
		WHERE ID_HOSTLOCAL = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		cliente.ID_HOSTWEB,
		cliente.AGUARDANDO_SYNC,
		cliente.CNPJ_CPF,
		cliente.NOME,
		cliente.FANTASIA,
		cliente.ENDERE,
		cliente.CIDADE,
		cliente.BAIRRO,
		cliente.CEP,
		cliente.UF,
		cliente.TELE1,
		cliente.CEL1,
		cliente.CONTATO,
		cliente.EMAIL,
		cliente.PLACA,
		cliente.NUMERO,
		cliente.CIDIBGE,
		cliente.COMPLE,
		cliente.STATUS,

		idHostLocal,
	)
	return erro
}

// Função para atualizar um cliente existente pelo ID_HOSTWEB
func (repositorio Clientes) AtualizarClientePorHostWeb(idHostWeb uint64, cliente modelos.Cliente) error {
	// Verifica se o status está vazio e define como "A"
	if cliente.STATUS == "" {
		cliente.STATUS = "A"
	}
	statement, erro := repositorio.db.Prepare(`
		UPDATE CLIENTE SET AGUARDANDO_SYNC = ?, CNPJ_CPF = ?, NOME = ?, FANTASIA = ?, ENDERE = ?, CIDADE = ?, BAIRRO = ?, CEP = ?, UF = ?, TELE1 = ?, CEL1 = ?, CONTATO = ?, EMAIL = ?, PLACA = ?, NUMERO = ?, CIDIBGE = ?, COMPLE = ?, STATUS=?,
		WHERE ID_HOSTWEB = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		cliente.AGUARDANDO_SYNC,
		cliente.CNPJ_CPF,
		cliente.NOME,
		cliente.FANTASIA,
		cliente.ENDERE,
		cliente.CIDADE,
		cliente.BAIRRO,
		cliente.CEP,
		cliente.UF,
		cliente.TELE1,
		cliente.CEL1,
		cliente.CONTATO,
		cliente.EMAIL,
		cliente.PLACA,
		cliente.NUMERO,
		cliente.CIDIBGE,
		cliente.COMPLE,
		cliente.STATUS,

		idHostWeb,
	)
	return erro
}

// BuscarClientesAguardandoSync busca clientes com AGUARDANDO_SYNC = "S"
func (repositorio Clientes) BuscarClientesAguardandoSync() ([]modelos.Cliente, error) {
	rows, _ := repositorio.db.Query("SELECT ID,ID_HOSTWEB, CNPJ_CPF, NOME, FANTASIA, ENDERE, CIDADE, BAIRRO, CEP, UF, TELE1, CEL1, CONTATO, EMAIL, PLACA, NUMERO, CIDIBGE, COMPLE, STATUS FROM CLIENTE WHERE AGUARDANDO_SYNC = 'S'")
	defer rows.Close()

	var clientes []modelos.Cliente
	for rows.Next() {
		var cliente modelos.Cliente
		rows.Scan(&cliente.ID, &cliente.ID_HOSTWEB, &cliente.CNPJ_CPF,
			&cliente.NOME, &cliente.FANTASIA, &cliente.ENDERE,
			&cliente.CIDADE, &cliente.BAIRRO, &cliente.CEP,
			&cliente.UF, &cliente.TELE1, &cliente.CEL1,
			&cliente.CONTATO, &cliente.EMAIL, &cliente.PLACA,
			&cliente.NUMERO, &cliente.CIDIBGE, &cliente.COMPLE,
			&cliente.STATUS)
		clientes = append(clientes, cliente)
	}
	return clientes, nil
}

// DesmarcarAguardandoSyncLocal desmarca o campo AGUARDANDO_SYNC no cliente local
func (repositorio Clientes) DesmarcarAguardandoSyncLocal(clienteID uint64) error {
	_, erro := repositorio.db.Exec("UPDATE CLIENTE SET AGUARDANDO_SYNC = '' WHERE ID = ?", clienteID)
	return erro
}
