package store

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

type ImageStore interface {
	Store(id string, raw []byte, thumbnail image.Image) error

	Thumbnails(ids ...string) ([]image.Image, error)

	Images(ids ...string) ([][]byte, error)

	Delete(id string) error
}

const (
	photoFolder   = "photos"
	rawName       = "raw"
	thumbnailName = "thumbnail.jpg"
)

type LocalStore struct {
	path string
}

func (s LocalStore) New(path string) {
	s.path = filepath.Join(path, photoFolder)
	err := os.MkdirAll(s.path, os.ModePerm)
	if err != nil {
		log.Fatalf("Can not create application storage path [%v]", s.path)
		panic("Photo storage invalid. See logs for more details.")
	}
}

func (s LocalStore) Store(id string, raw []byte, thumbnail image.Image) error {
	path := filepath.Join(s.path, photoFolder, id)
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
	b, err := export(thumbnail)
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

func export(image image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)
	return buf.Bytes(), err
}

func importTumb(b []byte) (image.Image, error) {
	return jpeg.Decode(bytes.NewReader(b))
}

func (s LocalStore) Delete(id string) error {
	path := filepath.Join(s.path, photoFolder, id)
	err := os.RemoveAll(path)
	return err
}

func (s LocalStore) Images(ids ...string) ([][]byte, error) {
	res := make([][]byte, len(ids))
	for i, id := range ids {
		path := filepath.Join(s.path, photoFolder, id, rawName)
		raw, err := os.ReadFile(path)
		res[i] = raw
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (s LocalStore) Thumbnails(ids ...string) ([]image.Image, error) {
	res := make([]image.Image, len(ids))
	for i, id := range ids {
		path := filepath.Join(s.path, photoFolder, id, thumbnailName)
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		res[i], err = importTumb(raw)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// TODO: implement AWS S3 storage and try GCP storage
