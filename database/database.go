package database

import (
	"database/sql"
	"os"
	"strings"

	_ "github.com/lib/pq" // Don't forget the driver!
)

type Database struct {
	conn *sql.DB
}

func CheckDbError(err error, message ...string) {
	if err != nil {
		message = append(message, err.Error())
		panic(strings.Join(message[:], " | "))
	}
}

func (db *Database) Init() {

	dbPath := os.Getenv("DB_URL")
	conn, err := sql.Open("postgres", dbPath)

	CheckDbError(err, "Error connecting to database")

	err = conn.Ping()
	CheckDbError(err, "Error pinging database")

	db.conn = conn
}

func (db *Database) Close() {
	defer db.conn.Close()
}
