package modelos

import (
	"errors"

	"github.com/badoux/checkmail"
)

// Fornecedor representa um fornecedor no sistema
type Fornecedor struct {
	ID       uint64 `json:"id,omitempty"`
	CNPJ_CPF string `json:"cnpj_cpf,omitempty"`
	NOME     string `json:"nome,omitempty"`
	FANTASIA string `json:"fantasia,omitempty"`
	ENDERE   string `json:"endere,omitempty"`
	CIDADE   string `json:"cidade,omitempty"`
	BAIRRO   string `json:"bairro,omitempty"`
	CEP      string `json:"cep,omitempty"`
	UF       string `json:"uf,omitempty"`
	TELE1    string `json:"tele1,omitempty"`
	CEL1     string `json:"cel1,omitempty"`
	CONTATO  string `json:"contato,omitempty"`
	EMAIL    string `json:"email,omitempty"`
	GRUPO    int    `json:"grupo,omitempty"`
	PLACA    string `json:"placa,omitempty"`
	NUMERO   string `json:"numero,omitempty"`
	CIDIBGE  string `json:"cidibge,omitempty"`
	COMPLE   string `json:"comple,omitempty"`
}

// Preparar vai chamar os métodos para validar e formatar o fornecedor recebido
func (fornecedor *Fornecedor) Preparar(etapa string) error {
	if erro := fornecedor.validar(etapa); erro != nil {
		return erro
	}
	return nil
}

func (fornecedor *Fornecedor) validar(etapa string) error {
	if fornecedor.NOME == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}

	if erro := checkmail.ValidateFormat(fornecedor.EMAIL); erro != nil {
		return errors.New("O e-mail inserido é inválido")
	}
	if fornecedor.CEL1 == "" {
		return errors.New("O numero do celular é obrigatório e não pode estar em branco")
	}

	// Adicione outras validações necessárias

	return nil
}
