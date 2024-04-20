package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() {
	db, err := sql.Open("sqlite3", "./go-crypto-price-tracker.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	DB = db
}
