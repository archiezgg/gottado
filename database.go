package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
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
	for i := 0; i <= 1000 && err != nil; i++ {
		if i == 0 {
			log.Println("Trying to connect to the database...")
		}
	}

	if err = db.Ping(); err != nil {
		panic(err)
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
