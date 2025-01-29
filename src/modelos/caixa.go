package modelos

import "time"

// Caixa representa o modelo da tabela CAIXA no banco de dados
type Caixa struct {
	ID         int       `json:"id"`
	NR_CAIXA   string    `json:"conta"`
	DOCUMENTO  string    `json:"documento"`
	DATA       time.Time `json:"data"`
	VALOR      float64   `json:"valor_cre"`
	VALORS     float64   `json:"valor_deb"`
	COMPLEMENT string    `json:"complemento"`
	CriadoEm   time.Time `json:"criadoEm"`
	IDUsu      int       `json:"id_usu"`
}

// LancamentoCaixaDuplo representa a estrutura do lançamento duplo
type LancamentoCaixaDuplo struct {
	ContaCredito string  `json:"conta_credito"`
	ContaDebito  string  `json:"conta_debito"`
	Documento    string  `json:"documento"`
	Valor        float64 `json:"valor"`
	Data         string  `json:"data"`
	Complement   string  `json:"complement"`
	IDUsu        int     `json:"id_usu"`
	// Verifique se este campo está sendo enviado corretamente no JSON
}
