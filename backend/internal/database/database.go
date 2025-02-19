package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type Service struct {
	db *sql.DB
}

var (
	database string
	username string
	password string
	port     string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	database = os.Getenv("DB_DATABASE")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	port = os.Getenv("DB_PORT")
}

func (s *Service) Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", username, password, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func (s *Service) Close() error {
	log.Println("Closing database connection")
	return s.db.Close()
}
