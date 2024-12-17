package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type Db *sql.DB

var (
	dbname   = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")

	connMaxLifetime = 5 * time.Minute
	maxOpenConns    = 10
	maxIdleConns    = 5
)

func New() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbname))
	if err != nil {
		// will be DSN parse error or some initialization error.
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	if err = TestConnection(db); err != nil {
		return nil, err
	}

	return db, nil
}

func TestConnection(d *sql.DB) error {
	if err := d.Ping(); err != nil {
		fmt.Println("Database ping failed: ", err)
	}

	fmt.Println("Database ping successful")
	return nil
}
