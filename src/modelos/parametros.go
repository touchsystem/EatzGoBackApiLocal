package modelos

import "errors"

// Parametro representa um parâmetro no sistema
type Parametro struct {
	ID     uint64  `json:"id,omitempty"`
	Nome   string  `json:"nome,omitempty"`
	Status string  `json:"status,omitempty"`
	Limite float64 `json:"limite,omitempty"`
}

// Validar valida os campos do modelo Parametro
func (p *Parametro) Validar() error {
	// Verificar se o campo Nome está preenchido
	//	if p.Nome == "" {
	//		return errors.New("o campo Nome é obrigatório")
	//	}

	// Verificar se o Status contém valores válidos
	statusPermitidos := map[string]bool{
		"R": true, // Moeda Principal (Real)
		"G": true, // Guarani
		"D": true, // Dólar
		"P": true, // Peso AR
		"S": true, // Sim
		"N": true, // Não
	}

	if _, permitido := statusPermitidos[p.Status]; !permitido {
		return errors.New("o campo Status possui um valor inválido. Valores permitidos: 'R', 'G', 'D', 'P', 'S', 'N'")
	}

	return nil
}
