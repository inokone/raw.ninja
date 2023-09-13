package photo

import (
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/image"
)

type Store struct {
	db *gorm.DB
	is image.Store
}

func (s *Store) New(db *gorm.DB, is image.Store) {
	s.db = db
	s.is = is
}

func (s *Store) Store(photo Photo) error {
	s.db.Save(&photo)
	err := s.is.Store(photo.ID.String(), photo.Raw, photo.Desc.Thumbnail)
	return err
}

func (s *Store) Get(id string) (Photo, error) {
	var photo Photo
	result := s.db.First(&photo, id)
	return photo, result.Error
}

func (s *Store) List(userID string) ([]Photo, error) {
	var photos []Photo
	result := s.db.Where("user_id = ?", userID).Find(&photos)
	return photos, result.Error
}

func (s *Store) Raw(id string) ([]byte, error) {
	return s.is.Image(id)
}
