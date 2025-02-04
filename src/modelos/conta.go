package modelos

// Conta representa uma conta associada a uma mesa
type Conta struct {
	MesaNumero     int          `json:"mesa_numero,omitempty"`
	NickUsuario    string       `json:"nick_usuario,omitempty"`
	NomeCliente    string       `json:"nome_cliente,omitempty"`
	IDCliente      *int         `json:"id_cliente,omitempty"`
	Celular        string       `json:"celular,omitempty"`
	DataHoraPedido string       `json:"data_hora_pedido,omitempty"`
	QtdPessoas     int          `json:"qtd_pessoas,omitempty"`
	HoraAbertura   string       `json:"hora_abertura,omitempty"`
	Vendas         []ContaVenda `json:"vendas,omitempty"` // Lista de vendas
	TotalConta     float64      `json:"total_conta"`      // Novo campo para total
}

// ContaVenda representa os itens vendidos na conta
type ContaVenda struct {
	CODM        string  `json:"codm,omitempty"`
	Descricao   string  `json:"descricao,omitempty"`
	DES2        string  `json:"des2,omitempty"`
	Celular     string  `json:"celular,omitempty"`
	CPFCliente  string  `json:"cpf_cliente,omitempty"`
	NomeCliente string  `json:"nome_cliente,omitempty"`
	IDCliente   *int    `json:"id_cliente,omitempty"`
	PV          float64 `json:"pv,omitempty"`
	PVProm      float64 `json:"pv_prom,omitempty"`
	QTD         float64 `json:"qtd,omitempty"`
	Nick        string  `json:"nick,omitempty"`
	Data        string  `json:"data,omitempty"`
}
