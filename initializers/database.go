package initializers

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DatabaseConnection() {
	database, err := sql.Open("mysql", "root:@tcp(localhost:3306)/belajar_api")

	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	DB = database
}
