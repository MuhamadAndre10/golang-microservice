package main

import (
	"auth-service/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth-service API server on port", webPort)

	// connect to database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to database. Exiting...")
		return
	}
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	log.Println("Listening on port", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Could not connect to database. Retrying...")
			counts++
		} else {
			log.Println("Connected to database")
			return connection
		}

		if counts > 10 {
			log.Println("Could not connect to database. Giving up")
			return nil
		}

		log.Println("Backing off for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}

}
