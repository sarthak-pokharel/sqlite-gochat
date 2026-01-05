package models

import "time"


type WebhookEvent struct {
	ID          int64      `json:"id" db:"id"`
	ChannelID   int64      `json:"channel_id" db:"channel_id"`
	EventType   string     `json:"event_type" db:"event_type"`
	Payload     string     `json:"payload" db:"payload"`
	Processed   bool       `json:"processed" db:"processed"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	Error       *string    `json:"error,omitempty" db:"error"`
}
