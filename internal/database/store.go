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
	GetAllChats() []models.Chat
	GetPromptByConversationID(conversationID string) (string, string, error)
	InsertGeneratedText(modelName, userID, conversationID, generatedText string) models.GeneratedText
	// InsertRepley() string
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

func (s *Store) GetPromptByConversationID(conversationID string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var prompt string
	query := `select text from prompts where conversation_id = $1`

	row := s.Db.QueryRowContext(ctx, query, conversationID)

	err := row.Scan(&prompt)
	if err != nil {
		return "", "", err
	}
	return prompt, conversationID, nil
}

func (s *Store) InsertGeneratedText(modelName, userID, conversationID, generatedText string) models.GeneratedText {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := s.Db.BeginTx(ctx, nil)

	if err != nil {
		s.ErrorLog.Println("error starting transaction:", err)
	}

	var generatedTextResponse models.GeneratedText

	queryGeneratedText := `INSERT INTO generated_texts 
		(user_id, conversation_id, timestamp, text) 
		VALUES ( $1, $2, $3, $4 ) RETURNING chat_id, user_id, conversation_id, timestamp, text`

	row := tx.QueryRowContext(ctx, queryGeneratedText, userID, conversationID, time.Now().Unix(), generatedText)

	_ = row.Scan(
		&generatedTextResponse.ChatID,
		&generatedTextResponse.UserID,
		&generatedTextResponse.ConversationID,
		&generatedTextResponse.Timestamp,
		&generatedTextResponse.Text,
	)

	queryMetadata := `INSERT INTO response_metadata (conversation_id, model_name) VALUES ( $1, $2 ) RETURNING model_name`
	row = tx.QueryRowContext(ctx, queryMetadata, conversationID, modelName)

	_ = row.Scan(&generatedTextResponse.Metadata.ModelName)

	if err = tx.Commit(); err != nil {
		s.ErrorLog.Fatal("error committing transaction:", err)
	}

	return generatedTextResponse
}

func (s *Store) GetAllChats() []models.Chat {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `select 
		p.chat_id,
		p.user_id,
		p.conversation_id,
		p.timestamp,
		p.text,
		gt.chat_id,
		gt.user_id,
		gt.conversation_id,
		gt.timestamp,
		gt.text,
		mt.model_name
	from prompts p
	join generated_texts gt on gt.conversation_id = p.conversation_id 
	join response_metadata mt on mt.conversation_id = p.conversation_id
	order by p.id desc`

	rows, err := s.Db.QueryContext(ctx, query)
	if err != nil {
		s.ErrorLog.Println("error querying all chats:", err)
	}

	var chats []models.Chat
	for rows.Next() {

		var prompt models.TextPrompt
		var generatedText models.GeneratedText

		err := rows.Scan(
			&prompt.ChatID,
			&prompt.UserID,
			&prompt.ConversationID,
			&prompt.Timestamp,
			&prompt.Text,
			&generatedText.ChatID,
			&generatedText.UserID,
			&generatedText.ConversationID,
			&generatedText.Timestamp,
			&generatedText.Text,
			&generatedText.Metadata.ModelName,
		)

		if err != nil {
			s.ErrorLog.Println("error scanning row:", err)
			continue
		}

		chats = append(chats, models.Chat{
			Prompt:        prompt,
			GeneratedText: generatedText,
		})
	}

	return chats
}
