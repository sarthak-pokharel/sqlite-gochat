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
	ID                int64             `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID    int64             `json:"conversation_id" gorm:"not null;index:idx_conv_created"`
	PlatformMessageID *string           `json:"platform_message_id,omitempty" gorm:"index"`
	SenderType        MessageSenderType `json:"sender_type" gorm:"not null"`
	SenderID          *int64            `json:"sender_id,omitempty"`
	Content           string            `json:"content" gorm:"not null;type:text"`
	MessageType       MessageType       `json:"message_type" gorm:"default:text"`
	MediaURL          *string           `json:"media_url,omitempty" gorm:"type:text"`
	Direction         MessageDirection  `json:"direction" gorm:"not null"`
	Status            MessageStatus     `json:"status" gorm:"default:received"`
	CreatedAt         time.Time         `json:"created_at" gorm:"autoCreateTime;index:idx_created,idx_conv_created"`
	DeliveredAt       *time.Time        `json:"delivered_at,omitempty"`
	ReadAt            *time.Time        `json:"read_at,omitempty"`
	Metadata          *string           `json:"metadata,omitempty" gorm:"type:text"`
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
