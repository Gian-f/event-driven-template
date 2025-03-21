package db

import (
	"finalizacao-pedido-svc/infra/logger"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	go_ora "github.com/sijms/go-ora/v2"
)

var SQLDB *sql.DB
var OracleDB *sql.DB

func ConnectSQLServer() {
	// Obtendo as credenciais do ambiente
	server := os.Getenv("SQL_SERVER_URL")
	user := os.Getenv("SQL_SERVER_USER")
	password := os.Getenv("SQL_SERVER_PASSWORD")
	database := os.Getenv("SQL_SERVER_DATABASE")

	// Criando a string de conexão
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
		server, user, password, database)

	// Abrindo a conexão com o banco de dados
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		logger.Log.Fatal("Erro ao abrir conexão com o SQL Server:", zap.Error(err))
	}

	// Testando a conexão
	if err := db.Ping(); err != nil {
		logger.Log.Fatal("Erro ao conectar ao SQL Server:", zap.Error(err))
	}

	logger.Log.Info("✅ Conexão bem-sucedida com o SQL Server!")

	SQLDB = db
}
func ConnectOracle() {
	username := os.Getenv("ORACLE_USER")
	password := os.Getenv("ORACLE_PASSWORD")
	hostname := os.Getenv("ORACLE_URL")
	port, _ := strconv.ParseInt(os.Getenv("ORACLE_PORT"), 10, 0)
	serviceName := os.Getenv("ORACLE_SERVICE_NAME")
	connStr := go_ora.BuildUrl(hostname, int(port), serviceName, username, password, nil)
	db, err := sql.Open("oracle", connStr)
	if err != nil {
		logger.Log.Fatal("Erro ao conectar ao Oracle", zap.Error(err))
	}

	err = db.Ping()
	if err != nil {
		logger.Log.Fatal("Erro ao conectar ao Oracle", zap.Error(err))
	}
	OracleDB = db
	logger.Log.Info("✅ Conexão bem-sucedida com o Oracle!")
}

func CloseConnections() {
	if SQLDB != nil {
		_ = SQLDB.Close()
	}
	if OracleDB != nil {
		_ = OracleDB.Close()
	}
}
