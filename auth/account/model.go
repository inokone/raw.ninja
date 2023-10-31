package account

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Account is a struct to store state for authentication
type Account struct {
	UserID             uuid.UUID `gorm:"type:uuid;primary_key"`
	FailedLoginCounter int
	FailedLoginLock    time.Time
	LastFailedLogin    time.Time
	ConfirmationToken  string `gorm:"type:varchar(100);index"`
	ConfirmationTTL    time.Time
	Confirmed          bool
	RecoveryToken      string `gorm:"type:varchar(100);index"`
	RecoveryTTL        time.Time
	LastRecovery       time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt
}

// ConfirmationResend is a struct for the message body of REST endpoint e-mail confirmation resend
type ConfirmationResend struct {
	Email string `json:"email" binding:"required,email"`
}

// Recovery is a struct for the message body of REST endpoint password reset request
type Recovery struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordReset is a struct for the message body of REST endpoint password reset request
type PasswordReset struct {
	Token    string `json:"token" binding:"required,uuid"`
	Password string `json:"password" binding:"required"`
}

// PasswordChange is a struct for the message body of REST endpoint password change
type PasswordChange struct {
	New string `json:"new" binding:"required"`
	Old string `json:"old" binding:"required"`
}
