package modelos

import (
	"database/sql"
	"errors"
)

// Mesa representa uma mesa no sistema
type Mesa struct {
	ID         uint64 `json:"id,omitempty"`
	MesaCartao int    `json:"mesa_cartao,omitempty"`
	Status     string `json:"status,omitempty"`
	IDUser     int    `json:"id_user,omitempty"`
	IDCli      int    `json:"id_cli,omitempty"`
	Nick       string `json:"nick,omitempty"`

	Abertura   sql.NullTime `json:"abertura,omitempty"`
	QtdPessoas int          `json:"qtd_pessoas,omitempty"`
	Turista    string       `json:"turista,omitempty"`
	Celular    string       `json:"celular,omitempty"`
	Apelido    string       `json:"apelido,omitempty"`
}

// Validar valida os campos da estrutura Mesa
func (m *Mesa) Validar() error {
	// Validar MesaCartao
	if m.MesaCartao <= 0 {
		return errors.New("o campo MesaCartao é obrigatório e deve ser maior que zero")
	}

	// Validar Status
	statusPermitidos := map[string]bool{
		"L": true, // Livre
		"O": true, // Ocupado
		"F": true, // Fechamento
		"I": true, // Inativa
	}

	if _, permitido := statusPermitidos[m.Status]; !permitido {
		return errors.New("o campo Status possui um valor inválido. Valores permitidos: 'L'ivre, 'O'cupado, 'F'echamento, 'I'nativa ")
	}

	return nil
}
