package auth

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PassHash  string    `gorm:"type:varchar(100)"`
	Phone     string    `gorm:"type:varchar(20)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func NewUser(email string, password string, phone string) (*User, error) {
	u := new(User)
	u.Email = email
	u.Phone = phone
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.PassHash = string(hash)
	return u, nil
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PassHash), []byte(password))
	print(err)
	return err == nil
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticatedUser struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	UserType     string    `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
}

type Registration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}
