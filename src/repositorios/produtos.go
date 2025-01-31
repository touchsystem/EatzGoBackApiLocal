package repositorios

import (
	"api/src/modelos"
	"api/src/utils"
	"database/sql"
	"fmt"
	"log"
)

// Produtos representa um repositório de produtos
type Produtos struct {
	db *sql.DB
}

// NovoRepositorioDeProdutos cria um repositório de produtos
func NovoRepositorioDeProdutos(db *sql.DB) *Produtos {
	return &Produtos{db}
}

// CriarProdutos insere um produto no banco de dados local e sincroniza com a nuvem
func (repositorio Produtos) CriarProdutos(produto modelos.Produto) (uint64, error) {
	// Prepara e executa o SQL de inserção

	statement, erro := repositorio.db.Prepare(`
		INSERT INTO PRODUTO (CODM, DES2, PV, IMPRE, FORNE, PS_L, PS_BR, BARRA, STATUS, UND, GRUPO, CNAE, CD_SERV, 
		CODPISCOFINS, SERIAL, NCM, EAN, CSON, ORIGEN, ALICOTA, DES2_COMPLE, ST_PROMO1, PC_PROMO, PROM1_HI, 
		PROM1_HF, PROM1_SEMANA, PROMO_CONSUMO, PROMO_PAGAR, ID_HOSTWEB)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	// Insere os valores para cada coluna do produto
	resultado, erro := statement.Exec(
		produto.CODM, produto.DES2, produto.PV, produto.IMPRE, produto.FORNE, produto.PS_L, produto.PS_BR,
		produto.BARRA, produto.STATUS, produto.UND, produto.GRUPO, produto.CNAE, produto.CD_SERV,
		produto.CODPISCOFINS, produto.SERIAL, produto.NCM, produto.EAN, produto.CSON, produto.ORIGEN,
		produto.ALICOTA, produto.DES2_COMPLE, produto.ST_PROMO1, produto.PC_PROMO, produto.PROM1_HI,
		produto.PROM1_HF, produto.PROM1_SEMANA, produto.PROMO_CONSUMO, produto.PROMO_PAGAR, produto.ID_HOSTLOCAL,
	)
	if erro != nil {
		return 0, erro
	}

	// Obtém o último ID inserido
	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	// Define o ID do produto com o valor gerado
	produto.ID = uint64(ultimoIDInserido)

	// Sincroniza com a nuvem e captura o ID gerado na nuvem
	idHostWeb, erro := utils.SincronizarProdutoNuvem(produto)
	if erro != nil {
		log.Printf("Erro... ao sincronizar com a nuvem: %v", erro)
		return uint64(ultimoIDInserido), nil
	}

	// Atualiza o campo ID_HOSTWEB com o ID retornado da nuvem
	_, erro = repositorio.db.Exec("UPDATE PRODUTO SET ID_HOSTWEB = ? WHERE ID = ?", idHostWeb, ultimoIDInserido)
	if erro != nil {
		return uint64(ultimoIDInserido), fmt.Errorf("Erro ao atualizar ID_HOSTWEB: %v", erro)
	}

	return uint64(ultimoIDInserido), nil
}

// BuscarProdutos traz produtos com filtro de status e filtros LIKE em DES2 e CODM
func (repositorio Produtos) BuscarProdutos(status string, filtro string) ([]modelos.Produto, error) {
	// Construir a query base com a condição de status
	query := `SELECT ID, CODM, DES2, PV, IMPRE, FORNE, PS_L, PS_BR, BARRA, STATUS, UND, GRUPO, CNAE, 
			  CD_SERV, CODPISCOFINS, SERIAL, NCM, EAN, CSON, ORIGEN, ALICOTA, DES2_COMPLE, ST_PROMO1, 
			  PC_PROMO, PROM1_HI, PROM1_HF, PROM1_SEMANA, PROMO_CONSUMO, PROMO_PAGAR, ID_HOSTLOCAL 
			  FROM PRODUTO WHERE STATUS = ?`

	// Adicionar condição de filtro para os campos DES2 e CODM
	var args []interface{}
	args = append(args, status)
	if filtro != "" {
		filtro = fmt.Sprintf("%%%s%%", filtro) // Adicionar % para o LIKE
		query += " AND (DES2 LIKE ? OR CODM LIKE ?)"
		args = append(args, filtro, filtro)
	}
	//fmt.Println(args, filtro)
	// Execute a consulta com filtros aplicados
	rows, erro := repositorio.db.Query(query, args...)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var produtos []modelos.Produto
	for rows.Next() {
		var produto modelos.Produto
		if erro = rows.Scan(
			&produto.ID,
			&produto.CODM,
			&produto.DES2,
			&produto.PV,
			&produto.IMPRE,
			&produto.FORNE,
			&produto.PS_L,
			&produto.PS_BR,
			&produto.BARRA,
			&produto.STATUS,
			&produto.UND,
			&produto.GRUPO,
			&produto.CNAE,
			&produto.CD_SERV,
			&produto.CODPISCOFINS,
			&produto.SERIAL,
			&produto.NCM,
			&produto.EAN,
			&produto.CSON,
			&produto.ORIGEN,
			&produto.ALICOTA,
			&produto.DES2_COMPLE,
			&produto.ST_PROMO1,
			&produto.PC_PROMO,
			&produto.PROM1_HI,
			&produto.PROM1_HF,
			&produto.PROM1_SEMANA,
			&produto.PROMO_CONSUMO,
			&produto.PROMO_PAGAR,
			&produto.ID_HOSTLOCAL,
		); erro != nil {
			return nil, erro
		}
		produtos = append(produtos, produto)
	}

	return produtos, nil
}

// BuscarPorID busca um produto no banco de dados pelo ID
func (repositorio Produtos) BuscarProdutoPorID(produtoID uint64) (modelos.Produto, error) {
	row := repositorio.db.QueryRow("SELECT ID,  CODM, DES1, DES2, PV, PV2, IMPRE, FORNE, CST1, CST2, STOCK1, STOCK2, STOCK3, STOCK4, COMPOSI, MINIMO, PS_L, PS_BR, COM_1, COM_2, BARRA, OB1, OB2, OB3, OB4, OB5, OBS1, OBS2, OBS3, OBS4, OBS5, STATUS, UND, CARDAPIO, GRUPO, FISCAL, FISCAL1, BAIXAD, SALDO_AN, CUSTO_CM, CUSTO_MD_AN, CUSTO_MD_AT, CNAE, CD_SERV, COMISSAO, CODPISCOFINS, SERIAL, QPCX, PV3, NCM, EAN, CSON, ORIGEN, ALICOTA, SERVICO, DES2_COMPLE, ST2, ST_PROMO1, PC_PROMO, PROM1_HI, PROM1_HF, PROM1_SEMANA, PROMO_CONSUMO, PROMO_PAGAR FROM PRODUTO WHERE ID = ?", produtoID)

	var produto modelos.Produto
	if erro := row.Scan(
		&produto.ID,
		&produto.CODM,
		&produto.DES1,
		&produto.DES2,
		&produto.PV,
		&produto.PV2,
		&produto.IMPRE,
		&produto.FORNE,
		&produto.CST1,
		&produto.CST2,
		&produto.STOCK1,
		&produto.STOCK2,
		&produto.STOCK3,
		&produto.STOCK4,
		&produto.COMPOSI,
		&produto.MINIMO,
		&produto.PS_L,
		&produto.PS_BR,
		&produto.COM_1,
		&produto.COM_2,
		&produto.BARRA,
		&produto.OB1,
		&produto.OB2,
		&produto.OB3,
		&produto.OB4,
		&produto.OB5,
		&produto.OBS1,
		&produto.OBS2,
		&produto.OBS3,
		&produto.OBS4,
		&produto.OBS5,
		&produto.STATUS,
		&produto.UND,
		&produto.CARDAPIO,
		&produto.GRUPO,
		&produto.FISCAL,
		&produto.FISCAL1,
		&produto.BAIXAD,
		&produto.SALDO_AN,
		&produto.CUSTO_CM,
		&produto.CUSTO_MD_AN,
		&produto.CUSTO_MD_AT,
		&produto.CNAE,
		&produto.CD_SERV,
		&produto.COMISSAO,
		&produto.CODPISCOFINS,
		&produto.SERIAL,
		&produto.QPCX,
		&produto.PV3,
		&produto.NCM,
		&produto.EAN,
		&produto.CSON,
		&produto.ORIGEN,
		&produto.ALICOTA,
		&produto.SERVICO,
		&produto.DES2_COMPLE,
		&produto.ST2,
		&produto.ST_PROMO1,
		&produto.PC_PROMO,
		&produto.PROM1_HI,
		&produto.PROM1_HF,
		&produto.PROM1_SEMANA,
		&produto.PROMO_CONSUMO,
		&produto.PROMO_PAGAR,
	); erro != nil {
		return modelos.Produto{}, erro
	}

	return produto, nil
}

// Atualizar altera um produto no banco de dados
func (repositorio Produtos) AtualizarProdutos(produtoID uint64, produto modelos.Produto) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE PRODUTO SET AGUARDANDO_SYNC = ?,  DES1 = ?, DES2 = ?, PV = ?, PV2 = ?, IMPRE = ?, FORNE = ?, CST1 = ?, CST2 = ?, STOCK1 = ?, STOCK2 = ?, STOCK3 = ?, STOCK4 = ?, COMPOSI = ?, MINIMO = ?, PS_L = ?, PS_BR = ?, COM_1 = ?, COM_2 = ?, BARRA = ?, OB1 = ?, OB2 = ?, OB3 = ?, OB4 = ?, OB5 = ?, OBS1 = ?, OBS2 = ?, OBS3 = ?, OBS4 = ?, OBS5 = ?, STATUS = ?, UND = ?, CARDAPIO = ?, GRUPO = ?, FISCAL = ?, FISCAL1 = ?, BAIXAD = ?, SALDO_AN = ?, CUSTO_CM = ?, CUSTO_MD_AN = ?, CUSTO_MD_AT = ?, CNAE = ?, CD_SERV = ?, COMISSAO = ?, CODPISCOFINS = ?, SERIAL = ?, QPCX = ?, PV3 = ?, NCM = ?, EAN = ?, CSON = ?, ORIGEN = ?, ALICOTA = ?, SERVICO = ?, DES2_COMPLE = ?, ST2 = ?, ST_PROMO1 = ?, PC_PROMO = ?, PROM1_HI = ?, PROM1_HF = ?, PROM1_SEMANA = ?, PROMO_CONSUMO = ?, PROMO_PAGAR = ? WHERE ID = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()
	//fmt.Println(produtoID)
	_, erro = statement.Exec(produto.AGUARDANDO_SYNC, produto.DES1, produto.DES2, produto.PV, produto.PV2, produto.IMPRE, produto.FORNE, produto.CST1, produto.CST2, produto.STOCK1, produto.STOCK2, produto.STOCK3, produto.STOCK4, produto.COMPOSI, produto.MINIMO, produto.PS_L, produto.PS_BR, produto.COM_1, produto.COM_2, produto.BARRA, produto.OB1, produto.OB2, produto.OB3, produto.OB4, produto.OB5, produto.OBS1, produto.OBS2, produto.OBS3, produto.OBS4, produto.OBS5, produto.STATUS, produto.UND, produto.CARDAPIO, produto.GRUPO, produto.FISCAL, produto.FISCAL1, produto.BAIXAD, produto.SALDO_AN, produto.CUSTO_CM, produto.CUSTO_MD_AN, produto.CUSTO_MD_AT, produto.CNAE, produto.CD_SERV, produto.COMISSAO, produto.CODPISCOFINS, produto.SERIAL, produto.QPCX, produto.PV3, produto.NCM, produto.EAN, produto.CSON, produto.ORIGEN, produto.ALICOTA, produto.SERVICO, produto.DES2_COMPLE, produto.ST2, produto.ST_PROMO1, produto.PC_PROMO, produto.PROM1_HI, produto.PROM1_HF, produto.PROM1_SEMANA, produto.PROMO_CONSUMO, produto.PROMO_PAGAR, produtoID)
	return erro
}

// DeletarProduto altera um produto no banco de dados, marcando-o como deletado
func (repositorio Produtos) DeletarProduto(produtoID uint64, produto modelos.Produto) error {
	statement, erro := repositorio.db.Prepare(
		`UPDATE PRODUTO SET AGUARDANDO_SYNC = ?, STATUS = ? WHERE ID = ?`,
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	// Atualizar o status do produto para "deletado" (D) e marcar como "aguardando sincronização" (S)
	_, erro = statement.Exec("S", "D", produtoID)
	return erro
}

//PARTE NOVA SYNC

// Função para buscar produtos não sincronizados (ID_HOSTWEB = 0)
func (repositorio Produtos) BuscarProdutosNaoSincronizados() ([]modelos.Produto, error) {
	rows, erro := repositorio.db.Query(`
		SELECT ID, ID_HOSTWEB, CODM, DES1, DES2, PV, PV2, IMPRE, FORNE, CST1, CST2, STOCK1, STOCK2, STOCK3, STOCK4, 
			   COMPOSI, MINIMO, PS_L, PS_BR, COM_1, COM_2, BARRA, OB1, OB2, OB3, OB4, OB5, OBS1, OBS2, OBS3, OBS4, OBS5, 
			   STATUS, UND, CARDAPIO, GRUPO, FISCAL, FISCAL1, BAIXAD, SALDO_AN, CUSTO_CM, CUSTO_MD_AN, CUSTO_MD_AT, 
			   CNAE, CD_SERV, COMISSAO, CODPISCOFINS, SERIAL, QPCX, PV3, NCM, EAN, CSON, ORIGEN, ALICOTA, SERVICO, 
			   DES2_COMPLE, ST2, ST_PROMO1, PC_PROMO, PROM1_HI, PROM1_HF, PROM1_SEMANA, PROMO_CONSUMO, PROMO_PAGAR
		FROM PRODUTO WHERE ID_HOSTWEB = 0
	`)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var produtos []modelos.Produto
	for rows.Next() {
		var produto modelos.Produto
		erro := rows.Scan(
			&produto.ID,
			&produto.ID_HOSTWEB,
			&produto.CODM,
			&produto.DES1,
			&produto.DES2,
			&produto.PV,
			&produto.PV2,
			&produto.IMPRE,
			&produto.FORNE,
			&produto.CST1,
			&produto.CST2,
			&produto.STOCK1,
			&produto.STOCK2,
			&produto.STOCK3,
			&produto.STOCK4,
			&produto.COMPOSI,
			&produto.MINIMO,
			&produto.PS_L,
			&produto.PS_BR,
			&produto.COM_1,
			&produto.COM_2,
			&produto.BARRA,
			&produto.OB1,
			&produto.OB2,
			&produto.OB3,
			&produto.OB4,
			&produto.OB5,
			&produto.OBS1,
			&produto.OBS2,
			&produto.OBS3,
			&produto.OBS4,
			&produto.OBS5,
			&produto.STATUS,
			&produto.UND,
			&produto.CARDAPIO,
			&produto.GRUPO,
			&produto.FISCAL,
			&produto.FISCAL1,
			&produto.BAIXAD,
			&produto.SALDO_AN,
			&produto.CUSTO_CM,
			&produto.CUSTO_MD_AN,
			&produto.CUSTO_MD_AT,
			&produto.CNAE,
			&produto.CD_SERV,
			&produto.COMISSAO,
			&produto.CODPISCOFINS,
			&produto.SERIAL,
			&produto.QPCX,
			&produto.PV3,
			&produto.NCM,
			&produto.EAN,
			&produto.CSON,
			&produto.ORIGEN,
			&produto.ALICOTA,
			&produto.SERVICO,
			&produto.DES2_COMPLE,
			&produto.ST2,
			&produto.ST_PROMO1,
			&produto.PC_PROMO,
			&produto.PROM1_HI,
			&produto.PROM1_HF,
			&produto.PROM1_SEMANA,
			&produto.PROMO_CONSUMO,
			&produto.PROMO_PAGAR,
		)
		if erro != nil {
			return nil, erro
		}
		produtos = append(produtos, produto)
	}

	return produtos, nil
}

// Função para verificar se o produto já existe na nuvem
func (repositorio Produtos) ProdutoExisteNaNuvem(idHostLocal uint64) bool {
	var existe int
	erro := repositorio.db.QueryRow("SELECT COUNT(*) FROM PRODUTO WHERE ID_HOSTLOCAL = ?", idHostLocal).Scan(&existe)
	return erro == nil && existe > 0
}

// Função para atualizar um produto existente pelo ID_HOSTWEB
func (repositorio Produtos) AtualizarProdutoPorHostWeb(idHostWeb uint64, produto modelos.Produto) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE PRODUTO SET 
			ID_HOSTLOCAL = ?, CODM = ?, DES1 = ?, DES2 = ?, PV = ?, PV2 = ?, IMPRE = ?, FORNE = ?, CST1 = ?, CST2 = ?, 
			STOCK1 = ?, STOCK2 = ?, STOCK3 = ?, STOCK4 = ?, COMPOSI = ?, MINIMO = ?, PS_L = ?, PS_BR = ?, 
			COM_1 = ?, COM_2 = ?, BARRA = ?, OB1 = ?, OB2 = ?, OB3 = ?, OB4 = ?, OB5 = ?, OBS1 = ?, OBS2 = ?, 
			OBS3 = ?, OBS4 = ?, OBS5 = ?, STATUS = ?, UND = ?, CARDAPIO = ?, GRUPO = ?, FISCAL = ?, FISCAL1 = ?, 
			BAIXAD = ?, SALDO_AN = ?, CUSTO_CM = ?, CUSTO_MD_AN = ?, CUSTO_MD_AT = ?, CNAE = ?, CD_SERV = ?, 
			COMISSAO = ?, CODPISCOFINS = ?, SERIAL = ?, QPCX = ?, PV3 = ?, NCM = ?, EAN = ?, CSON = ?, ORIGEN = ?, 
			ALICOTA = ?, SERVICO = ?, DES2_COMPLE = ?, ST2 = ?, ST_PROMO1 = ?, PC_PROMO = ?, PROM1_HI = ?, 
			PROM1_HF = ?, PROM1_SEMANA = ?, PROMO_CONSUMO = ?, PROMO_PAGAR = ?, CRIADOEM = ?
		WHERE ID_HOSTWEB = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		produto.ID_HOSTLOCAL,
		produto.CODM,
		produto.DES1,
		produto.DES2,
		produto.PV,
		produto.PV2,
		produto.IMPRE,
		produto.FORNE,
		produto.CST1,
		produto.CST2,
		produto.STOCK1,
		produto.STOCK2,
		produto.STOCK3,
		produto.STOCK4,
		produto.COMPOSI,
		produto.MINIMO,
		produto.PS_L,
		produto.PS_BR,
		produto.COM_1,
		produto.COM_2,
		produto.BARRA,
		produto.OB1,
		produto.OB2,
		produto.OB3,
		produto.OB4,
		produto.OB5,
		produto.OBS1,
		produto.OBS2,
		produto.OBS3,
		produto.OBS4,
		produto.OBS5,
		produto.STATUS,
		produto.UND,
		produto.CARDAPIO,
		produto.GRUPO,
		produto.FISCAL,
		produto.FISCAL1,
		produto.BAIXAD,
		produto.SALDO_AN,
		produto.CUSTO_CM,
		produto.CUSTO_MD_AN,
		produto.CUSTO_MD_AT,
		produto.CNAE,
		produto.CD_SERV,
		produto.COMISSAO,
		produto.CODPISCOFINS,
		produto.SERIAL,
		produto.QPCX,
		produto.PV3,
		produto.NCM,
		produto.EAN,
		produto.CSON,
		produto.ORIGEN,
		produto.ALICOTA,
		produto.SERVICO,
		produto.DES2_COMPLE,
		produto.ST2,
		produto.ST_PROMO1,
		produto.PC_PROMO,
		produto.PROM1_HI,
		produto.PROM1_HF,
		produto.PROM1_SEMANA,
		produto.PROMO_CONSUMO,
		produto.PROMO_PAGAR,
		produto.CRIADOEM,
		idHostWeb,
	)
	return erro
}

// Atualiza o ID_HOSTWEB do produto após sincronização bem-sucedida com a nuvem
func (repositorio Produtos) AtualizarProdutoHostWeb(idLocal uint64, idNuvem uint64) error {
	_, erro := repositorio.db.Exec(`
		UPDATE PRODUTO SET ID_HOSTWEB = ? WHERE ID = ?
	`, idNuvem, idLocal)
	return erro
}

// Função para verificar se o produto já existe localmente
func (repositorio Produtos) ProdutoExisteNoLocal(idHostWeb uint64) bool {
	var existe int
	erro := repositorio.db.QueryRow("SELECT COUNT(*) FROM PRODUTO WHERE ID_HOSTWEB = ?", idHostWeb).Scan(&existe)
	return erro == nil && existe > 0
}

// Função para atualizar um produto existente pelo ID_HOSTLOCAL
func (repositorio Produtos) AtualizarProdutoPorHostLocal(idHostLocal uint64, produto modelos.Produto) error {
	statement, erro := repositorio.db.Prepare(`
		UPDATE PRODUTO SET 
			ID_HOSTWEB = ?, AGUARDANDO_SYNC = ?, CODM = ?, DES1 = ?, DES2 = ?, PV = ?, PV2 = ?, IMPRE = ?, 
			FORNE = ?, CST1 = ?, CST2 = ?, STOCK1 = ?, STOCK2 = ?, STOCK3 = ?, STOCK4 = ?, COMPOSI = ?, 
			MINIMO = ?, PS_L = ?, PS_BR = ?, COM_1 = ?, COM_2 = ?, BARRA = ?, OB1 = ?, OB2 = ?, OB3 = ?, 
			OB4 = ?, OB5 = ?, OBS1 = ?, OBS2 = ?, OBS3 = ?, OBS4 = ?, OBS5 = ?, STATUS = ?, UND = ?, 
			CARDAPIO = ?, GRUPO = ?, FISCAL = ?, FISCAL1 = ?, BAIXAD = ?, SALDO_AN = ?, CUSTO_CM = ?, 
			CUSTO_MD_AN = ?, CUSTO_MD_AT = ?, CNAE = ?, CD_SERV = ?, COMISSAO = ?, CODPISCOFINS = ?, 
			SERIAL = ?, QPCX = ?, PV3 = ?, NCM = ?, EAN = ?, CSON = ?, ORIGEN = ?, ALICOTA = ?, 
			SERVICO = ?, DES2_COMPLE = ?, ST2 = ?, ST_PROMO1 = ?, PC_PROMO = ?, PROM1_HI = ?, 
			PROM1_HF = ?, PROM1_SEMANA = ?, PROMO_CONSUMO = ?, PROMO_PAGAR = ?, CRIADOEM = ?
		WHERE ID_HOSTLOCAL = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(
		produto.ID_HOSTWEB,
		produto.AGUARDANDO_SYNC,
		produto.CODM,
		produto.DES1,
		produto.DES2,
		produto.PV,
		produto.PV2,
		produto.IMPRE,
		produto.FORNE,
		produto.CST1,
		produto.CST2,
		produto.STOCK1,
		produto.STOCK2,
		produto.STOCK3,
		produto.STOCK4,
		produto.COMPOSI,
		produto.MINIMO,
		produto.PS_L,
		produto.PS_BR,
		produto.COM_1,
		produto.COM_2,
		produto.BARRA,
		produto.OB1,
		produto.OB2,
		produto.OB3,
		produto.OB4,
		produto.OB5,
		produto.OBS1,
		produto.OBS2,
		produto.OBS3,
		produto.OBS4,
		produto.OBS5,
		produto.STATUS,
		produto.UND,
		produto.CARDAPIO,
		produto.GRUPO,
		produto.FISCAL,
		produto.FISCAL1,
		produto.BAIXAD,
		produto.SALDO_AN,
		produto.CUSTO_CM,
		produto.CUSTO_MD_AN,
		produto.CUSTO_MD_AT,
		produto.CNAE,
		produto.CD_SERV,
		produto.COMISSAO,
		produto.CODPISCOFINS,
		produto.SERIAL,
		produto.QPCX,
		produto.PV3,
		produto.NCM,
		produto.EAN,
		produto.CSON,
		produto.ORIGEN,
		produto.ALICOTA,
		produto.SERVICO,
		produto.DES2_COMPLE,
		produto.ST2,
		produto.ST_PROMO1,
		produto.PC_PROMO,
		produto.PROM1_HI,
		produto.PROM1_HF,
		produto.PROM1_SEMANA,
		produto.PROMO_CONSUMO,
		produto.PROMO_PAGAR,
		produto.CRIADOEM,
		idHostLocal,
	)
	return erro
}

// CriarProdutoSoLocal insere uma nova entrada de produto no banco de dados local

// Criar insere um produto no banco de dados
func (repositorio Produtos) CriarProdutosSoLocal(produto modelos.Produto) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		` INSERT INTO PRODUTO (CODM, DES2, PV, IMPRE, FORNE, PS_L, PS_BR, BARRA, STATUS, UND, GRUPO, CNAE, CD_SERV, ` +
			`CODPISCOFINS, SERIAL, NCM, EAN, CSON, ORIGEN, ALICOTA, DES2_COMPLE, ST_PROMO1, PC_PROMO, PROM1_HI, PROM1_HF, PROM1_SEMANA, ` +
			`PROMO_CONSUMO, PROMO_PAGAR) ` +
			`VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	// Insira os valores para cada coluna
	resultado, erro := statement.Exec(
		produto.CODM, produto.DES2, produto.PV, produto.IMPRE, produto.FORNE, produto.PS_L, produto.PS_BR,
		produto.BARRA, produto.STATUS, produto.UND, produto.GRUPO, produto.CNAE, produto.CD_SERV,
		produto.CODPISCOFINS, produto.SERIAL, produto.NCM, produto.EAN, produto.CSON, produto.ORIGEN, produto.ALICOTA,
		produto.DES2_COMPLE, produto.ST_PROMO1, produto.PC_PROMO, produto.PROM1_HI, produto.PROM1_HF, produto.PROM1_SEMANA,
		produto.PROMO_CONSUMO, produto.PROMO_PAGAR)

	if erro != nil {
		fmt.Println(erro)
		return 0, erro
	}

	// Obtém o último ID inserido
	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// BuscarProdutosAguardandoSync busca produtos com AGUARDANDO_SYNC = "S"
func (repositorio Produtos) BuscarProdutosAguardandoSync() ([]modelos.Produto, error) {
	rows, _ := repositorio.db.Query(`
		SELECT ID, ID_HOSTWEB, CODM, DES1, DES2, PV, PV2, IMPRE, FORNE, CST1, CST2, STOCK1, STOCK2, STOCK3, STOCK4, 
		       COMPOSI, MINIMO, PS_L, PS_BR, COM_1, COM_2, BARRA, OB1, OB2, OB3, OB4, OB5, OBS1, OBS2, 
		       OBS3, OBS4, OBS5, STATUS, UND, CARDAPIO, GRUPO, FISCAL, FISCAL1, BAIXAD, SALDO_AN, CUSTO_CM, 
		       CUSTO_MD_AN, CUSTO_MD_AT, CNAE, CD_SERV, COMISSAO, CODPISCOFINS, SERIAL, QPCX, PV3, NCM, 
		       EAN, CSON, ORIGEN, ALICOTA, SERVICO, DES2_COMPLE, ST2, ST_PROMO1, PC_PROMO, PROM1_HI, 
		       PROM1_HF, PROM1_SEMANA, PROMO_CONSUMO, PROMO_PAGAR, CRIADOEM
		FROM PRODUTO WHERE AGUARDANDO_SYNC = 'S'
	`)
	defer rows.Close()

	var produtos []modelos.Produto
	for rows.Next() {
		var produto modelos.Produto
		rows.Scan(
			&produto.ID,
			&produto.ID_HOSTWEB,
			&produto.CODM,
			&produto.DES1,
			&produto.DES2,
			&produto.PV,
			&produto.PV2,
			&produto.IMPRE,
			&produto.FORNE,
			&produto.CST1,
			&produto.CST2,
			&produto.STOCK1,
			&produto.STOCK2,
			&produto.STOCK3,
			&produto.STOCK4,
			&produto.COMPOSI,
			&produto.MINIMO,
			&produto.PS_L,
			&produto.PS_BR,
			&produto.COM_1,
			&produto.COM_2,
			&produto.BARRA,
			&produto.OB1,
			&produto.OB2,
			&produto.OB3,
			&produto.OB4,
			&produto.OB5,
			&produto.OBS1,
			&produto.OBS2,
			&produto.OBS3,
			&produto.OBS4,
			&produto.OBS5,
			&produto.STATUS,
			&produto.UND,
			&produto.CARDAPIO,
			&produto.GRUPO,
			&produto.FISCAL,
			&produto.FISCAL1,
			&produto.BAIXAD,
			&produto.SALDO_AN,
			&produto.CUSTO_CM,
			&produto.CUSTO_MD_AN,
			&produto.CUSTO_MD_AT,
			&produto.CNAE,
			&produto.CD_SERV,
			&produto.COMISSAO,
			&produto.CODPISCOFINS,
			&produto.SERIAL,
			&produto.QPCX,
			&produto.PV3,
			&produto.NCM,
			&produto.EAN,
			&produto.CSON,
			&produto.ORIGEN,
			&produto.ALICOTA,
			&produto.SERVICO,
			&produto.DES2_COMPLE,
			&produto.ST2,
			&produto.ST_PROMO1,
			&produto.PC_PROMO,
			&produto.PROM1_HI,
			&produto.PROM1_HF,
			&produto.PROM1_SEMANA,
			&produto.PROMO_CONSUMO,
			&produto.PROMO_PAGAR,
			&produto.CRIADOEM,
		)
		produtos = append(produtos, produto)
	}
	return produtos, nil
}

// DesmarcarAguardandoSyncLocal desmarca o campo AGUARDANDO_SYNC no produto local
func (repositorio Produtos) DesmarcarAguardandoSyncLocal(produtoID uint64) error {
	_, erro := repositorio.db.Exec("UPDATE PRODUTO SET AGUARDANDO_SYNC = '' WHERE ID = ?", produtoID)
	return erro
}

// BuscarPreferidos busca produtos com STATUS = 'C' e OB4 = 'S'
func (repositorio Produtos) BuscarPreferidos() ([]modelos.Produto, error) {
	query := `SELECT ID, CODM, DES2, PV, IMPRE, FORNE, PS_L, PS_BR, BARRA, STATUS, UND, GRUPO, CNAE, 
			  CD_SERV, CODPISCOFINS, SERIAL, NCM, EAN, CSON, ORIGEN, ALICOTA, DES2_COMPLE, ST_PROMO1, 
			  PC_PROMO, PROM1_HI, PROM1_HF, PROM1_SEMANA, PROMO_CONSUMO, PROMO_PAGAR, ID_HOSTLOCAL 
			  FROM PRODUTO WHERE STATUS = 'C' AND OB4 = 'S'`

	rows, erro := repositorio.db.Query(query)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var produtos []modelos.Produto
	for rows.Next() {
		var produto modelos.Produto
		if erro = rows.Scan(
			&produto.ID,
			&produto.CODM,
			&produto.DES2,
			&produto.PV,
			&produto.IMPRE,
			&produto.FORNE,
			&produto.PS_L,
			&produto.PS_BR,
			&produto.BARRA,
			&produto.STATUS,
			&produto.UND,
			&produto.GRUPO,
			&produto.CNAE,
			&produto.CD_SERV,
			&produto.CODPISCOFINS,
			&produto.SERIAL,
			&produto.NCM,
			&produto.EAN,
			&produto.CSON,
			&produto.ORIGEN,
			&produto.ALICOTA,
			&produto.DES2_COMPLE,
			&produto.ST_PROMO1,
			&produto.PC_PROMO,
			&produto.PROM1_HI,
			&produto.PROM1_HF,
			&produto.PROM1_SEMANA,
			&produto.PROMO_CONSUMO,
			&produto.PROMO_PAGAR,
			&produto.ID_HOSTLOCAL,
		); erro != nil {
			return nil, erro
		}
		produtos = append(produtos, produto)
	}

	return produtos, nil
}

// BuscarProdutosPorGrupo retorna produtos com base no ID do grupo
func (repositorio Produtos) BuscarProdutosPorGrupo(grupoID uint64) ([]modelos.Produto, error) {
	query := `SELECT ID, CODM, DES2, PV, STATUS, GRUPO 
              FROM PRODUTO 
              WHERE GRUPO = ?`

	rows, erro := repositorio.db.Query(query, grupoID)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var produtos []modelos.Produto
	for rows.Next() {
		var produto modelos.Produto
		if erro = rows.Scan(
			&produto.ID,
			&produto.CODM,
			&produto.DES2,
			&produto.PV,
			&produto.STATUS,
			&produto.GRUPO,
		); erro != nil {
			return nil, erro
		}
		produtos = append(produtos, produto)
	}

	return produtos, nil
}

// BuscarProdutoPorCODM busca um produto no banco de dados pelo CODM
func (repositorio Produtos) BuscarProdutoPorCODM(codm string) (modelos.Produto, error) {

	row := repositorio.db.QueryRow(`
		SELECT ID, CODM, DES1, DES2, PV, PV2, IMPRE, FORNE, CST1, CST2, STOCK1, STOCK2, STOCK3, STOCK4, 
		       COMPOSI, MINIMO, PS_L, PS_BR, COM_1, COM_2, BARRA, OB1, OB2, OB3, OB4, OB5, OBS1, OBS2, OBS3, 
		       OBS4, OBS5, STATUS, UND, CARDAPIO, GRUPO, FISCAL, FISCAL1, BAIXAD, SALDO_AN, CUSTO_CM, 
		       CUSTO_MD_AN, CUSTO_MD_AT, CNAE, CD_SERV, COMISSAO, CODPISCOFINS, SERIAL, QPCX, PV3, NCM, EAN, 
		       CSON, ORIGEN, ALICOTA, SERVICO, DES2_COMPLE, ST2, ST_PROMO1, PC_PROMO, PROM1_HI, PROM1_HF, 
		       PROM1_SEMANA, PROMO_CONSUMO, PROMO_PAGAR 
		FROM PRODUTO 
		WHERE TRIM(CODM) = ?
	`, codm)

	var produto modelos.Produto
	if erro := row.Scan(
		&produto.ID,
		&produto.CODM,
		&produto.DES1,
		&produto.DES2,
		&produto.PV,
		&produto.PV2,
		&produto.IMPRE,
		&produto.FORNE,
		&produto.CST1,
		&produto.CST2,
		&produto.STOCK1,
		&produto.STOCK2,
		&produto.STOCK3,
		&produto.STOCK4,
		&produto.COMPOSI,
		&produto.MINIMO,
		&produto.PS_L,
		&produto.PS_BR,
		&produto.COM_1,
		&produto.COM_2,
		&produto.BARRA,
		&produto.OB1,
		&produto.OB2,
		&produto.OB3,
		&produto.OB4,
		&produto.OB5,
		&produto.OBS1,
		&produto.OBS2,
		&produto.OBS3,
		&produto.OBS4,
		&produto.OBS5,
		&produto.STATUS,
		&produto.UND,
		&produto.CARDAPIO,
		&produto.GRUPO,
		&produto.FISCAL,
		&produto.FISCAL1,
		&produto.BAIXAD,
		&produto.SALDO_AN,
		&produto.CUSTO_CM,
		&produto.CUSTO_MD_AN,
		&produto.CUSTO_MD_AT,
		&produto.CNAE,
		&produto.CD_SERV,
		&produto.COMISSAO,
		&produto.CODPISCOFINS,
		&produto.SERIAL,
		&produto.QPCX,
		&produto.PV3,
		&produto.NCM,
		&produto.EAN,
		&produto.CSON,
		&produto.ORIGEN,
		&produto.ALICOTA,
		&produto.SERVICO,
		&produto.DES2_COMPLE,
		&produto.ST2,
		&produto.ST_PROMO1,
		&produto.PC_PROMO,
		&produto.PROM1_HI,
		&produto.PROM1_HF,
		&produto.PROM1_SEMANA,
		&produto.PROMO_CONSUMO,
		&produto.PROMO_PAGAR,
	); erro != nil {
		return modelos.Produto{}, erro
	}

	return produto, nil
}
