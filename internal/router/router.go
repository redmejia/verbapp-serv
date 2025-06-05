package router

import (
	"net/http"

	"github.com/redmejia/internal/handlers"
	"github.com/redmejia/internal/middleware"
)

func Router(app *handlers.App) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/chat", middleware.IsAuthorized(app, app.ChatHandler))
	mux.HandleFunc("/v1/chat/prompt", middleware.IsAuthorized(app, app.PromptHandler))
	mux.HandleFunc("/v1/chat/generated/text", middleware.IsAuthorized(app, app.PromptHandler))
	// mux.HandleFunc("/v1/chat/generated/image", middleware.IsAuthorized(app, app.PromptHandler))

	return middleware.Logger(mux)
}
