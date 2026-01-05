package models

import "time"

type Organization struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Slug      string    `json:"slug" db:"slug" validate:"required,lowercase,alphanum"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	Metadata  *string   `json:"metadata,omitempty" db:"metadata"`
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
