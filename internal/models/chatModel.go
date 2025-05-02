package models

import "time"

// Human ask question to machine
type Human struct {
	UserID         string    `json:"user_id"`
	ConversationID string    `json:"conversation_id,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
	Text           string    `json:"text"`
}

// Machine response to the human text question
type Machine struct {
	UserID         string    `json:"user_id"`
	ConversationID string    `json:"conversation_id,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
	Text           string    `json:"text"`
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
