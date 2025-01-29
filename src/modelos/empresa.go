package modelos

// Empresa representa o modelo da tabela EMPRESA no banco de dados
type Empresa struct {
	ID_EMPRESA   int    `json:"id_empresa,omitempty"`
	DATA_SISTEMA string `json:"data_sistema,omitempty"`
}
