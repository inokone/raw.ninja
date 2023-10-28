package descriptor

import (
	"gorm.io/gorm"
)

// Storer is an interface for persisting `Descriptor` entities.
type Storer interface {
	Store(desc *Descriptor) error
	Delete(id string) error
}

// GORMStorer is an implementation of `Storer` using GORM library.
type GORMStorer struct {
	db *gorm.DB
}

// Store is a method of `GORMStorer` for persisting the `Descriptor` object provided as a parameter.
func (s *GORMStorer) Store(desc *Descriptor) error {
	result := s.db.Create(&desc)
	return result.Error
}

// Delete is a method of `GORMStorer` for deleting a `Descriptor` object from persistence specified by the ID provided as a parameter.
func (s *GORMStorer) Delete(id string) error {
	var desc Descriptor
	result := s.db.Where("id = ?", id).Delete(&desc)
	return result.Error
}
