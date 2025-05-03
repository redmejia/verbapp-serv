package router

import (
	"net/http"

	"github.com/redmejia/internal/handlers"
)

func Router(app *handlers.App) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/chat", app.ChatHandler)
	return mux
}
