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
func (s Service) CreateUpload(usr user.User, photoIDs []string) (*Collection, error) {
	var (
		u      Collection
		err    error
		photos []photo.Photo
	)
	photos, err = createPhotos(photoIDs)
	if err != nil {
		return nil, err
	}
	u = Collection{
		Type:   Upload,
		Owner:  usr,
		Photos: photos,
		Name:   time.Now().Local().Local().Format("YYYY-MM-DD"),
	}
	if err = s.s.Store(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateAlbum is a method of `Service` creating a persisted album type collection
func (s Service) CreateAlbum(usr user.User, name string, photoIDs []string) (*Collection, error) {
	var (
		u      Collection
		err    error
		photos []photo.Photo
	)
	photos, err = createPhotos(photoIDs)
	if err != nil {
		return nil, err
	}
	u = Collection{
		Type:   Album,
		Owner:  usr,
		Photos: photos,
		Name:   name,
	}
	if err = s.s.Store(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func createPhotos(photoIDs []string) ([]photo.Photo, error) {
	var (
		res []photo.Photo = make([]photo.Photo, len(photoIDs))
		err error
	)
	for i, id := range photoIDs {
		res[i] = photo.Photo{}
		res[i].ID, err = uuid.Parse(id)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
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
