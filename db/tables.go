package db

import (
	"database/sql"
	"log"
)

func createTable(db *sql.DB) {

	log.Println("Creating tables")
	_, errx := db.Exec(schema)
	if errx != nil {
		log.Fatal(errx.Error())
	}
	log.Println("Tables initialized")

}
