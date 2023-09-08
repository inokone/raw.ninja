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

func (s *Store) Store(photo Photo) {
	s.db.Save(&photo)
	s.is.Store(photo.ID.String(), photo.Raw, photo.Desc.Thumbnail)
}
