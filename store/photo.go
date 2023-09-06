package store

import (
	_ "github.com/lib/pq"

	"github.com/inokone/photostorage/photo"
)

type Store interface {
	Store(ownerId string, photos ...photo.Photo) error

	Retrieve(ownerId string, filenames ...string) ([]photo.Photo, error)

	List(ownerId string) ([]photo.Descriptor, error)
}

type PostgresPhotoStore struct {
	url   string
	image ImageStore
}

func (p PostgresPhotoStore) Store(ownerId string, photos ...photo.Photo) {
	// TODO: fill
}
