package models

import "time"

type Platform string

const (
	PlatformWhatsApp  Platform = "whatsapp"
	PlatformTelegram  Platform = "telegram"
	PlatformInstagram Platform = "instagram"
	PlatformFacebook  Platform = "facebook"
	PlatformSMS       Platform = "sms"
	PlatformEmail     Platform = "email"
	PlatformWeb       Platform = "web"
)

type ChannelStatus string

const (
	ChannelStatusActive   ChannelStatus = "active"
	ChannelStatusInactive ChannelStatus = "inactive"
	ChannelStatusError    ChannelStatus = "error"
	ChannelStatusPending  ChannelStatus = "pending"
)

type ChatChannel struct {
	ID                int64         `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationID    int64         `json:"organization_id" gorm:"not null;index"`
	Platform          Platform      `json:"platform" gorm:"not null;index"`
	Name              string        `json:"name" gorm:"not null;size:100"`
	AccountIdentifier string        `json:"account_identifier" gorm:"not null"`
	Status            ChannelStatus `json:"status" gorm:"default:active;index"`
	WebhookSecret     *string       `json:"webhook_secret,omitempty" gorm:"type:text"`
	AccessToken       *string       `json:"-" gorm:"type:text"`
	Config            *string       `json:"config,omitempty" gorm:"type:text"`
	CreatedAt         time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	LastMessageAt     *time.Time    `json:"last_message_at,omitempty"`
	IsActive          bool          `json:"is_active" gorm:"default:true"`
}

type CreateChannelRequest struct {
	OrganizationID    int64    `json:"organization_id" validate:"required,gt=0"`
	Platform          Platform `json:"platform" validate:"required,oneof=whatsapp telegram instagram facebook sms email web"`
	Name              string   `json:"name" validate:"required,min=2,max=100"`
	AccountIdentifier string   `json:"account_identifier" validate:"required"`
	WebhookSecret     *string  `json:"webhook_secret,omitempty"`
	AccessToken       *string  `json:"access_token,omitempty"`
	Config            *string  `json:"config,omitempty"`
}

type UpdateChannelRequest struct {
	Name          *string        `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Status        *ChannelStatus `json:"status,omitempty" validate:"omitempty,oneof=active inactive error pending"`
	WebhookSecret *string        `json:"webhook_secret,omitempty"`
	AccessToken   *string        `json:"access_token,omitempty"`
	Config        *string        `json:"config,omitempty"`
	IsActive      *bool          `json:"is_active,omitempty"`
}
