package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/redmejia/internal/models"
)

type Store struct {
	Db                *sql.DB
	InfoLog, ErrorLog *log.Logger
}

type ChatStore interface {
	GetAllChats() []models.Chat
	GetPromptByConversationID(conversationID string) (promptFound bool, text string, id string, err error)
	GetReplyPromptByConversationID(conversationID string) (string, string, error)
	GetGeneratedTextByConversationID(conversationID string) (string, error) // get the generated text to insert into reply table
	InsertPrompt(prompt *models.TextPrompt) error                           // prompt with user input text
	InsertReplyPromptWithReplyText(replyPrompt *models.ReplyText) error     // insert prompt with reply ai generated text
	InsertGeneratedText(modelName, userID, conversationID, generatedText string) models.GeneratedText
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
    
		rp.chat_id,
		rp.user_id,
		rp.conversation_id,
		rp.timestamp,
		rp.reply,
		rp.text,
    
		gt.chat_id,
		gt.user_id,
		gt.conversation_id,
		gt.timestamp,
		gt.text,
    
		COALESCE(mt.model_name,'Gemini 2.0-flash'),
    
    GREATEST(
        COALESCE(p.timestamp, 0),
        COALESCE(rp.timestamp, 0),
        COALESCE(gt.timestamp, 0)
    ) as last_timestamp
    
	from prompts p
	right join generated_texts gt on gt.conversation_id = p.conversation_id 
	left join reply_prompts rp on rp.conversation_id = gt.conversation_id
	left join response_metadata mt on mt.conversation_id = p.conversation_id
	order by last_timestamp asc`

	rows, err := s.Db.QueryContext(ctx, query)
	if err != nil {
		s.ErrorLog.Println("error querying all chats:", err)
	}

	var chats []models.Chat
	for rows.Next() {
		// var prompt models.TextPrompt
		// var reply models.ReplyText
		// var generatedText models.GeneratedText

		var (
			// prompt
			promptChatID,
			promptUserID,
			promptConversationID,
			promptText sql.NullString
			promptTimestamp sql.NullInt64

			// reply
			replyChatID,
			replyUserID,
			replyConversationID,
			replyReply,
			replyText sql.NullString
			replyTimestamp sql.NullInt64

			// generated text
			genChatID,
			genUserID,
			genConversationID,
			genText,
			genModelName sql.NullString
			genTimestamp sql.NullInt64

			lastTimestamp sql.NullInt64 // not used
		)

		err := rows.Scan(
			&promptChatID,
			&promptUserID,
			&promptConversationID,
			&promptTimestamp,
			&promptText,
			&replyChatID,
			&replyUserID,
			&replyConversationID,
			&replyTimestamp,
			&replyReply,
			&replyText,
			&genChatID,
			&genUserID,
			&genConversationID,
			&genTimestamp,
			&genText,
			&genModelName,
			&lastTimestamp,
		)

		if err != nil {
			s.ErrorLog.Println("error scanning row:", err)
			continue
		}

		chats = append(chats, models.Chat{
			Prompt: models.TextPrompt{
				ChatID:         promptChatID.String,
				UserID:         promptUserID.String,
				ConversationID: promptConversationID.String,
				Timestamp:      promptTimestamp.Int64,
				Text:           promptText.String,
			},
			ReplyText: models.ReplyText{
				ChatID:         replyChatID.String,
				UserID:         replyUserID.String,
				ConversationID: replyConversationID.String,
				Timestamp:      replyTimestamp.Int64,
				Reply:          replyReply.String,
				Text:           replyText.String,
			},
			GeneratedText: models.GeneratedText{
				ChatID:         genChatID.String,
				UserID:         genUserID.String,
				ConversationID: genConversationID.String,
				Timestamp:      genTimestamp.Int64,
				Text:           genText.String,
				Metadata: models.Metadata{
					ModelName: genModelName.String,
				},
			},
		})
	}

	return chats
}

func (s *Store) GetPromptByConversationID(conversationID string) (bool, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var prompt string
	query := `select text from prompts where conversation_id = $1`

	row := s.Db.QueryRowContext(ctx, query, conversationID)

	err := row.Scan(&prompt)
	if err != nil {
		return false, "", "", err
	}

	return true, prompt, conversationID, nil
}

func (s *Store) GetReplyPromptByConversationID(conversationID string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `select reply, text from reply_prompts where conversation_id = $1`
	row := s.Db.QueryRowContext(ctx, query, conversationID)

	var (
		reply               string
		text                string
		formatedPromptReply string
	)

	err := row.Scan(&reply, &text)

	if err != nil {
		return "", "", err
	}

	formatedPromptReply = fmt.Sprintf("%s\n%s", reply, text)

	return formatedPromptReply, conversationID, nil

}

func (s *Store) GetGeneratedTextByConversationID(conversationID string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var generatedText string
	query := "select text from generated_texts where conversation_id = $1"

	row := s.Db.QueryRowContext(ctx, query, conversationID)

	err := row.Scan(&generatedText)

	if err != nil {
		return "", err
	}

	return generatedText, nil
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

func (s *Store) InsertReplyPromptWithReplyText(replyPrompt *models.ReplyText) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var converstionID string
	var chatID string

	query := `INSERT INTO reply_prompts (user_id, timestamp, reply, text) VALUES ( $1, $2, $3, $4 ) RETURNING chat_id, conversation_id`
	row := s.Db.QueryRowContext(
		ctx,
		query,
		replyPrompt.UserID,
		replyPrompt.Timestamp,
		replyPrompt.Reply,
		replyPrompt.Text,
	)

	err := row.Scan(&chatID, &converstionID)

	if err != nil {
		return err
	}

	replyPrompt.ChatID = chatID
	replyPrompt.ConversationID = converstionID

	return nil
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
