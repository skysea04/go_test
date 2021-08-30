package db_client

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DBClient *sql.DB

func InitialiseDBConnection() {
	db, err := sql.Open("mysql", "root:root@/test")
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	DBClient = db
}
