package model

import "time"

type Role string

const (
	RoleUser   Role = "user"
	RoleAdmin  Role = "admin"
	RoleSystem Role = "system"
)

type User struct {
	UserID            string                 `json:"user_id" db:"user_id"` // UUID
	Email             string                 `json:"email" db:"email"`
	EncryptedPassword string                 `json:"-" db:"encrypted_password"` // omit from JSON responses
	FullName          string                 `json:"full_name,omitempty" db:"full_name"`
	Role              Role                   `json:"role" db:"role"`
	IsActive          bool                   `json:"is_active" db:"is_active"`
	Timezone          *string                `json:"timezone,omitempty" db:"timezone"`
	Metadata          map[string]interface{} `json:"metadata,omitempty" db:"metadata"` // JSONB
	LastLogin         *time.Time             `json:"last_login,omitempty" db:"last_login"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`
}

type UserToken struct {
	TokenID        string    `json:"token_id" db:"token_id"` // UUID
	UserID         string    `json:"user_id" db:"user_id"`   // UUID
	EncryptedToken string    `json:"-" db:"encrypted_token"` // omit from JSON responses
	IssuedAt       time.Time `json:"issued_at" db:"issued_at"`
	ExpiresAt      time.Time `json:"expires_at" db:"expires_at"`
	Revoked        bool      `json:"revoked" db:"revoked"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
