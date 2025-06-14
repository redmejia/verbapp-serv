package models

type Chat struct {
	Prompt        TextPrompt    `json:"prompt,omitempty"`
	ReplyText     ReplyText     `json:"reply_text,omitempty"` // reply to ai generated text
	GeneratedText GeneratedText `json:"generated_text"`
}

// User prmpt text
type TextPrompt struct {
	ChatID         string `json:"chat_id,omitempty"` // Optional, can be used to link to a specific chat
	UserID         string `json:"user_id,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	Timestamp      int64  `json:"timestamp,omitempty"`
	Text           string `json:"text,omitempty"`
}

// Response from AI
type GeneratedText struct {
	ChatID         string   `json:"chat_id,omitempty"` // Optional, can be used to link to a specific chat
	UserID         string   `json:"user_id,omitempty"`
	ConversationID string   `json:"conversation_id,omitempty"`
	Timestamp      int64    `json:"timestamp,omitempty"`
	Text           string   `json:"text,omitempty"`
	Metadata       Metadata `json:"metadata,omitempty"`
}

// Reply prompt text response USER reply to generated text
type ReplyText struct {
	ChatID         string `json:"chat_id,omitempty"` // Optional, can be used to link to a specific chat
	UserID         string `json:"user_id,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	Timestamp      int64  `json:"timestamp,omitempty"`
	Reply          string `json:"reply,omitempty"` // Generareted text
	Text           string `json:"text,omitempty"`  // user prompt
}

type Metadata struct {
	ModelName              string `json:"model_name"`
	FinishReason           string `json:"finish_reason,omitempty"`
	PromptTokenCount       int    `json:"promt_token_count,omitempty"`
	CompletationTokenCount int    `json:"conpletation_token_count,omitempty"`
	TotalTokenCount        int    `json:"total_token_count,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
