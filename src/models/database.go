package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDB() {
	var err error
	connStr := "user=youruser dbname=yourdb password=yourpassword host=yourhost sslmode=disable"
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}
