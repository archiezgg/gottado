package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbHost     = "database"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "root"
	dbName     = "gottado"
)

var db *sql.DB

func initDB() {
	var err error

	dbInfo := fmt.Sprintf("host=%v port=%v user=%v "+
		"password=%v dbname=%v sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = sql.Open("postgres", dbInfo)
	for err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		try := 1
		for try <= 6 && err != nil {
			log.Printf("Establishing connection to the database... %d\nExiting after 5 tries.", try)
			time.Sleep(10 * time.Second)
			err = db.Ping()
			try++
			if try == 6 {
				panic(err)
			}
		}
	}
	log.Println("Successfully connected to the database.")

	createTable()
}

func createTable() {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title VARCHAR(20) NOT NULL,
			content VARCHAR(255) NOT NULL,
			urgent BOOLEAN NOT NULL
		);
	`
	_, err := db.Exec(query)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Table 'tasks' has been created.")
	}
}
