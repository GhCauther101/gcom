package main

import (
	"database/sql"
	"gcom/cmd/api"
	"gcom/config"
	"gcom/db"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db, err := db.NewMySQLStorage(mysql.Config{
		User: config.Envs.DbUser,
		Passwd: config.Envs.DbPassword,
		Addr: config.Envs.DbAddress,
		DBName: config.Envs.DbName,
		Net: "tcp",
		AllowNativePasswords: true,
		ParseTime: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	apiServer := api.NewAPIServer(":8080", db)
	if err := apiServer.Run(); err != nil {
		log.Fatal("[api] could not launch api server")
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[db] succefully connected")
}