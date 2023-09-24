package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (s Repository) Create(user User) error {
	result := s.db.Create(&user)
	return result.Error
}

func (s Repository) ByEmail(email string) (User, error) {
	var user User
	result := s.db.Where(&User{Email: email}).First(&user)
	return user, result.Error
}

func (s Repository) ByID(id uuid.UUID) (User, error) {
	var user User
	result := s.db.Where(&User{ID: id}).First(&user)
	return user, result.Error
}

func (s Repository) Delete(email string) error {
	var user User
	result := s.db.Where(&User{Email: email}).Delete(&user)
	return result.Error
}
