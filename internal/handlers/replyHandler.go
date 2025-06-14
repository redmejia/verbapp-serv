package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/redmejia/internal/models"
)

func (app *App) ReplyGeneratedTextHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

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

		// GET the text from generated_texts
		generatedText, err := app.DB.GetGeneratedTextByConversationID(prompt.ConversationID)
		if err != nil {
			app.ErrorLog.Println("error getting conversation text", err)
		}

		app.InfoLog.Println(generatedText)

		var reply models.ReplyText
		reply.UserID = userID
		reply.Reply = generatedText
		reply.Text = prompt.Text
		reply.Timestamp = time.Now().Unix()

		// insert record reply generated text and the reply user text into reply_prompts
		err = app.DB.InsertReplyPromptWithReplyText(&reply)
		if err != nil {
			app.ErrorLog.Println("error inserting reply text into database", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(reply)

	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:    http.StatusMethodNotAllowed,
			Message: "The HTTP method used is not supported for this resource.",
		})
	}

}
