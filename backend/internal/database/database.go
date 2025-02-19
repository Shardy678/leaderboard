package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	port     = os.Getenv("DB_PORT")
)

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
