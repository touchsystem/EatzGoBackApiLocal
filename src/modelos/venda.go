package modelos

import (
	"errors"
)

type Venda struct {
	ID              int     `json:"id,omitempty"`
	CODM            string  `json:"codm,omitempty"`
	STATUS          string  `json:"status,omitempty"`
	IMPRESSORA      string  `json:"impressora,omitempty"`
	MONITOR_PRD     string  `json:"monitor_prd,omitempty"`
	STATUS_TP_VENDA string  `json:"status_tp_venda,omitempty"`
	STATUS_TP_HARD  string  `json:"status_tp_hard,omitempty"`
	MESA            int     `json:"mesa,omitempty"`
	CELULAR         string  `json:"celular,omitempty"`
	CPF_CLIENTE     string  `json:"cpf_cliente,omitempty"`
	NOME_CLIENTE    string  `json:"nome_cliente,omitempty"`
	ID_CLIENTE      int     `json:"id_cliente,omitempty"`
	PV              float64 `json:"pv,omitempty"`
	PV_PROM         float64 `json:"pv_prom,omitempty"`
	QTD             float64 `json:"qtd,omitempty"`
	ID_USER         int     `json:"id_user,omitempty"`
	NICK            string  `json:"nick,omitempty"`
	DATA            string  `json:"data,omitempty"`
	OBS             string  `json:"obs,omitempty"`
	OBS2            string  `json:"obs2,omitempty"`
	OBS3            string  `json:"obs3,omitempty"`
	STATUS_PGTO     string  `json:"status_pgto,omitempty"`
	DATA_IFOOD      *string `json:"data_ifood,omitempty"`
	STATUS_IFOOD    string  `json:"status_ifood,omitempty"`
	ID_IFOOD        string  `json:"id_ifood,omitempty"`
	COMPLEM_CODM    string  `json:"complem_codm,omitempty"`
	CHAVE           string  `json:"chave,omitempty"`
	CRIADOEM        string  `json:"criadoem,omitempty"`
}

func (venda *Venda) Preparar() error {
	// Validação de campos obrigatórios
	if venda.CODM == "" {
		return errors.New("O campo CODM é obrigatório")
	}
	//	if venda.PV <= 0 {
	//		return errors.New("O valor PV deve ser maior que zero")
	//	}

	// Validação de STATUS_TP_VENDA
	valoresPermitidos := map[string]bool{
		"P": true, // Pedido terminal
		"M": true, // Pedido dispositivo móvel
		"Y": true, // Delivery
		"D": true, // Venda direta
	}
	if !valoresPermitidos[venda.STATUS_TP_VENDA] {
		return errors.New("O campo STATUS_TP_VENDA aceita apenas os valores: P, M, Y, D")
	}

	return nil
}
