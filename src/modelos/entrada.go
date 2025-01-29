package modelos

type Entrada struct {
	IDEntra  int     `json:"id_entra,omitempty"`
	IDFor    int     `json:"id_for"`
	Data     string  `json:"data"`
	VLTO     float64 `json:"vlto"`
	NAT      string  `json:"nat"`
	NRNota   int     `json:"nr_nota"`
	DTCan    string  `json:"dt_can,omitempty"`
	MTCan    string  `json:"mt_can,omitempty"`
	IDUsu    int     `json:"id_usu"`
	Status   string  `json:"status,omitempty"`
	Emissao  string  `json:"emissao"`
	Serie    string  `json:"serie,omitempty"`
	Subserie string  `json:"subserie,omitempty"`
	Modelo   string  `json:"modelo,omitempty"`
	Tipo     string  `json:"tipo,omitempty"`
	CriadoEm string  `json:"criado_em,omitempty"`
}

// EntradaComProdutos encapsula uma entrada e uma lista de produtos

type EntradaComProdutos struct {
	Entrada  Entrada      `json:"entrada"`
	Produtos []ProdutoEnt `json:"produtos"`
	Titulos  []Pagar      `json:"titulos"`
}

// ProdutoEnt representa um produto relacionado a uma entrada
type ProdutoEnt struct {
	IDEntra int     `json:"id_entra,omitempty"`
	IDFor   int     `json:"id_for"`
	IDUsu   int     `json:"id_usu"`
	Status  string  `json:"status"`
	CODM    string  `json:"codm"`
	QTDE    float64 `json:"qtde"`
	TOTAL   float64 `json:"total"`
}
