package image

import (
	"log"
	"os"
	"path/filepath"
)

type Store interface {
	Store(id string, raw []byte, thumbnail []byte) error

	Thumbnail(id string) ([]byte, error)

	Image(id string) ([]byte, error)

	Delete(id string) error
}

const (
	photoFolder   = "photos"
	rawName       = "raw"
	thumbnailName = "thumbnail.jpg"
)

type LocalStore struct {
	Path string
}

func NewLocalStore(path string) (*LocalStore, error) {
	fs := new(LocalStore)
	fs.Path = filepath.Join(path, photoFolder)
	if err := os.MkdirAll(fs.Path, os.ModePerm); err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *LocalStore) Store(id string, raw []byte, thumbnail []byte) error {
	path := filepath.Join(s.Path, photoFolder, id)
	var err error
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Can not create path [%v] for image [%v]", path, id)
		return err
	}
	if err = write(filepath.Join(path, rawName), raw); err != nil {
		log.Fatalf("Can not write raw to path [%v] for image [%v]", path, id)
		return err
	}
	if err = write(filepath.Join(path, thumbnailName), thumbnail); err != nil {
		log.Fatalf("Can not write thumbnail to path [%v] for image [%v]", path, id)
		return err
	}
	return nil
}

func write(path string, content []byte) error {
	return os.WriteFile(path, content, 0755)
}

func (s *LocalStore) Delete(id string) error {
	path := filepath.Join(s.Path, photoFolder, id)
	return os.RemoveAll(path)
}

func (s *LocalStore) Image(id string) ([]byte, error) {
	path := filepath.Join(s.Path, photoFolder, id, rawName)
	return os.ReadFile(path)
}

func (s *LocalStore) Thumbnail(id string) ([]byte, error) {
	path := filepath.Join(s.Path, photoFolder, id, thumbnailName)
	thumbnail, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return thumbnail, nil
}

// TODO: implement AWS S3 storage and try GCP storage
