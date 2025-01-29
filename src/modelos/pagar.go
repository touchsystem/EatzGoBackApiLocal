package modelos

// Pagar representa o modelo da tabela PAGAR no banco de dados
type Pagar struct {
	ID_PAGAR        int     `json:"id_pagar,omitempty"`
	ID_HOSTLOCAL    int     `json:"id_hostlocal,omitempty"`
	ID_HOSTWEB      int     `json:"id_hostweb,omitempty"`
	STATUS          string  `json:"STATUS,omitempty"`
	DOC_L           string  `json:"doc_l,omitempty"`
	DATA            string  `json:"data,omitempty"`
	DT_VEN          string  `json:"dt_ven,omitempty"`
	VL_TIT          float64 `json:"vl_tit,omitempty"`
	VL_SAL          float64 `json:"vl_sal,omitempty"`
	ID_FOR          int     `json:"id_for,omitempty"`
	ID_ENTRA        int     `json:"id_entra,omitempty"`
	IDUsu           int     `json:"id_usu,omitempty"`
	DESCONTO        float64 `json:"desconto,omitempty"`
	ACRESCIMO       float64 `json:"acrescimo,omitempty"`
	OBS             string  `json:"obs,omitempty"`
	SITUA           string  `json:"situa,omitempty"`
	CNPJ            string  `json:"cnpj,omitempty"`
	ContraPartida   string  `json:"contra_partida,omitempty"`
	DT_PGTO         string  `json:"dt_pgto,omitempty"` // Data de pagamento
	NOME_FORNECEDOR string  `json:"nome_fornecedor,omitempty"`
}
