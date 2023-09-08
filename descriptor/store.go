package descriptor

import (
	"gorm.io/gorm"
)

type DescriptorStore struct {
	db *gorm.DB
}

func (s *DescriptorStore) New(db *gorm.DB) {
	s.db = db
}

func (s DescriptorStore) Store(desc Descriptor) error {
	result := s.db.Create(&desc)
	return result.Error
}

func (s DescriptorStore) Retrieve(ids ...string) ([]Descriptor, error) {
	var descs []Descriptor
	result := s.db.Where("id IN ?", ids).Find(&descs)
	return descs, result.Error
}

func (s DescriptorStore) Delete(id string) error {
	var desc Descriptor
	result := s.db.Where("id = ?", id).Delete(&desc)
	return result.Error
}
