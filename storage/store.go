package store

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"photo/photo"

	"github.com/google/uuid"
)

type Store interface {
	Store(userId string, photos ...photo.Photo) error

	Retrieve(userId string, filenames ...string) ([]photo.Photo, error)

	List(userId string) ([]photo.Descriptor, error)
}

func (s LocalStorage) generateFileName(prefix string) string {
	id := uuid.New()
	return prefix + id.String()
}

type ImageStorage interface {
	Store(id string, image image.Image, thumbnail image.Image) error

	Thumbnails(ids ...string) ([]image.Image, error)

	Images(ids ...string) ([]image.Image, error)
}

type LocalStorage struct {
	path string
}

func (s LocalStorage) New(path string) {
	s.path = filepath.Join(path, "photos")
	err := os.MkdirAll(s.path, os.ModePerm)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can not create application storage path [%v]", s.path))
		panic("Photo storage invalid. See logs for more details.")
	}
}

func (s LocalStorage) Store(id string, image image.Image, thumbnail image.Image) error {
	path := filepath.Join(s.path, id)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can not create path [%v] for user [%v]", s.path))
		return err
	}
	// TODO: finish up
	return nil
}

type S3Storage struct {
	account string
	bucket  string
	prefix  string
}

// TODO: implement S3 storage
