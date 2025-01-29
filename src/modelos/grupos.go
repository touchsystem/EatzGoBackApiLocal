package modelos

// Grupo representa o modelo para a tabela  GRUPOS
type Grupo struct {
	ID              uint64 `json:"id,omitempty"`           // Identificador único
	ID_HOSTWEB      int    `json:"id_hostweb,omitempty"`   // ID local
	ID_HOSTLOCAL    int    `json:"id_hostlocal,omitempty"` // ID local
	AGUARDANDO_SYNC string `json:"aguardando_sync,omitempty"`
	COD_GP          string `json:"cod_gp,omitempty"`         // Código do grupo
	NOME            string `json:"nome,omitempty"`           // Nome do grupo
	TIPO            string `json:"tipo,omitempty"`           // Tipo do grupo
	CONTA_CONTABIL  string `json:"conta_contabil,omitempty"` // Conta contábil associada
}
