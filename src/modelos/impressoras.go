package modelos

// Impressora representa o modelo para a tabela IMPRESSORAS
type Impressora struct {
	ID              int    `json:"id,omitempty"`              // Identificador único
	ID_HOSTLOCAL    int    `json:"id_hostlocal,omitempty"`    // ID local
	AGUARDANDO_SYNC string `json:"aguardando_sync,omitempty"` // Indicador de sincronização
	COD_IMP         string `json:"cod_imp,omitempty"`         // Código da impressora
	NOME            string `json:"nome,omitempty"`            // Nome da impressora
	END_IMP         string `json:"end_imp,omitempty"`         // Nome da impressora
	END_SER         string `json:"end_ser,omitempty"`         // Endereço ou localização da impressora
}
