package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redmejia/internal/database"
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

	defaultPort := os.Getenv("SERVER_PORT")
	defaultHost := os.Getenv("SERVER_HOST")
	jwtSecretKey := os.Getenv("JWT_KEY")
	uid := os.Getenv("USER_ID")

	flag.StringVar(&port, "port", defaultPort, "Sever port")
	flag.StringVar(&host, "host", defaultHost, "Sever host")
	flag.StringVar(&jwtKey, "key", jwtSecretKey, "JWT key")
	flag.StringVar(&userID, "uid", uid, "User id")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := database.StoreConnection()
	if err != nil {
		errorLog.Fatal("error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// token gerenation when user registers
	token, err := security.GenerateToken(jwtKey, userID)
	if err != nil {
		errorLog.Fatal("error generating token:", err)
		os.Exit(1)
	}

	infoLog.Println("TOKEN: ", token)

	geminiKey := os.Getenv("GEMINI_API_KEY")

	infoLog.Println("PORT:", defaultPort)
	infoLog.Println("HOST:", defaultHost)
	infoLog.Println("JWT_SECRET:", jwtSecretKey)
	infoLog.Println("USER_ID:", uid)

	app := &handlers.App{
		InfoLog:   infoLog,
		ErrorLog:  errorLog,
		GeminiKey: geminiKey,
		JwtKey:    jwtKey,
		DB: &database.Store{
			Db:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
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
