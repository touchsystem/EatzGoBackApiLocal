package modelos

// Receber representa o modelo da tabela RECEBER no banco de dados
type Receber struct {
	ID_RECEBER    int     `json:"id_receber,omitempty"`
	DOC_L         string  `json:"doc_l,omitempty"`
	DATA          string  `json:"data,omitempty"`
	DT_VEN        string  `json:"dt_ven,omitempty"`
	VL_TIT        float64 `json:"vl_tit,omitempty"`
	VL_SAL        float64 `json:"vl_sal,omitempty"`
	ID_CLI        int     `json:"id_cli,omitempty"`
	IDUsu         int     `json:"id_usu,omitempty"`
	DESCONTO      float64 `json:"desconto,omitempty"`
	ACRESCIMO     float64 `json:"acrescimo,omitempty"`
	OBS           string  `json:"obs,omitempty"`
	SITUA         string  `json:"situa,omitempty"`
	CNPJ          string  `json:"cnpj,omitempty"`
	ContraPartida string  `json:"contra_partida,omitempty"`
	DT_PGTO       string  `json:"dt_pgto,omitempty"` // Data de pagamento
	NOME_CLIENTE  string  `json:"nome_cliente,omitempty"`
}
