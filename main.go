// to open postgres cli use 'psql'

package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func initDB() {

	err1 := godotenv.Load() // This will look for the .env file in the current directory
	if err1 != nil {
		log.Fatal("Error loading .env file")
	}

	dbpath := fmt.Sprintf("user=%s dbname=%s host=%s port=%s sslmode=%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBNAME"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBSSL"))

	fmt.Println("connected to database")

	connStr := dbpath //"user=axellmartinez dbname=mydb  host=localhost port=5432 sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	initDB()

	fmt.Println("hello world")
}
