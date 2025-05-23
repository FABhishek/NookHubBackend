package db

import (
	"Nookhub/config"
	"database/sql"
	"fmt"
	"os"

	"log"

	"cloud.google.com/go/firestore"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var FirebaseClient *firestore.Client
var ConfigData []byte

const (
	firebaseConfigFile = "firebaseconfig.json"
)

func Initialize() {
	// connStr := "user=username password=password dbname=mydatabase host=localhost sslmode=disable"
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s",
		config.GetString("user"),
		config.GetString("password"),
		config.GetString("name"),
		config.GetString("host"),
		config.GetString("port"))

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	if err = DB.Ping(); err != nil {
		fmt.Println(dsn)
		panic("Database connection error: " + err.Error())
	}

	fmt.Print(dsn)
	fmt.Print("Postgres DB connection established successfully\n\n")

	//firebase
	ConfigData, err = os.ReadFile(firebaseConfigFile)
	if err != nil {
		log.Fatalf("Failed to read service account file: %v", err)
	}
}
