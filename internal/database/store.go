package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/redmejia/internal/models"
)

type Store struct {
	Db                *sql.DB
	InfoLog, ErrorLog *log.Logger
}

type ChatStore interface {
	InsertPrompt(prompt *models.TextPrompt) error
	InsertGeneratedText() string
	InsertRepley() string
}

func (s *Store) InsertPrompt(prompt *models.TextPrompt) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var converstionID string
	var chatID string

	query := `INSERT INTO prompts (user_id, timestamp, text) VALUES ( $1, $2, $3 ) RETURNING chat_id, conversation_id`
	row := s.Db.QueryRowContext(ctx, query, prompt.UserID, prompt.Timestamp, prompt.Text)

	err := row.Scan(&chatID, &converstionID)

	if err != nil {
		return err
	}

	prompt.ChatID = chatID
	prompt.ConversationID = converstionID

	return nil
}

func (s *Store) InsertGeneratedText() string {
	return ""
}

func (s *Store) InsertRepley() string {
	return ""
}
