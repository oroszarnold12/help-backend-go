package main

import (
	"help/cmd/api"
	"help/config"
	"help/db"
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {
	dbConfig := mysql.Config{
		User:                 config.Env.DbUser,
		Passwd:               config.Env.DbPassword,
		DBName:               config.Env.DbName,
		Addr:                 config.Env.DbAddress,
		AllowNativePasswords: true,
	}
	mysqlDb, err := db.NewMySqlDb(dbConfig)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	api := api.NewApi(config.Env.ApiPort, mysqlDb)
	if err := api.Run(); err != nil {
		log.Fatalf("Cannot create api: %v", err)
	}
}
