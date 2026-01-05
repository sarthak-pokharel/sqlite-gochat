package models

import "time"

type Organization struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;not null" validate:"required,lowercase,alphanum"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	Metadata  *string   `json:"metadata,omitempty" gorm:"type:text"`
}

type CreateOrganizationRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=100"`
	Slug     string  `json:"slug" validate:"required,lowercase,alphanum"`
	Metadata *string `json:"metadata,omitempty"`
}

type UpdateOrganizationRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	IsActive *bool   `json:"is_active,omitempty"`
	Metadata *string `json:"metadata,omitempty"`
}
