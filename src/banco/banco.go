package banco

import (
	"api/src/config"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // Driver
)

// Conectar abre a conexão com o banco de dados e a retorna
func Conectar(cdEmp string) (*sql.DB, error) {
	stringConexao := config.GetStringConexaoBanco(cdEmp)
	db, erro := sql.Open("mysql", stringConexao)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}

// ConectarPorEmpresa abre a conexão com o banco de dados da empresa específica
func ConectarPorEmpresa(cdEmp string) (*sql.DB, error) {
	// Buscar o nome do banco de dados da empresa no .env
	//nomeBanco := os.Getenv(fmt.Sprintf("CD_EMP_%s", cdEmp))
	nomeBanco := "databaseLOCAL"

	if nomeBanco == "" {
		return nil, fmt.Errorf("Banco de dados não encontrado para a empresa %s", cdEmp)
	}

	// Construir a string de conexão usando as variáveis de ambiente
	usuario := os.Getenv("DB_USUARIO")
	senha := os.Getenv("DB_SENHA")
	host := os.Getenv("localhost") // Adicione o DB_HOST no seu .env se necessário (padrão localhost)

	// String de conexão no formato: usuario:senha@tcp(host)/nomeBanco
	conexao := fmt.Sprintf("%s:%s@tcp(%s)/%s", usuario, senha, host, nomeBanco)

	// Abrir conexão com o banco de dados da empresa
	db, erro := sql.Open("mysql", conexao)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}
