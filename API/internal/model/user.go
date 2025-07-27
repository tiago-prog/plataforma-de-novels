package model

import (
	"time"
)

type Role string
type UserStatus string

const (
	RoleReader          Role       = "reader"
	RoleCreator         Role       = "creator"
	RoleModerator       Role       = "moderator"
	RoleAdmin           Role       = "admin"
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

type User struct {
	ID          int        `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Password    string     `json:"password"`
	Role        Role       `json:"role"`
	Avatar      string     `json:"avatar,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	Verified    bool       `json:"verified"`
	Status      UserStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
