package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMySqlDb(config mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = checkDbConnection(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func checkDbConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}

	log.Println("Database connection initialized")
	return nil
}
