package main

import (
	"log"
	"net/http"
	"os"

	"github.com/redmejia/internal/handlers"
	"github.com/redmejia/internal/router"
)

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &handlers.App{
		Host:     "localhost",
		Port:     ":8080",
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	srv := &http.Server{
		Addr:     app.Port,
		ErrorLog: app.ErrorLog,
		Handler:  router.Router(app),
	}
	app.InfoLog.Println("Starting server on port", app.Port)
	if err := srv.ListenAndServe(); err != nil {
		app.ErrorLog.Fatal(err)
		os.Exit(1)
	}

}
