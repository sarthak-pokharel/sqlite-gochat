package models

import "time"


type ExternalUser struct {
	ID               int64      `json:"id" db:"id"`
	ChannelID        int64      `json:"channel_id" db:"channel_id"`
	PlatformUserID   string     `json:"platform_user_id" db:"platform_user_id"`
	PlatformUsername *string    `json:"platform_username,omitempty" db:"platform_username"`
	DisplayName      *string    `json:"display_name,omitempty" db:"display_name"`
	PhoneNumber      *string    `json:"phone_number,omitempty" db:"phone_number"`
	Email            *string    `json:"email,omitempty" db:"email"`
	AvatarURL        *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Metadata         *string    `json:"metadata,omitempty" db:"metadata"`
	FirstSeenAt      time.Time  `json:"first_seen_at" db:"first_seen_at"`
	LastSeenAt       *time.Time `json:"last_seen_at,omitempty" db:"last_seen_at"`
	IsBlocked        bool       `json:"is_blocked" db:"is_blocked"`
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
