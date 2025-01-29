package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

// Clientes representa um repositório de clientes
type Clientes struct {
	db *sql.DB
}

// NovoRepositorioDeClientes cria um repositório de clientes
func NovoRepositorioDeClientes(db *sql.DB) *Clientes {
	return &Clientes{db}
}

// Criar insere um cliente no banco de dados
func (repositorio Clientes) Criar(cliente modelos.Cliente) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"insert into CLIENTE (CNPJ_CPF, NOME, FANTASIA, ENDERE, CIDADE, BAIRRO, CEP, UF, TELE1, CEL1,  CONTATO, EMAIL, GRUPO, PLACA, NUMERO, CIDIBGE, COMPLE) values( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)

	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(cliente.CNPJ_CPF, cliente.NOME, cliente.FANTASIA, cliente.ENDERE, cliente.CIDADE, cliente.BAIRRO, cliente.CEP, cliente.UF, cliente.TELE1, cliente.CEL1, cliente.CONTATO, cliente.EMAIL, cliente.GRUPO, cliente.PLACA, cliente.NUMERO, cliente.CIDIBGE, cliente.COMPLE)
	if erro != nil {
		fmt.Println(erro)
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// Buscar traz todos os clientes que atendem um filtro de nome
func (repositorio Clientes) Buscar(nome string) ([]modelos.Cliente, error) {
	//	nome = "%" + nome + "%"
	nome = fmt.Sprintf("%%%s%%", nome) // %nomeOuNick%

	//	linhas, erro := repositorio.db.Query(
	//		"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?",
	//		nomeOuNick, nomeOuNick,
	//	)

	rows, erro := repositorio.db.Query("SELECT ID, CNPJ_CPF, NOME,FANTASIA, 	ENDERE, 	CIDADE, 	BAIRRO, 	CEP, 	UF, 	TELE1, 	 	CEL1, 	CONTATO, EMAIL, PLACA, NUMERO, CIDIBGE, COMPLE FROM CLIENTE WHERE NOME LIKE ?", nome)

	if erro != nil {

		return nil, erro
	}
	defer rows.Close()

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
		); erro != nil {
			println(erro)
			return nil, erro

		}
		clientes = append(clientes, cliente)
	}

	return clientes, nil
}

// BuscarPorID busca um cliente no banco de dados pelo ID
func (repositorio Clientes) BuscarPorID(clienteID uint64) (modelos.Cliente, error) {
	row := repositorio.db.QueryRow(`SELECT 
    ID, 
    CNPJ_CPF, 
    NOME, 
    FANTASIA, 
    ENDERE, 
    CIDADE, 
    BAIRRO, 
    CEP, 
    UF, 
    TELE1, 
    CEL1, 
    CONTATO, 
    EMAIL, 
    GRUPO, 
    PLACA, 
    NUMERO, 
    CIDIBGE, 
    COMPLE 
FROM CLIENTE 
WHERE ID = ?`, clienteID)

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
		&cliente.GRUPO,
		&cliente.PLACA,
		&cliente.NUMERO,
		&cliente.CIDIBGE,
		&cliente.COMPLE,
	); erro != nil {
		return modelos.Cliente{}, erro
	}

	return cliente, nil
}

// Atualizar altera um cliente no banco de dados
func (repositorio Clientes) Atualizar(clienteID uint64, cliente modelos.Cliente) error {
	statement, erro := repositorio.db.Prepare(
		"update CLIENTE set CNPJ_CPF = ?, NOME = ?, FANTASIA = ?, ENDERE = ?, CIDADE = ?, BAIRRO = ?, CEP = ?, UF = ?, TELE1 = ?, CEL1 = ?, CONTATO = ?, EMAIL = ?,  PLACA = ?, NUMERO = ?, CIDIBGE = ?, COMPLE = ? where ID = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(cliente.CNPJ_CPF, cliente.NOME, cliente.FANTASIA, cliente.ENDERE, cliente.CIDADE, cliente.BAIRRO, cliente.CEP, cliente.UF, cliente.TELE1, cliente.CEL1, cliente.CONTATO, cliente.EMAIL, cliente.PLACA, cliente.NUMERO, cliente.CIDIBGE, cliente.COMPLE, clienteID)
	return erro
}

// Deletar exclui um cliente no banco de dados
func (repositorio Clientes) Deletar(clienteID uint64) error {
	statement, erro := repositorio.db.Prepare("delete from CLIENTE where ID = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(clienteID)
	return erro
}
