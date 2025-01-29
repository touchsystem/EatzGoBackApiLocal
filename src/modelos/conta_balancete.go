package modelos

import (
	"errors"
	"strings"
)

// ContaBalancete representa uma conta balancete no sistema
type ContaBalancete struct {
	ID        uint64  `json:"id,omitempty"`
	AG        string  `json:"ag,omitempty"`
	Conta     string  `json:"conta,omitempty"`
	Descricao string  `json:"descricao,omitempty"`
	Valor     float64 `json:"valor_cre,omitempty"`
	Valors    float64 `json:"valor_deb,omitempty"`
	Valora    float64 `json:"saldo_ant,omitempty"`
}

// ContaBalancete representa uma conta balancete no sistema
type ResultadoBalancete struct {
	ID        uint64  `json:"id,omitempty"`
	AG        string  `json:"ag,omitempty"`
	Conta     string  `json:"conta,omitempty"`
	Descricao string  `json:"descricao,omitempty"`
	Valor_cre float64 `json:"valor_cre,omitempty"`
	Valor_deb float64 `json:"valor_deb,omitempty"`
	Saldo_ant float64 `json:"saldo_ant,omitempty"`
	Saldo_atu float64 `json:"saldo_atu,omitempty"`
}

// Preparar formata e valida a conta balancete
func (contaBalancete *ContaBalancete) Preparar() error {
	contaBalancete.AG = strings.TrimSpace(contaBalancete.AG)
	contaBalancete.Conta = strings.TrimSpace(contaBalancete.Conta)
	contaBalancete.Descricao = strings.TrimSpace(contaBalancete.Descricao)

	if contaBalancete.Conta == "" {
		return errors.New("A conta é obrigatória")
	}

	if contaBalancete.Descricao == "" {
		return errors.New("A descrição é obrigatória")
	}

	return nil
}
