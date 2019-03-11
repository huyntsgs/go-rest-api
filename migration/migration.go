package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// migration inits database schema. Creates necessary tables.
// Needs to setup correct database information to .env file and then runs migration.
// Database and tables will be created.
func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")
	dbServer := os.Getenv("DB_SERVER")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPass, dbServer, dbPort)
	db, err := sql.Open(dbDriver, dbUrl)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer db.Close()
	fmt.Println("Open db connection successfully")

	_, err = db.Exec("DROP SCHEMA if exists store")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Dropped database")
	_, err = db.Exec("CREATE DATABASE store")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Database created successfully")

	_, err = db.Exec("USE store")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("DB store is selected")

	_, err = db.Exec(`CREATE Table products(product_id int unsigned NOT NULL AUTO_INCREMENT, product_name nvarchar(100), image varchar(200), 
		rate tinyint unsigned, created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (product_id))`)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Table products created successfully")

	_, err = db.Exec(`CREATE Table users(user_id int unsigned NOT NULL AUTO_INCREMENT, user_name nvarchar(100), email varchar(200),
	 	password varchar(200), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (user_id))`)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Table users created successfully")

}
