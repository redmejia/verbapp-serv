package handlers

import "log"

type App struct {
	InfoLog, ErrorLog *log.Logger
	GeminiKey         string
	JwtKey            string
}
