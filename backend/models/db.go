package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:ops123!@tcp(localhost:3306)/browser_monitor")
	if err != nil {
		panic(err)
	}
}
