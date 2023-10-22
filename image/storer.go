package image

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// Writer is an interface for changing images (RAW or processed).
type Writer interface {
	Store(id string, raw []byte, thumbnail []byte) error

	Delete(id string) error
}

// Loader is an interface for loading images (RAW or processed).
type Loader interface {
	LoadThumbnail(id string) ([]byte, error)

	LoadImage(id string) ([]byte, error)
}

// Storer is an interface for types that can store images (RAW or processed).
type Storer interface {
	Writer

	Loader

	UsedSpace(ids []string) (int64, error)
}

const (
	imageFolder   = "photos"
	rawName       = "raw"
	thumbnailName = "thumbnail.jpg"
)

// LocalStorer is an implementation of the Storer interface as pointer.
// that stores images on the local disk.
type LocalStorer struct {
	path string
}

// NewLocalStorer creates a new LocalStorer with the specified storage path.
func NewLocalStorer(path string) (*LocalStorer, error) {
	ip := filepath.Join(path, imageFolder)
	fs := LocalStorer{
		path: ip,
	}
	if err := os.MkdirAll(ip, os.ModePerm); err != nil {
		return &LocalStorer{}, err
	}
	return &fs, nil
}

// Store stores an image on the local disk.
func (s *LocalStorer) Store(id string, raw []byte, thumbnail []byte) error {
	path := filepath.Join(s.path, imageFolder, id)
	var err error
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		log.Error().Err(err).Str("path", path).Str("id", id).Msg("Failed to create path for image store.")
		return err
	}
	if err = write(filepath.Join(path, rawName), raw); err != nil {
		log.Error().Err(err).Str("path", path).Str("id", id).Msg("Failed to write raw")
		return err
	}
	if err = write(filepath.Join(path, thumbnailName), thumbnail); err != nil {
		log.Error().Err(err).Str("path", path).Str("id", id).Msg("Failed to write thumbnail")
		return err
	}
	return nil
}

func write(path string, content []byte) error {
	return os.WriteFile(path, content, 0o755)
}

// Delete deletes a image on the local disk.
func (s *LocalStorer) Delete(id string) error {
	path := filepath.Join(s.path, imageFolder, id)
	return os.RemoveAll(path)
}

// LoadImage loads a image specified by the id from the local disk.
func (s *LocalStorer) LoadImage(id string) ([]byte, error) {
	path := filepath.Join(s.path, imageFolder, id, rawName)
	return os.ReadFile(path)
}

// LoadThumbnail loads the thumbnail of the image specified by the id from the local disk.
func (s *LocalStorer) LoadThumbnail(id string) ([]byte, error) {
	path := filepath.Join(s.path, imageFolder, id, thumbnailName)
	thumbnail, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return thumbnail, nil
}

// UsedSpace calculates the disk space needed to store the images with the provided IDs
func (s *LocalStorer) UsedSpace(ids []string) (int64, error) {
	var sum int64
	for _, id := range ids {
		path := filepath.Join(s.path, imageFolder, id)
		c, err := getFolderSize(path)
		if err != nil {
			return 0, err
		}
		sum += c
	}
	return sum, nil
}

func getFolderSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return size, nil
}
