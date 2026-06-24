package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectToDB() {
	// Carrega .env se existir (não-fatal — em produção usa env vars do sistema)
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("Aviso: arquivo .env não encontrado, usando variáveis de ambiente do sistema")
		}
	}

	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		// Fallback para variáveis individuais
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")

		sslmode := os.Getenv("DB_SSLMODE")

		connStr = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
			user, password, dbname, host, port, sslmode)
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Testa a conexão
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Erro ao verificar a conexão com o banco de dados: %v", err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")
}
