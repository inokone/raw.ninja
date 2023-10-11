package photo

import (
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/image"
)

type Repository struct {
	db *gorm.DB
	ir image.Repository
}

func (s *Repository) Create(photo Photo) (uuid.UUID, error) {
	s.db.Save(&photo)
	err := s.ir.Create(photo.ID.String(), photo.Raw, photo.Desc.Thumbnail)
	return photo.ID, err
}

func (s *Repository) Get(id string) (Photo, error) {
	var photo Photo
	result := s.db.Preload("Desc.Metadata").First(&photo, "id = ?", id)
	return photo, result.Error
}

func (s *Repository) All(userID string) ([]Photo, error) {
	var photos []Photo
	result := s.db.Preload("Desc.Metadata").Where("user_id = ?", userID).Find(&photos)
	for i := 0; i < len(photos); i++ {
		thumb, err := s.ir.Thumbnail(photos[i].ID.String())
		if err != nil {
			return nil, err
		}
		photos[i].Desc.Thumbnail = thumb
	}
	return photos, result.Error
}

func (s *Repository) Raw(id string) ([]byte, error) {
	return s.ir.Image(id)
}

func (s *Repository) Thumbnail(id string) ([]byte, error) {
	return s.ir.Thumbnail(id)
}
