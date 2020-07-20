package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "golangdemo"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	CheckForError(err)

	return db
}
