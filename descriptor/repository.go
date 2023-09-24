package descriptor

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (s *Repository) New(db *gorm.DB) {
	s.db = db
}

func (s Repository) Create(desc Descriptor) error {
	result := s.db.Create(&desc)
	return result.Error
}

func (s Repository) Delete(id string) error {
	var desc Descriptor
	result := s.db.Where("id = ?", id).Delete(&desc)
	return result.Error
}
