package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

// ContaRepository representa um repositório para operações da conta
type ContaRepository struct {
	db *sql.DB
}

// NovoRepositorioDeContas cria um novo repositório de contas
func NovoRepositorioDeContas(db *sql.DB) *ContaRepository {
	return &ContaRepository{db}
}

// BuscarVendasPorMesa busca as vendas ativas por número da mesa com informações adicionais
func (repositorio ContaRepository) BuscarVendasPorMesa(mesa int) ([]modelos.ContaVenda, error) {
	query := `
		SELECT 
			V.CODM, P.DES2, V.CELULAR, V.CPF_CLIENTE, V.NOME_CLIENTE, 
			V.ID_CLIENTE, V.PV, V.PV_PROM, V.NICK, V.DATA
		FROM VENDA V
		LEFT JOIN PRODUTO P ON V.CODM = P.CODM
		WHERE V.MESA = ? AND V.STATUS = 'A'`

	rows, err := repositorio.db.Query(query, mesa)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendas []modelos.ContaVenda
	for rows.Next() {
		var venda modelos.ContaVenda
		err = rows.Scan(
			&venda.CODM, &venda.DES2, &venda.Celular,
			&venda.CPFCliente, &venda.NomeCliente, &venda.IDCliente,
			&venda.PV, &venda.PVProm, &venda.Nick, &venda.Data,
		)
		if err != nil {
			return nil, err
		}
		vendas = append(vendas, venda)
	}

	return vendas, nil
}

// BuscarMesaPorNumero busca os detalhes da mesa pelo número
func (repositorio ContaRepository) BuscarMesaPorNumero(mesa int) (modelos.Mesa, error) {
	//	fmt.Printf("Buscando mesa: %d\n", mesa)
	query := `SELECT ID, MESA_CARTAO, STATUS, ID_USER, ID_CLI, NICK, ABERTURA, QTD_PESSOAS, TURISTA, CELULAR, APELIDO 
          FROM MESA WHERE MESA_CARTAO = ?`

	//	fmt.Printf("Executando query: %s com parâmetro: %d\n", query, mesa)

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

	//	fmt.Printf("Mesa encontrada: %+v\n", mesaInfo)
	return mesaInfo, nil
}
