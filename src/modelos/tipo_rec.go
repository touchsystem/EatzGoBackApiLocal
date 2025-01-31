package modelos

import (
	"fmt"
)

type TipoRecebimento struct {
	ID     uint64  `json:"id,omitempty"`
	Nome   string  `json:"nome,omitempty"`
	Cambio float64 `json:"cambio,omitempty"`
	FtConv string  `json:"ft_conv,omitempty"` // Pode ser '*', '/' ou 'N'
	Status string  `json:"status,omitempty"`  // Pode ser 'A' (Ativo) ou 'I' (Inativo)
}

// Valida se o TipoRecebimento possui valores corretos
func (t *TipoRecebimento) Validar() error {
	if t.FtConv != "*" && t.FtConv != "/" && t.FtConv != "N" {
		return fmt.Errorf("FtConv inválido: deve ser '*', '/' ou 'N'")
	}

	if t.Status != "A" && t.Status != "I" {
		return fmt.Errorf("Status inválido: deve ser 'A' (Ativo) ou 'I' (Inativo)")
	}

	return nil
}
