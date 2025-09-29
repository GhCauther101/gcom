package main

import (
	// "fmt"
	"fmt"
	"gcom/config"
	"gcom/db"
	"log"
	"os"
	"strconv"

	mysqlConfig "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := db.NewMySQLStorage(mysqlConfig.Config{
		User:                 config.Envs.DbUser,
		Passwd:               config.Envs.DbPassword,
		Addr:                 config.Envs.DbAddress,
		DBName:               config.Envs.DbName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	drive, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	
	m, err :=migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		drive,
	)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[1]

	if cmd == "up" {
		if err :=m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		if err :=m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	//check
	if cmd == "apply" {
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if err :=m.Force(version); err != nil && err != nil {
			log.Fatal(err)
		}
	}
}