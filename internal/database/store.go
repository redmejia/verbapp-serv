package database

import (
	"database/sql"
	"log"
)

type Store struct {
	Db                *sql.DB
	InfoLog, ErrorLog *log.Logger
}

type ChatStore interface {
	InsertText() string
	InsertGeneratedText() string
	InsertRepley() string
}

func (s *Store) InsertText() string {
	return ""
}
func (s *Store) InsertGeneratedText() string {
	return ""
}

func (s *Store) InsertRepley() string {
	return ""
}
