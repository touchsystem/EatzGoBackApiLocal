package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"errors"
	"time"
)

// Pagar representa o repositório de pagar
type Pagar struct {
	db *sql.DB
}

// NovoRepositorioDePagar cria um novo repositório de pagar
func NovoRepositorioDePagar(db *sql.DB) *Pagar {
	return &Pagar{db}
}

// CriarPagar insere um novo registro de pagar na tabela PAGAR e gera o lançamento duplo no CAIXA
func (repositorio Pagar) CriarPagar(pagar modelos.Pagar) error {
	// Iniciar transação
	tx, erro := repositorio.db.Begin()
	if erro != nil {
		return erro
	}

	// Inserir o registro de pagar
	statement, erro := tx.Prepare(
		`INSERT INTO PAGAR (DOC_L, DATA, DT_VEN, VL_TIT, VL_SAL, ID_FOR,  OBS, SITUA, CNPJ, id_usu) 
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if erro != nil {
		tx.Rollback()
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		pagar.DOC_L,
		pagar.DATA,
		pagar.DT_VEN,
		pagar.VL_TIT,
		pagar.VL_SAL,
		pagar.ID_FOR,
		//	pagar.DESCONTO,
		//	pagar.ACRESCIMO,
		pagar.OBS,
		pagar.SITUA,
		pagar.CNPJ,
		pagar.IDUsu,
	)
	if erro != nil {
		tx.Rollback()
		return erro
	}

	// Criar lançamentos duplos no caixa
	statementDebito, erro := tx.Prepare(
		`INSERT INTO CAIXA (NR_CAIXA, DOCUMENTO, DATA, VALORS, COMPLEMENT, id_usu) 
		 VALUES (?, ?, ?, ?, ?, ?)`)
	if erro != nil {
		tx.Rollback()
		return erro
	}
	defer statementDebito.Close()

	_, erro = statementDebito.Exec(
		pagar.ContraPartida, // Conta de débito
		pagar.DOC_L,         // Documento
		pagar.DATA,          // Data
		pagar.VL_TIT,        // Valor a débito
		pagar.OBS,           // Complemento
		pagar.IDUsu,         // ID do usuário
	)
	if erro != nil {
		tx.Rollback()
		return erro
	}

	// Crédito
	statementCredito, erro := tx.Prepare(
		`INSERT INTO CAIXA (NR_CAIXA, DOCUMENTO, DATA, VALOR, COMPLEMENT, id_usu) 
		 VALUES (?, ?, ?, ?, ?, ?)`)
	if erro != nil {
		tx.Rollback()
		return erro
	}
	defer statementCredito.Close()

	_, erro = statementCredito.Exec(
		"211101",     // Conta de crédito (fornecedores)
		pagar.DOC_L,  // Documento
		pagar.DATA,   // Data
		pagar.VL_TIT, // Valor a crédito
		pagar.OBS,    // Complemento
		pagar.IDUsu,  // ID do usuário
	)
	if erro != nil {
		tx.Rollback()
		return erro
	}

	// Confirmar transação
	if erro = tx.Commit(); erro != nil {
		return erro
	}

	return nil
}

// BuscarTituloPorID busca um título de pagamento pelo ID
func (repositorio Pagar) BuscarTituloPorID(idPagar int) (modelos.Pagar, error) {
	var titulo modelos.Pagar

	linha := repositorio.db.QueryRow("SELECT ID_PAGAR, DOC_L, VL_SAL FROM PAGAR WHERE ID_PAGAR = ?", idPagar)
	erro := linha.Scan(&titulo.ID_PAGAR, &titulo.DOC_L, &titulo.VL_SAL)

	if erro == sql.ErrNoRows {
		return titulo, errors.New("Título não encontrado")
	}

	if erro != nil {
		return titulo, erro
	}

	return titulo, nil
}
func (repositorio Pagar) QuitarTitulo(tituloExistente modelos.Pagar, pagar modelos.Pagar) error {
	// Iniciar transação
	tx, erro := repositorio.db.Begin()
	if erro != nil {
		return erro
	}

	// Atualizar o saldo do título no banco de dados
	statement, erro := tx.Prepare(
		`UPDATE PAGAR SET VL_SAL = ?, DT_PGTO = ? WHERE ID_PAGAR = ?`,
	)
	if erro != nil {
		tx.Rollback()
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		tituloExistente.VL_SAL,
		pagar.DT_PGTO,

		pagar.ID_PAGAR,
	)
	if erro != nil {
		tx.Rollback()
		return erro
	}

	// Lançamento para o pagamento (normal)
	statementDebito, erro := tx.Prepare(
		`INSERT INTO CAIXA (NR_CAIXA, DOCUMENTO, DATA, VALORS , COMPLEMENT, id_usu) 
         VALUES (?, ?, ?, ?, ?, ?)`)
	if erro != nil {
		tx.Rollback()
		return erro
	}
	defer statementDebito.Close()

	_, erro = statementDebito.Exec(
		"211101",                        // CONTA FIXA
		tituloExistente.DOC_L,           // Documento
		pagar.DT_PGTO,                   // Data de pagamento
		(pagar.VL_TIT + pagar.DESCONTO), // Valor pago
		pagar.OBS,                       // Complemento
		pagar.IDUsu,                     // ID do usuário
	)
	if erro != nil {
		tx.Rollback()
		return erro
	}

	// Lançamento de crédito para pagamento
	statementCredito, erro := tx.Prepare(
		`INSERT INTO CAIXA (NR_CAIXA, DOCUMENTO, DATA, VALOR, COMPLEMENT, id_usu) 
         VALUES (?, ?, ?, ?, ?, ?)`)
	if erro != nil {
		tx.Rollback()
		return erro
	}
	defer statementCredito.Close()

	_, erro = statementCredito.Exec(
		pagar.ContraPartida,             // CONTA DESPESAS
		tituloExistente.DOC_L,           // Documento
		pagar.DT_PGTO,                   // Data de pagamento
		(pagar.VL_TIT + pagar.DESCONTO), // Valor pago
		pagar.OBS,                       // Complemento
		pagar.IDUsu,                     // ID do usuário
	)
	if erro != nil {
		tx.Rollback()
		return erro
	}

	// Verificar e lançar o desconto (se aplicável)
	if pagar.DESCONTO > 0 {
		// Débito na conta informada
		_, erro = statementDebito.Exec(
			pagar.ContraPartida,   // Conta de débito
			tituloExistente.DOC_L, // Documento
			pagar.DT_PGTO,         // Data de pagamento
			pagar.DESCONTO,        // Valor do desconto
			"Desconto Contas a Pagar Fornecedor "+pagar.CNPJ, // Observação
			pagar.IDUsu, // ID do usuário
		)
		if erro != nil {
			tx.Rollback()
			return erro
		}

		// Crédito na conta fixa de desconto
		_, erro = statementCredito.Exec(
			"312201",              // Conta de crédito para descontos
			tituloExistente.DOC_L, // Documento
			pagar.DT_PGTO,         // Data de pagamento
			pagar.DESCONTO,        // Valor do desconto
			"Desconto Contas a Pagar Fornecedor "+pagar.CNPJ, // Observação
			pagar.IDUsu, // ID do usuário
		)
		if erro != nil {
			tx.Rollback()
			return erro
		}
	}

	// Verificar e lançar o acréscimo (se aplicável)
	if pagar.ACRESCIMO > 0 {
		// Débito na conta fixa de acréscimo
		_, erro = statementDebito.Exec(
			"331106",                             // Conta de débito para acréscimos
			tituloExistente.DOC_L,                // Documento
			pagar.DT_PGTO,                        // Data de pagamento
			pagar.ACRESCIMO,                      // Valor do acréscimo
			"Juros dados Fornecedor "+pagar.CNPJ, // Observação
			pagar.IDUsu,                          // ID do usuário
		)
		if erro != nil {
			tx.Rollback()
			return erro
		}

		// Crédito na conta informada no JSON
		_, erro = statementCredito.Exec(
			pagar.ContraPartida,                  // Conta de crédito informada no JSON
			tituloExistente.DOC_L,                // Documento
			pagar.DT_PGTO,                        // Data de pagamento
			pagar.ACRESCIMO,                      // Valor do acréscimo
			"Juros dados Fornecedor "+pagar.CNPJ, // Observação
			pagar.IDUsu,                          // ID do usuário
		)
		if erro != nil {
			tx.Rollback()
			return erro
		}
	}

	// Confirmar transação
	if erro = tx.Commit(); erro != nil {
		return erro
	}

	return nil
}

// BuscarTitulosPorFornecedor busca os títulos de um fornecedor por período e situação (em aberto ou pagos)
func (repositorio Pagar) BuscarTitulosPorFornecedor(idFornecedor int, dataInicio string, dataFim string, situacao string) ([]modelos.Pagar, error) {
	var rows *sql.Rows
	var erro error

	// Montar a query básica com os campos comuns
	query := `SELECT 
			p.ID_PAGAR, 
			p.DOC_L, 
			p.DATA, 
			p.DT_VEN, 
			p.VL_TIT, 
			p.VL_SAL, 
			p.ID_FOR, 
			p.OBS, 
			f.CNPJ_CPF AS CNPJ, 
			f.NOME AS NOME_FORNECEDOR`

	// Incluir `DT_PGTO` na query quando a situação for "pago"
	if situacao == "pago" {
		query += `, p.DT_PGTO`
	}

	query += ` FROM 
			PAGAR p
		JOIN 
			FORNECEDOR f ON p.ID_FOR = f.ID 
		WHERE 
			p.ID_FOR = ? 
			AND p.DATA BETWEEN ? AND ? 
			AND p.VL_SAL`

	// Ajustar a condição de acordo com a situação
	if situacao == "aberto" {
		query += " > 0"
	} else if situacao == "pago" {
		query += " = 0"
	} else {
		return nil, errors.New("situação inválida: use 'aberto' para títulos com saldo ou 'pago' para quitados")
	}

	// Executar a consulta com os parâmetros
	rows, erro = repositorio.db.Query(query, idFornecedor, dataInicio, dataFim)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var titulos []modelos.Pagar
	for rows.Next() {
		var titulo modelos.Pagar

		// Realizar o `Scan` com `DT_PGTO` se a situação for "pago"
		if situacao == "pago" {
			if erro = rows.Scan(
				&titulo.ID_PAGAR,
				&titulo.DOC_L,
				&titulo.DATA,
				&titulo.DT_VEN,
				&titulo.VL_TIT,
				&titulo.VL_SAL,
				&titulo.ID_FOR,
				&titulo.OBS,
				&titulo.CNPJ,
				&titulo.NOME_FORNECEDOR,
				&titulo.DT_PGTO, // Inclui `DT_PGTO` apenas para "pago"
			); erro != nil {
				return nil, erro
			}
		} else {
			// Excluir `DT_PGTO` do `Scan` se a situação for "aberto"
			if erro = rows.Scan(
				&titulo.ID_PAGAR,
				&titulo.DOC_L,
				&titulo.DATA,
				&titulo.DT_VEN,
				&titulo.VL_TIT,
				&titulo.VL_SAL,
				&titulo.ID_FOR,
				&titulo.OBS,
				&titulo.CNPJ,
				&titulo.NOME_FORNECEDOR,
			); erro != nil {
				return nil, erro
			}
		}
		titulos = append(titulos, titulo)
	}

	return titulos, nil
}

// BuscarTitulosPorPeriodo busca os títulos de todos os fornecedores por período e situação (em aberto ou pagos)
func (repositorio Pagar) BuscarTitulosPorPeriodo(dataInicio string, dataFim string, situacao string) ([]modelos.Pagar, error) {
	// Validar formato das datas
	_, err := time.Parse("2006-01-02", dataInicio)
	if err != nil {
		return nil, errors.New("data_inicio inválida, use o formato YYYY-MM-DD")
	}
	_, err = time.Parse("2006-01-02", dataFim)
	if err != nil {
		return nil, errors.New("data_fim inválida, use o formato YYYY-MM-DD")
	}

	// Montar a query baseada na situação
	query := `SELECT 
			p.ID_PAGAR, 
			p.DOC_L, 
			p.DATA, 
			p.DT_VEN, 
			p.VL_TIT, 
			p.VL_SAL, 
			p.ID_FOR, 
			p.OBS, 
			f.CNPJ_CPF AS CNPJ, 
			f.NOME AS NOME_FORNECEDOR`

	// Ajustar a query para incluir `DT_PGTO` quando a situação for "pago"
	if situacao == "pago" {
		query += `, p.DT_PGTO`
	}

	query += ` FROM 
			PAGAR p
		JOIN 
			FORNECEDOR f ON p.ID_FOR = f.ID 
		WHERE 
			p.DATA BETWEEN ? AND ? 
			AND p.VL_SAL`

	// Ajustar a condição de acordo com a situação
	if situacao == "aberto" {
		query += " > 0"
	} else if situacao == "pago" {
		query += " = 0"
	} else {
		return nil, errors.New("situação inválida: use 'aberto' para títulos com saldo ou 'pago' para quitados")
	}

	// Executar a consulta com os parâmetros
	rows, erro := repositorio.db.Query(query, dataInicio, dataFim)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	// Mapear resultados para a estrutura de retorno
	var titulos []modelos.Pagar
	for rows.Next() {
		var titulo modelos.Pagar
		if situacao == "pago" {
			// Incluir `DT_PGTO` no `Scan` quando a situação for "pago"
			if erro = rows.Scan(
				&titulo.ID_PAGAR,
				&titulo.DOC_L,
				&titulo.DATA,
				&titulo.DT_VEN,
				&titulo.VL_TIT,
				&titulo.VL_SAL,
				&titulo.ID_FOR,
				&titulo.OBS,
				&titulo.CNPJ,
				&titulo.NOME_FORNECEDOR,
				&titulo.DT_PGTO,
			); erro != nil {
				return nil, erro
			}
		} else {
			// Excluir `DT_PGTO` do `Scan` quando a situação for "aberto"
			if erro = rows.Scan(
				&titulo.ID_PAGAR,
				&titulo.DOC_L,
				&titulo.DATA,
				&titulo.DT_VEN,
				&titulo.VL_TIT,
				&titulo.VL_SAL,
				&titulo.ID_FOR,
				&titulo.OBS,
				&titulo.CNPJ,
				&titulo.NOME_FORNECEDOR,
			); erro != nil {
				return nil, erro
			}
		}
		titulos = append(titulos, titulo)
	}

	return titulos, nil
}
