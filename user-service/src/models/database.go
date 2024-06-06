package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDB() {
	var err error
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"))
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}
