package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redmejia/internal/handlers"
	"github.com/redmejia/internal/router"
	"github.com/redmejia/internal/security"
)

func main() {

	var (
		port   string
		host   string
		jwtKey string
		userID string // test debug
	)

	defaultPort := "8080"
	defaultHost := "127.0.0.1"
	flag.StringVar(&port, "port", defaultPort, "Sever port")
	flag.StringVar(&host, "host", defaultHost, "Sever host")
	flag.StringVar(&jwtKey, "key", "", "JWT key")
	flag.StringVar(&userID, "uid", "", "User id")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// token gerenation when user registers
	token, err := security.GenerateToken(jwtKey, userID)
	if err != nil {
		errorLog.Fatal("error generating token:", err)
		os.Exit(1)
	}

	infoLog.Println("TOKEN: ", token)

	app := &handlers.App{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		JwtKey:   jwtKey,
	}

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%s", port),
		ErrorLog: app.ErrorLog,
		Handler:  router.Router(app),
	}

	app.InfoLog.Println("Starting server on port", port)
	if err := srv.ListenAndServe(); err != nil {
		app.ErrorLog.Fatal(err)
		os.Exit(1)
	}

}
