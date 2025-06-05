package handlers

import (
	"log"

	"github.com/redmejia/internal/database"
)

type App struct {
	InfoLog, ErrorLog *log.Logger
	GeminiKey         string
	JwtKey            string
	DB                database.ChatStore
}
