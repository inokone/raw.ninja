package common

import (
	"golang.org/x/crypto/bcrypt"

	"encoding/base64"
)

type User struct {
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
