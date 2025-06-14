package router

import (
	"net/http"

	"github.com/redmejia/internal/handlers"
	"github.com/redmejia/internal/middleware"
)

func Router(app *handlers.App) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/chats", middleware.IsAuthorized(app, app.ChatHandler))                    // GET
	mux.HandleFunc("/v1/chat/prompt", middleware.IsAuthorized(app, app.PromptHandler))            // POST
	mux.HandleFunc("/v1/chat/reply", middleware.IsAuthorized(app, app.ReplyGeneratedTextHandler)) // POST
	mux.HandleFunc("/v1/chat/ai/resp", middleware.IsAuthorized(app, app.AITextHandler))           // GET

	// test
	// mux.HandleFunc("/v1/prompt", app.PromptHandler)
	// mux.HandleFunc("/v1/chat/ai/text", middleware.IsAuthorized(app, app.AITextHandler))
	// mux.HandleFunc("/v1/chat/ai/image", middleware.IsAuthorized(app, app.PromptHandler))

	return middleware.Logger(mux)
}
