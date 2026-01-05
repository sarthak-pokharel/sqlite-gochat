package models

import "time"

type ExternalUser struct {
	ID               int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	ChannelID        int64      `json:"channel_id" gorm:"not null;index"`
	PlatformUserID   string     `json:"platform_user_id" gorm:"not null;index"`
	PlatformUsername *string    `json:"platform_username,omitempty"`
	DisplayName      *string    `json:"display_name,omitempty"`
	PhoneNumber      *string    `json:"phone_number,omitempty" gorm:"index"`
	Email            *string    `json:"email,omitempty"`
	AvatarURL        *string    `json:"avatar_url,omitempty" gorm:"type:text"`
	Metadata         *string    `json:"metadata,omitempty" gorm:"type:text"`
	FirstSeenAt      time.Time  `json:"first_seen_at" gorm:"autoCreateTime"`
	LastSeenAt       *time.Time `json:"last_seen_at,omitempty"`
	IsBlocked        bool       `json:"is_blocked" gorm:"default:false"`
}

type CreateExternalUserRequest struct {
	ChannelID        int64  `validate:"required,gt=0"`
	PlatformUserID   string `validate:"required"`
	PlatformUsername *string
	DisplayName      *string
	PhoneNumber      *string
	Email            *string `validate:"omitempty,email"`
	AvatarURL        *string `validate:"omitempty,url"`
	Metadata         *string
}

type UpdateExternalUserRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	AvatarURL   *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
	Metadata    *string `json:"metadata,omitempty"`
	IsBlocked   *bool   `json:"is_blocked,omitempty"`
}
