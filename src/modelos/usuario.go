package modelos

import (
	"api/src/seguranca"
	"errors"
	"strings"
	"time"
)

// Usuario representa um usuário utilizando a rede social
type Usuario struct {
	ID              uint64    `json:"id,omitempty"`
	ID_HOSTWEB      uint64    `json:"id_hostweb,omitempty"`
	ID_HOSTLOCAL    uint64    `json:"id_hostlocal,omitempty"`
	AGUARDANDO_SYNC string    `json:"aguardando_sync,omitempty"`
	Nome            string    `json:"nome,omitempty"`
	Nick            string    `json:"nick,omitempty"`
	Email           string    `json:"email,omitempty"`
	Senha           string    `json:"senha,omitempty"`
	CriadoEm        time.Time `json:"CriadoEm,omitempty"`
	CDEMP           string    `json:"CDEMP,omitempty"`
	Nivel           uint64    `json:"nivel,omitempty"`
}

// Preparar vai chamar os métodos para validar e formatar o usuário recebido
func (usuario *Usuario) Preparar(etapa string) error {
	if erro := usuario.validar(etapa); erro != nil {
		return erro
	}

	if erro := usuario.formatar(etapa); erro != nil {
		return erro
	}

	return nil
}

func (usuario *Usuario) validar(etapa string) error {
	if usuario.Nome == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}

	//if usuario.Nick == "" {
	//	return errors.New("O nick é obrigatório e não pode estar em branco")
	//}

	if usuario.Email == "" {
		return errors.New("O Login  é obrigatório e não pode estar em branco")
	}

	//if erro := checkmail.ValidateFormat(usuario.Email); erro != nil {
	//	return errors.New("O e-mail inserido é inválido")
	//}

	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("A senha é obrigatória e não pode estar em branco")
	}
	if etapa == "cadastro" && usuario.Nivel == 0 {
		return errors.New("o nivel do usuário é obrigatória e não pode estar em branco")
	}
	return nil
}

func (usuario *Usuario) formatar(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)
	usuario.CDEMP = strings.TrimSpace(usuario.CDEMP)

	if etapa == "cadastro" {
		senhaComHash, erro := seguranca.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}

		usuario.Senha = string(senhaComHash)
	}

	return nil
}
