package models

type Chat struct {
	ChatID        string        `json:"chat_id"`
	Prompt        TextPrompt    `json:"prompt"`
	GeneratedText GeneratedText `json:"generated_text"`
}

// User prmpt text
type TextPrompt struct {
	ChatID         string `json:"chat_id,omitempty"` // Optional, can be used to link to a specific chat
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id,omitempty"`
	Timestamp      int64  `json:"timestamp"`
	Text           string `json:"text"`
}

// Response from AI
type GeneratedText struct {
	ChatID         string   `json:"chat_id,omitempty"` // Optional, can be used to link to a specific chat
	UserID         string   `json:"user_id"`
	ConversationID string   `json:"conversation_id,omitempty"`
	ResponseID     string   `json:"response_id,omitempty"`
	Timestamp      int64    `json:"timestamp"`
	Text           string   `json:"text"`
	Metadata       Metadata `json:"metadata,omitempty"`
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
