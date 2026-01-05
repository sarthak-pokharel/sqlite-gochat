package models

import "time"

type ConversationStatus string

const (
	ConversationStatusOpen     ConversationStatus = "open"
	ConversationStatusPending  ConversationStatus = "pending"
	ConversationStatusResolved ConversationStatus = "resolved"
	ConversationStatusClosed   ConversationStatus = "closed"
)

type ConversationPriority string

const (
	PriorityLow    ConversationPriority = "low"
	PriorityNormal ConversationPriority = "normal"
	PriorityHigh   ConversationPriority = "high"
	PriorityUrgent ConversationPriority = "urgent"
)

type Conversation struct {
	ID                   int64                `json:"id" gorm:"primaryKey;autoIncrement"`
	ChannelID            int64                `json:"channel_id" gorm:"not null;index"`
	ExternalUserID       int64                `json:"external_user_id" gorm:"not null;index"`
	AssignedToExternalID *string              `json:"assigned_to_external_id,omitempty" gorm:"index"`
	Status               ConversationStatus   `json:"status" gorm:"default:open;index"`
	Priority             ConversationPriority `json:"priority" gorm:"default:normal"`
	Subject              *string              `json:"subject,omitempty"`
	FirstMessageAt       *time.Time           `json:"first_message_at,omitempty"`
	LastMessageAt        *time.Time           `json:"last_message_at,omitempty"`
	ResolvedAt           *time.Time           `json:"resolved_at,omitempty"`
	CreatedAt            time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time            `json:"updated_at" gorm:"autoUpdateTime;index"`
	Metadata             *string              `json:"metadata,omitempty" gorm:"type:text"`
}

type CreateConversationRequest struct {
	ChannelID      int64                `validate:"required,gt=0"`
	ExternalUserID int64                `validate:"required,gt=0"`
	Subject        *string              `validate:"omitempty,max=200"`
	Priority       ConversationPriority `validate:"omitempty,oneof=low normal high urgent"`
}

type UpdateConversationRequest struct {
	AssignedToExternalID *string               `json:"assigned_to_external_id,omitempty"`
	Status               *ConversationStatus   `json:"status,omitempty" validate:"omitempty,oneof=open pending resolved closed"`
	Priority             *ConversationPriority `json:"priority,omitempty" validate:"omitempty,oneof=low normal high urgent"`
	Subject              *string               `json:"subject,omitempty" validate:"omitempty,max=200"`
	Metadata             *string               `json:"metadata,omitempty"`
}
