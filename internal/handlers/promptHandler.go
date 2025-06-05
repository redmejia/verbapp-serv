package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/redmejia/internal/models"
)

func (app *App) PromptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle POST request
		userID := r.Context().Value("user_id").(string)

		var prompt models.TextPrompt
		err := json.NewDecoder(r.Body).Decode(&prompt)
		if err != nil {
			app.ErrorLog.Println("error decoding requested body", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
			})
		}

		prompt.UserID = userID
		prompt.Timestamp = time.Now().Unix()

		err = app.DB.InsertPrompt(&prompt)
		if err != nil {
			app.ErrorLog.Println("error inserting prompt into database", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(prompt)

	}
}
