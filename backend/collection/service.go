package collection

import (
	"time"

	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/photo"
)

// Service for handling collection related business logic
type Service struct {
	s Storer
}

// NewService is a function that creates a new instance of the service
func NewService(s Storer) *Service {
	return &Service{
		s: s,
	}
}

// CreateUpload is a method of `Service` creating a persisted upload type collection
func (s Service) CreateUpload(usr user.User, photoIDs []uuid.UUID) (*Collection, error) {
	return s.createCollection(usr, time.Now().Local().Format("2006-01-02 15:04:05"), Upload, photoIDs)
}

func (s Service) createCollection(usr user.User, name string, ct Type, photoIDs []uuid.UUID) (*Collection, error) {
	var (
		u         *Collection
		err       error
		photos    []photo.Photo
		thumbnail *uuid.UUID
	)
	photos = createPhotos(photoIDs)
	if len(photoIDs) > 0 {
		thumbnail = &photoIDs[0]
	}
	u = &Collection{
		Type:        ct,
		User:        usr,
		Photos:      photos,
		Name:        name,
		ThumbnailID: thumbnail,
	}
	if err = s.s.Store(u); err != nil {
		return nil, err
	}
	return u, nil
}

// CreateAlbum is a method of `Service` creating a persisted album type collection
func (s Service) CreateAlbum(usr user.User, name string, photoIDs []uuid.UUID) (*Collection, error) {
	return s.createCollection(usr, name, Album, photoIDs)
}

func createPhotos(photoIDs []uuid.UUID) []photo.Photo {
	var res []photo.Photo = make([]photo.Photo, len(photoIDs))
	for i, id := range photoIDs {
		res[i] = photo.Photo{
			ID: id,
		}
	}
	return res
}

// SetProperties is a method of `Service` for updating name and tags of a persisted collection
func (s Service) SetProperties(collectionID uuid.UUID, name string, tags []string) error {
	var (
		c   *Collection
		err error
	)
	c, err = s.s.ByID(collectionID)
	if err != nil {
		return err
	}

	c.Tags = tags
	c.Name = name
	return s.s.Update(c)
}
