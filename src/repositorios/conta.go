package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

// ContaRepository representa o repositório de contas
type ContaRepository struct {
	db *sql.DB
}

// NovoRepositorioDeContas cria um novo repositório de contas
func NovoRepositorioDeContas(db *sql.DB) *ContaRepository {
	return &ContaRepository{db}
}

// BuscarVendasPorMesa busca as vendas ativas por número da mesa com informações adicionais
func (repositorio ContaRepository) BuscarVendasPorMesa(mesa int) ([]modelos.ContaVenda, float64, error) {
	query := `
		SELECT 
			V.CODM, P.DES2, V.CELULAR, V.CPF_CLIENTE, V.NOME_CLIENTE, 
			V.ID_CLIENTE, V.PV, V.QTD, V.PV_PROM, V.NICK, V.DATA
		FROM VENDA V
		LEFT JOIN PRODUTO P ON V.CODM = P.CODM
		WHERE V.MESA = ? AND V.STATUS = 'A'`

	rows, err := repositorio.db.Query(query, mesa)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var vendas []modelos.ContaVenda
	var totalConta float64

	for rows.Next() {
		var venda modelos.ContaVenda
		err = rows.Scan(
			&venda.CODM, &venda.DES2, &venda.Celular,
			&venda.CPFCliente, &venda.NomeCliente, &venda.IDCliente,
			&venda.PV, &venda.QTD, &venda.PVProm, &venda.Nick, &venda.Data,
		)
		if err != nil {
			return nil, 0, err
		}

		// Cálculo do total: QTD * PV (Preço unitário)
		totalConta += venda.QTD * venda.PV

		vendas = append(vendas, venda)
	}

	return vendas, totalConta, nil
}

// BuscarMesaPorNumero busca os detalhes da mesa pelo número
func (repositorio ContaRepository) BuscarMesaPorNumero(mesa int) (modelos.Mesa, error) {
	query := `SELECT ID, MESA_CARTAO, STATUS, ID_USER, ID_CLI, NICK, ABERTURA, QTD_PESSOAS, TURISTA, CELULAR, APELIDO 
          FROM MESA WHERE MESA_CARTAO = ?`

	var mesaInfo modelos.Mesa
	err := repositorio.db.QueryRow(query, mesa).Scan(
		&mesaInfo.ID, &mesaInfo.MesaCartao, &mesaInfo.Status, &mesaInfo.IDUser,
		&mesaInfo.IDCli, &mesaInfo.Nick, &mesaInfo.Abertura, &mesaInfo.QtdPessoas,
		&mesaInfo.Turista, &mesaInfo.Celular, &mesaInfo.Apelido,
	)

	if err != nil {
		fmt.Printf("Erro ao buscar mesa: %v\n", err)
		return modelos.Mesa{}, err
	}

	return mesaInfo, nil
}

// BuscarParametroPorID busca um parâmetro pelo ID no banco de dados
func (repositorio ContaRepository) BuscarParametroPorID(parametroID int) (modelos.Parametro, error) {
	query := `SELECT ID, NOME, STATUS, LIMITE FROM PARAMETROS WHERE id = ?`

	var parametro modelos.Parametro
	err := repositorio.db.QueryRow(query, parametroID).Scan(
		&parametro.ID, &parametro.Nome, &parametro.Status, &parametro.Limite,
	)

	if err != nil {
		fmt.Printf("Erro ao buscar parâmetro: %v\n", err)
		return modelos.Parametro{}, err
	}

	return parametro, nil
}
