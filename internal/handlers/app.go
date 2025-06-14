package handlers

import (
	"log"

	"github.com/redmejia/internal/database"
)

const (
	REPLY  = "reply"
	PROMPT = "prompt"
)

type App struct {
	InfoLog, ErrorLog *log.Logger
	GeminiKey         string
	JwtKey            string
	DB                database.ChatStore
}
