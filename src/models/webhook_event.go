package models

import "time"


type WebhookEvent struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	ChannelID   int64      `json:"channel_id" gorm:"not null;index"`
	EventType   string     `json:"event_type" gorm:"not null"`
	Payload     string     `json:"payload" gorm:"not null;type:text"`
	Processed   bool       `json:"processed" gorm:"default:false;index:idx_processed"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime;index:idx_processed"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	Error       *string    `json:"error,omitempty" gorm:"type:text"`
}
