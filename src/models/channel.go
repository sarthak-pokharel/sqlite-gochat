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
	ID                int64         `json:"id" db:"id"`
	OrganizationID    int64         `json:"organization_id" db:"organization_id"`
	Platform          Platform      `json:"platform" db:"platform"`
	Name              string        `json:"name" db:"name"`
	AccountIdentifier string        `json:"account_identifier" db:"account_identifier"`
	Status            ChannelStatus `json:"status" db:"status"`
	WebhookSecret     *string       `json:"webhook_secret,omitempty" db:"webhook_secret"`
	AccessToken       *string       `json:"-" db:"access_token"`
	Config            *string       `json:"config,omitempty" db:"config"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`
	LastMessageAt     *time.Time    `json:"last_message_at,omitempty" db:"last_message_at"`
	IsActive          bool          `json:"is_active" db:"is_active"`
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
