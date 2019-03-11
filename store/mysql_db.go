package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type (
	MySqlDB struct {
		DB *sql.DB
	}
)

// Connect creates connection to database server.
// The database informations are loaded from .env file.
func (mysql *MySqlDB) Connect() {
	log.Println("Start db connection")
	var err error

	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbServer := os.Getenv("DB_SERVER")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbServer, dbPort, dbName)

	mysql.DB, err = sql.Open(dbDriver, dbUrl)

	if err != nil {
		log.Println("Error in db connection")
		panic(err.Error())
	} else {
		log.Println("Connected to db")
	}
}

// InitDB loads .env file to environment variable and
// starts connect to database.
func InitDB() *MySqlDB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Can not load env file")
	}
	mysql := new(MySqlDB)
	mysql.Connect()
	return mysql
}
