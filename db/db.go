package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	// "github.com/bwmarrin/discordgo"
	"errors"

	// "path/filepath"
	"strings"
	// "github.com/yakabuff/discord-dl/models"
)

const POSTGRES = "postgres"
const SQLITE = "sqlite3"

var UniqueConstraintError = errors.New("Unique constraint error")

type Db struct {
	DbConnection *sql.DB
}

func Init_db(path string) (*Db, error) {
	var err error
	//Check if DB exists
	if path == "" {
		//If path empty, fallback and check default db location
		_, err = os.Stat("archive.db")
	} else {
		_, err = os.Stat(path)
	}

	driver := determineDbType(path)
	log.Println(path)
	var dbConn *sql.DB
	var file *os.File
	if err == nil {
		//Exists
		if path == "" {
			dbConn, err = sql.Open(driver, "archive.db?_foreign_keys=on")
		} else {
			dbConn, err = sql.Open(driver, path+"?_foreign_keys=on")
		}
		if err != nil {
			return nil, err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if path == "" {
			path = "archive.db"
			file, err = os.Create(path)
		} else {
			file, err = os.Create(path) // Create SQLite file
		}
		if err != nil {
			log.Println("could not create db file")
			log.Fatal(err.Error())
		}
		file.Close()
		dbConn, err = sql.Open("sqlite3", path+"?_foreign_keys=on")
		if err != nil {
			return nil, err
		}
		createTable(dbConn)
		//*message_id | channel_id | guild| | date | content | media | sender_id | reply_to //
		// 234234242  | 23489353   | 324242 | 1231 |asdfasdfs | <urL> | 234242 | 234756//
	} else {
		//Panic
		log.Fatal(err.Error())
	}
	db := Db{DbConnection: dbConn}
	return &db, err
}
func determineDbType(path string) string {
	if strings.HasPrefix(path, "postgres://") {
		return POSTGRES
	}
	return SQLITE
}
