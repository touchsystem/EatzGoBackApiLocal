package modelos

import (
	"encoding/json"
	"errors"
	"time"
)

type CxRecebTipo struct {
	ID        int     `json:"id"`
	IDCxReceb int     `json:"id_cx_receb"`
	IDHostWeb int     `json:"id_hostweb"`
	IDTipoRec int     `json:"id_tipo_rec"`
	MoedaNac  float64 `json:"moeda_nac"`
	MoedaExt  float64 `json:"moeda_ext"`
}

type CxReceb struct {
	ID        int       `json:"id"`
	IDHostWeb int       `json:"id_hostweb"`
	Status    string    `json:"status"`
	Data      time.Time `json:"data"`
	IDCli     int       `json:"id_cli"`
	IDUser    int       `json:"id_user"`
	Total     float64   `json:"total"`
	Troco     float64   `json:"troco"`
	Mesa      int       `json:"mesa"`
	NRPessoas int       `json:"nr_pessoas"`
	CriadoEm  time.Time `json:"criado_em"`
}

// Método para Unmarshal JSON com formato de data personalizado
func (c *CxReceb) UnmarshalJSON(data []byte) error {
	type Alias CxReceb
	aux := &struct {
		Data string `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse manual do campo "data"
	if aux.Data != "" {
		parsedTime, err := time.Parse("2006-01-02", aux.Data)
		if err != nil {
			return errors.New("formato de data inválido")
		}
		c.Data = parsedTime
	}

	return nil
}
