package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthenticationState is a struct to store state for authentication
type AuthenticationState struct {
	UserID                uuid.UUID `gorm:"type:uuid;primary_key"`
	FailedLoginCounter    int
	FailedLoginLock       time.Time
	LastFailedLogin       time.Time
	EmailConfirmationHash string `gorm:"type:varchar(100);uniqueIndex"`
	EmailConfirmationTTL  time.Time
	EmailConfirmed        bool
	PasswordResetHash     string `gorm:"type:varchar(100);uniqueIndex"`
	PasswordResetTTL      time.Time
	LastPasswordReset     time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt
}

// ConfirmationResend is a struct for the message body of REST endpoint e-mail confirmation resend
type ConfirmationResend struct {
	Email string `json:"email"`
}
