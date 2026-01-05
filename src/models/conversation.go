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
	ID                   int64                `json:"id" db:"id"`
	ChannelID            int64                `json:"channel_id" db:"channel_id"`
	ExternalUserID       int64                `json:"external_user_id" db:"external_user_id"`
	AssignedToExternalID *string              `json:"assigned_to_external_id,omitempty" db:"assigned_to_external_id"`
	Status               ConversationStatus   `json:"status" db:"status"`
	Priority             ConversationPriority `json:"priority" db:"priority"`
	Subject              *string              `json:"subject,omitempty" db:"subject"`
	FirstMessageAt       *time.Time           `json:"first_message_at,omitempty" db:"first_message_at"`
	LastMessageAt        *time.Time           `json:"last_message_at,omitempty" db:"last_message_at"`
	ResolvedAt           *time.Time           `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt            time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at" db:"updated_at"`
	Metadata             *string              `json:"metadata,omitempty" db:"metadata"`
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
