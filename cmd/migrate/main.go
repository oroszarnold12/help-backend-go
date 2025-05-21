package main

import (
	"help/config"
	"help/db"
	"log"
	"os"

	mySqlDriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mySqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbConfig := mySqlDriver.Config{
		User:                 config.Env.DbUser,
		Passwd:               config.Env.DbPassword,
		DBName:               config.Env.DbName,
		Addr:                 config.Env.DbAddress,
		AllowNativePasswords: true,
	}
	mySqlDb, err := db.NewMySqlDb(dbConfig)
	if err != nil {
		log.Fatalf("Cannot connect to database %v", err)
	}

	driver, err := mySqlMigrate.WithInstance(mySqlDb, &mySqlMigrate.Config{})
	if err != nil {
		log.Fatalf("Cannot create database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("Cannot to create migration: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Please specify migration type: up | down")
	}

	migrationType := os.Args[len(os.Args)-1]
	switch migrationType {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration successful")
	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration successful")
	default:
		log.Fatalf("Wrong migration type: %s. Correct value are up | down", migrationType)
	}
}
