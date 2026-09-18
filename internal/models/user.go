package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "SUPER_ADMIN"
	RoleAdmin      UserRole = "ADMIN"
	RoleUser       UserRole = "USER"
)

type UserStatus string

const (
	StatusActive UserStatus = "active"
	StatusDelete UserStatus = "delete"
)

type Authentication struct {
	IsResetPassword bool       `json:"isResetPassword"`
	OneTimeCode     *int       `json:"oneTimeCode"`
	ExpireAt        *time.Time `json:"expireAt"`
}

type User struct {
	ID             string         `json:"id" gorm:"primaryKey;size:64"`
	Name           string         `json:"name" gorm:"not null"`
	Role           UserRole       `json:"role" gorm:"type:varchar(32);default:'USER'"`
	Contact        string         `json:"contact,omitempty"`
	Email          string         `json:"email" gorm:"uniqueIndex;not null"`
	Password       string         `json:"-" gorm:"not null"`
	Image          string         `json:"image,omitempty"`
	Status         UserStatus     `json:"status" gorm:"type:varchar(16);default:'active'"`
	Verified       bool           `json:"verified" gorm:"default:false"`
	Authentication datatypes.JSON `json:"-" gorm:"type:jsonb"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = "user-" + uuid.NewString()[:8]
	}
	if u.Role == "" {
		u.Role = RoleUser
	}
	if u.Status == "" {
		u.Status = StatusActive
	}
	return nil
}

type ResetToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	UserID    string    `json:"userId" gorm:"index;not null"`
	CreatedAt time.Time `json:"createdAt"`
}

func (ResetToken) TableName() string {
	return "reset_tokens"
}
