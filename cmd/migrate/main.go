package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/zhetkerbaevan/personal-blog/internal/config"
	"github.com/zhetkerbaevan/personal-blog/internal/db"
)

func main() {
	db, err := db.NewPostgreSQLStorage(config.Config{
			DBHost: config.Envs.DBHost,
			DBPort: config.Envs.DBPort,
			DBUser: config.Envs.DBUser,
			DBName: config.Envs.DBName,
			DBPassword: config.Envs.DBPassword,
		})
	if err != nil {
		log.Fatal(err)
	}	

	//create driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}	
	
	//create migration
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver)
	if err != nil {
		log.Fatal(err)
	}		

	//get last argument from cmd
	cmd := os.Args[len(os.Args) - 1]
	
	//migrate up or down
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("migrate-up")
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("migrate-down")
	}
}