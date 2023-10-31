package account

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Account is a struct to store state for authentication
type Account struct {
	UserID                uuid.UUID `gorm:"type:uuid;primary_key"`
	FailedLoginCounter    int
	FailedLoginLock       time.Time
	LastFailedLogin       time.Time
	EmailConfirmationHash string `gorm:"type:varchar(100);index"`
	EmailConfirmationTTL  time.Time
	EmailConfirmed        bool
	PasswordResetHash     string `gorm:"type:varchar(100);index"`
	PasswordResetTTL      time.Time
	LastPasswordReset     time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt
}

// ConfirmationResend is a struct for the message body of REST endpoint e-mail confirmation resend
type ConfirmationResend struct {
	Email string `json:"email" binding:"required,email"`
}

// RequestPwdReset is a struct for the message body of REST endpoint password reset request
type RequestPwdReset struct {
	Email string `json:"email" binding:"required,email"`
}
