package image

import (
	"image"
	"log"
	"os"
	"path/filepath"
)

type Store interface {
	Store(id string, raw []byte, thumbnail image.Image) error

	Thumbnail(id string) (image.Image, error)

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
	err := os.MkdirAll(fs.Path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *LocalStore) Store(id string, raw []byte, thumbnail image.Image) error {
	path := filepath.Join(s.Path, photoFolder, id)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalf("Can not create path [%v] for image [%v]", path, id)
		return err
	}
	err = write(filepath.Join(path, rawName), raw)
	if err != nil {
		log.Fatalf("Can not write raw to path [%v] for image [%v]", path, id)
		return err
	}
	b, err := ExportJpeg(thumbnail)
	if err != nil {
		log.Fatalf("Can not export thumbnail to JPG for image [%v]", id)
		return err
	}
	err = write(filepath.Join(path, thumbnailName), b)
	if err != nil {
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
	err := os.RemoveAll(path)
	return err
}

func (s *LocalStore) Image(id string) ([]byte, error) {
	path := filepath.Join(s.Path, photoFolder, id, rawName)
	return os.ReadFile(path)
}

func (s *LocalStore) Thumbnail(id string) (image.Image, error) {
	path := filepath.Join(s.Path, photoFolder, id, thumbnailName)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ImportJpeg(raw)
}

// TODO: implement AWS S3 storage and try GCP storage
