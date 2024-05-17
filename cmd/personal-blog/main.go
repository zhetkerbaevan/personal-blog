package main

import (
	"fmt"
	"log"

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
	fmt.Println("connected to db", db)
}