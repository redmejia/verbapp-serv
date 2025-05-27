package router

import (
	"net/http"

	"github.com/redmejia/internal/handlers"
	"github.com/redmejia/internal/middleware"
)

func Router(app *handlers.App) http.Handler {
	mux := http.NewServeMux()

	// mux.HandleFunc("/v1/chat", middleware.IsAuthorized(app, app.ChatHandler))
	mux.HandleFunc("/v1/chat", app.ChatHandler)

	return middleware.Logger(mux)
}
