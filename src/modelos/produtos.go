package modelos

import (
	"errors"
)

// Produto representa um produto no sistema
type Produto struct {
	ID              uint64  `json:"id,omitempty"`
	ID_HOSTLOCAL    int     `json:"id_hostlocal,omitempty"`
	ID_HOSTWEB      int     `json:"id_hostweb,omitempty"`
	AGUARDANDO_SYNC string  `json:"AGUARDANDO_SYNC,omitempty"`
	CODM            string  `json:"codm,omitempty"`
	DES1            string  `json:"des1,omitempty"`
	DES2            string  `json:"des2,omitempty"`
	PV              float64 `json:"pv,omitempty"`
	PV2             float64 `json:"pv2,omitempty"`
	IMPRE           string  `json:"impre,omitempty"`
	FORNE           float64 `json:"forne,omitempty"`
	CST1            float64 `json:"cst1,omitempty"`
	CST2            float64 `json:"cst2,omitempty"`
	STOCK1          float64 `json:"stock1,omitempty"`
	STOCK2          float64 `json:"stock2,omitempty"`
	STOCK3          float64 `json:"stock3,omitempty"`
	STOCK4          float64 `json:"stock4,omitempty"`
	COMPOSI         string  `json:"composi,omitempty"`
	MINIMO          float64 `json:"minimo,omitempty"`
	PS_L            float64 `json:"ps_l,omitempty"`
	PS_BR           float64 `json:"ps_br,omitempty"`
	COM_1           float64 `json:"com_1,omitempty"`
	COM_2           float64 `json:"com_2,omitempty"`
	BARRA           string  `json:"barra,omitempty"`
	OB1             string  `json:"ob1,omitempty"`
	OB2             string  `json:"ob2,omitempty"`
	OB3             string  `json:"ob3,omitempty"`
	OB4             string  `json:"ob4,omitempty"`
	OB5             string  `json:"ob5,omitempty"`
	OBS1            string  `json:"obs1,omitempty"`
	OBS2            string  `json:"obs2,omitempty"`
	OBS3            string  `json:"obs3,omitempty"`
	OBS4            string  `json:"obs4,omitempty"`
	OBS5            string  `json:"obs5,omitempty"`
	STATUS          string  `json:"status,omitempty"`
	UND             string  `json:"und,omitempty"`
	CARDAPIO        string  `json:"cardapio,omitempty"`
	GRUPO           string  `json:"grupo,omitempty"`
	FISCAL          string  `json:"fiscal,omitempty"`
	FISCAL1         string  `json:"fiscal1,omitempty"`
	BAIXAD          float64 `json:"baixad,omitempty"`
	SALDO_AN        float64 `json:"saldo_an,omitempty"`
	CUSTO_CM        float64 `json:"custo_cm,omitempty"`
	CUSTO_MD_AN     float64 `json:"custo_md_an,omitempty"`
	CUSTO_MD_AT     float64 `json:"custo_md_at,omitempty"`
	CNAE            string  `json:"cnae,omitempty"`
	CD_SERV         string  `json:"cd_serv,omitempty"`
	COMISSAO        float64 `json:"comissao,omitempty"`
	CODPISCOFINS    string  `json:"codpiscofins,omitempty"`
	SERIAL          string  `json:"serial,omitempty"`
	QPCX            int     `json:"qpcx,omitempty"`
	PV3             float64 `json:"pv3,omitempty"`
	NCM             string  `json:"ncm,omitempty"`
	EAN             string  `json:"ean,omitempty"`
	CSON            string  `json:"cson,omitempty"`
	ORIGEN          string  `json:"origen,omitempty"`
	ALICOTA         float64 `json:"alicota,omitempty"`
	SERVICO         string  `json:"servico,omitempty"`
	DES2_COMPLE     string  `json:"des2_comple,omitempty"`
	ST2             string  `json:"st2,omitempty"`
	ST_PROMO1       string  `json:"st_promo1,omitempty"`
	PC_PROMO        float64 `json:"pc_promo,omitempty"`
	PROM1_HI        string  `json:"prom1_hi,omitempty"` // Assume as string for JSON handling
	PROM1_HF        string  `json:"prom1_hf,omitempty"` // Assume as string for JSON handling
	PROM1_SEMANA    string  `json:"prom1_semana,omitempty"`
	PROMO_CONSUMO   float64 `json:"promo_consumo,omitempty"`
	PROMO_PAGAR     float64 `json:"promo_pagar,omitempty"`
	CRIADOEM        string  `json:"criadoem,omitempty"` // Assume as string for JSON handling
}

// Preparar vai chamar os métodos para validar e formatar o produto recebido
func (produto *Produto) Preparar(etapa string) error {
	return produto.Validar(etapa)
}

func (produto *Produto) Validar(etapa string) error {
	if (produto.STATUS != "C") && (produto.STATUS != "P") && (produto.STATUS != "I") && (produto.STATUS != "D") {
		//	fmt.Println(produto.STATUS + "P")
		return errors.New("Status Inválido " + produto.STATUS)

	}

	if produto.CODM == "" {
		return errors.New("O código do produto é obrigatório e não pode estar em branco")
	}

	if produto.DES2 == "" {
		return errors.New("A descrição do produto é obrigatória e não pode estar em branco")
	}

	// Adicione outras validações necessárias

	return nil
}
