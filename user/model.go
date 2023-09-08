package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"encoding/base64"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email    string
	PassHash string
	Phone    string
}

func (u *User) New(email string, password string, phone string) error {
	u.Email = email
	u.Phone = phone
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassHash = base64.RawStdEncoding.EncodeToString(hash)
	return nil
}
