package handlers

import "log"

type App struct {
	Host              string
	Port              string
	InfoLog, ErrorLog *log.Logger
	ApiKey            string
	JwtKey            string
}
