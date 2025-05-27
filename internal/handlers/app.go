package handlers

import "log"

type App struct {
	InfoLog, ErrorLog *log.Logger
	ApiKey            string // not needed
	JwtKey            string
}
