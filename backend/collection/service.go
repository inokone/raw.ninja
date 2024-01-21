package collection

import (
	"fmt"
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

// InvalidPhotoID is an error for invalid, malformed of non-existing IDs of photos
type InvalidPhotoID struct {
	ID string
}

// Error is the string representation of an `InvalidPhotoID` error
func (e InvalidPhotoID) Error() string { return fmt.Sprintf("invalid photo ID [%v]", e.ID) }

// CreateUpload is a method of `Service` creating a persisted upload type collection
func (s Service) CreateUpload(usr user.User, photoIDs []uuid.UUID) (*Collection, error) {
	return s.createCollection(usr, time.Now().Local().Format("2006-01-02 15:04"), Upload, nil, photoIDs)
}

func (s Service) createCollection(usr user.User, name string, ct Type, tags []string, photoIDs []uuid.UUID) (*Collection, error) {
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
		Tags:        tags,
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
func (s Service) CreateAlbum(usr user.User, name string, tags []string, photoIDs []uuid.UUID) (*Collection, error) {
	return s.createCollection(usr, name, Album, tags, photoIDs)
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

// Update manages changes of a collection. Current supported fields are name, tags and photos
func (s Service) Update(cl *Collection, cr Resp) (*Collection, error) {
	var err error
	if cr.Tags != nil {
		cl.Tags = cr.Tags
	}
	if len(cr.Name) > 0 {
		cl.Name = cr.Name
	}
	// update the photos
	if cr.Photos != nil {
		cl.Photos, err = s.updatePhotos(cr.Photos)
		if err != nil {
			return nil, err
		}
	}
	// check whether the thumbnail of the album is still associated to the album
	if len(cl.Photos) > 0 {
		cl.ThumbnailID = &cl.Photos[0].ID
	} else {
		cl.ThumbnailID = nil
	}
	err = s.s.Update(cl)
	return cl, err
}

func (s Service) updatePhotos(updated []photo.Response) ([]photo.Photo, error) {
	var (
		err error
		id  uuid.UUID
		ids []uuid.UUID
	)

	ids = make([]uuid.UUID, len(updated))
	for i, rl := range updated {
		id, err = uuid.Parse(rl.ID)
		if err != nil {
			return nil, InvalidPhotoID{rl.ID}
		}
		ids[i] = id
	}

	return createPhotos(ids), nil
}

// SearchAlbums searchews for album collections of a user based on a query string
func (s Service) SearchAlbums(usrID uuid.UUID, query string) ([]ListItem, error) {
	return s.s.Search(usrID, Album, query)
}

// SearchUploads searchews for upload collections of a user based on a query string
func (s Service) SearchUploads(usrID uuid.UUID, query string) ([]ListItem, error) {
	return s.s.Search(usrID, Upload, query)
}
