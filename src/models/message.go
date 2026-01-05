package models

import "time"


type MessageSenderType string

const (
	SenderExternal MessageSenderType = "external"
	SenderInternal MessageSenderType = "internal"
	SenderSystem   MessageSenderType = "system"
)


type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeFile     MessageType = "file"
	MessageTypeLocation MessageType = "location"
	MessageTypeContact  MessageType = "contact"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeSystem   MessageType = "system"
)


type MessageDirection string

const (
	DirectionInbound  MessageDirection = "inbound"
	DirectionOutbound MessageDirection = "outbound"
)


type MessageStatus string

const (
	MessageStatusReceived  MessageStatus = "received"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)


type Message struct {
	ID                int64             `json:"id" db:"id"`
	ConversationID    int64             `json:"conversation_id" db:"conversation_id"`
	PlatformMessageID *string           `json:"platform_message_id,omitempty" db:"platform_message_id"`
	SenderType        MessageSenderType `json:"sender_type" db:"sender_type"`
	SenderID          *int64            `json:"sender_id,omitempty" db:"sender_id"`
	Content           string            `json:"content" db:"content"`
	MessageType       MessageType       `json:"message_type" db:"message_type"`
	MediaURL          *string           `json:"media_url,omitempty" db:"media_url"`
	Direction         MessageDirection  `json:"direction" db:"direction"`
	Status            MessageStatus     `json:"status" db:"status"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	DeliveredAt       *time.Time        `json:"delivered_at,omitempty" db:"delivered_at"`
	ReadAt            *time.Time        `json:"read_at,omitempty" db:"read_at"`
	Metadata          *string           `json:"metadata,omitempty" db:"metadata"`
}


type CreateMessageRequest struct {
	ConversationID int64       `json:"conversation_id" validate:"required,gt=0"`
	Content        string      `json:"content" validate:"required,min=1"`
	MessageType    MessageType `json:"message_type" validate:"omitempty,oneof=text image video audio file location contact sticker"`
	MediaURL       *string     `json:"media_url,omitempty" validate:"omitempty,url"`
	SenderID       *int64      `json:"sender_id,omitempty"`
	Metadata       *string     `json:"metadata,omitempty"`
}


type InboundMessageRequest struct {
	ChannelID         int64       `validate:"required,gt=0"`
	PlatformMessageID *string     `validate:"required"`
	PlatformUserID    string      `validate:"required"`
	Content           string      `validate:"required,min=1"`
	MessageType       MessageType `validate:"required,oneof=text image video audio file location contact sticker"`
	MediaURL          *string     `validate:"omitempty,url"`
	Metadata          *string
}


type MessageListQuery struct {
	ConversationID int64  `query:"conversation_id" validate:"required,gt=0"`
	Limit          int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset         int    `query:"offset" validate:"omitempty,min=0"`
	Before         *int64 `query:"before"` 
}
