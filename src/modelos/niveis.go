package modelos

import (
	"errors"
	"strings"
)

// NivelAcesso representa um nível de acesso no sistema
type NivelAcesso struct {
	ID         uint64 `json:"id,omitempty"`
	Codigo     string `json:"codigo,omitempty"`
	NomeAcesso string `json:"nome_acesso,omitempty"`
	Nivel      int    `json:"nivel,omitempty"`
}

// Preparar formata e valida o nível de acesso
func (nivel *NivelAcesso) Preparar() error {
	nivel.Codigo = strings.TrimSpace(nivel.Codigo)
	nivel.NomeAcesso = strings.TrimSpace(nivel.NomeAcesso)

	if nivel.NomeAcesso == "" {
		return errors.New("O nome do acesso é obrigatório")
	}

	return nil
}
