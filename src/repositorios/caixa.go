package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

type Caixa struct {
	db *sql.DB
}

func NovoRepositorioDeCaixa(db *sql.DB) *Caixa {
	return &Caixa{db}
}

func (repositorio Caixa) BuscarVendasMesa(mesa int) ([]modelos.VendaRecebimento, error) {
	fmt.Println(mesa)
	query := `
SELECT TRIM(v.CODM), COALESCE(TRIM(p.DES2), 'Produto Desconhecido') AS PRODUTO, 
       v.PV, v.QTD, v.NICK, v.DESCONTO, v.ID_USER, v.CRIADOEM
FROM VENDA v
LEFT JOIN PRODUTO p ON TRIM(v.CODM) = TRIM(p.CODM)
WHERE v.MESA = ? AND v.STATUS = 'A';


	`

	rows, err := repositorio.db.Query(query, mesa)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendas []modelos.VendaRecebimento
	for rows.Next() {
		var venda modelos.VendaRecebimento
		if err := rows.Scan(&venda.CODM, &venda.Produto, &venda.PV, &venda.QTD, &venda.Nick, &venda.Desconto, &venda.IDUser, &venda.CriadoEm); err != nil {
			return nil, err
		}
		vendas = append(vendas, venda)
	}

	return vendas, nil
}
