package modelos

type VendaRecebimento struct {
	CODM     string  `json:"codm"`
	Produto  string  `json:"produto"`
	PV       float64 `json:"pv"`
	QTD      float64 `json:"qtd"`
	Nick     string  `json:"nick"`
	Desconto float64 `json:"desconto"`
	IDUser   int     `json:"id_user"`
	CriadoEm string  `json:"criadoem"`
}
