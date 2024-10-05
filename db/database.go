package db

import (
	"Nookhub/config"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Initialize() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.GetString("user"),
		config.GetString("password"),
		config.GetString("host"),
		config.GetString("port"),
		config.GetString("name"))

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	if err = DB.Ping(); err != nil {
		fmt.Println(dsn)
		panic("Database connection error: " + err.Error())
	}

	fmt.Print(dsn)
	fmt.Print("DB connection established successfully\n\n")
}
