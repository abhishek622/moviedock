package model

import "time"

type Role string

const (
	RoleUser   Role = "user"
	RoleAdmin  Role = "admin"
	RoleSystem Role = "system"
)

type User struct {
	UserID       string                 `json:"user_id"`
	FullName     string                 `json:"full_name"`
	Email        string                 `json:"email"`
	PasswordHash string                 `json:"password_hash"`
	Role         Role                   `json:"role"`
	IsActive     bool                   `json:"is_active"`
	Timezone     string                 `json:"timezone"`
	Metadata     map[string]interface{} `json:"metadata"`
	LastLogin    time.Time              `json:"last_login"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
