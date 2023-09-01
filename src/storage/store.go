package store

import (
	"photo"
)

type Store interface {
	Store(userId string, photos ...photo.Photo) (string, error)

	Retrieve(userId string, filenames ...string) (photo.Photo, error)

	List(userId string) ([]photo.Descriptor, error)
}
