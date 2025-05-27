package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redmejia/internal/models"
)

func (app *App) ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle POST request
		var human models.TextPrompt
		err := json.NewDecoder(r.Body).Decode(&human)
		if err != nil {
			app.ErrorLog.Println("error decoding rquested body", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
			})
		}

		var time = time.Now().Unix()
		human.Timestamp = time

		jsonByte, _ := json.Marshal(&human)
		app.InfoLog.Println("PROMPT>> ", string(jsonByte))

		chatID := uuid.New()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.GeneratedText{
			ChatID:    chatID.String(),
			Text:      "Nice! That sounds fun.",
			Timestamp: time,
			Metadata: models.Metadata{
				ModelName: "verbapp v1.0.0 light",
			},
		})

	}
}
