package auth

import (
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func (s *Store) New(db *gorm.DB) {
	s.db = db
}

func (s Store) Store(user User) error {
	result := s.db.Create(&user)
	return result.Error
}

func (s Store) User(email string) (User, error) {
	var user User
	result := s.db.Where(&User{Email: email}).First(&user)
	return user, result.Error
}

func (s Store) Delete(email string) error {
	var user User
	result := s.db.Where(&User{Email: email}).Delete(&user)
	return result.Error
}
