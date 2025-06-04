package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Store struct {
	Db                *sql.DB
	InfoLog, ErrorLog *log.Logger
}

const (
	OpenConns = 10
	IdleConns = 3
	LifeTime  = 60 * time.Second
)

func StoreConnection() (*sql.DB, error) {

	port, _ := strconv.Atoi(os.Getenv("DBPORT"))
	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		os.Getenv("HOSTNAME"), port, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("pgx", connection)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db.SetMaxOpenConns(OpenConns)
	db.SetMaxIdleConns(IdleConns)
	db.SetConnMaxLifetime(LifeTime)

	return db, err

}
