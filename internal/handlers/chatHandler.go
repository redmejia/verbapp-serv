package handlers

import (
	"encoding/json"
	"net/http"
)

func (app *App) ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		chats := app.DB.GetAllChats()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(chats)
	}
}
