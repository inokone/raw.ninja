package image

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

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
	return os.ReadFile(path)
}

// SupportsPresign indicates whether the store supports presign
func (s *LocalStorer) SupportsPresign() bool {
	return false
}

// PresignThumbnail makes a presigned request that can be used to get a thumbnail.
func (s *LocalStorer) PresignThumbnail(id string) (*PresignedRequest, error) { // nolint:revive
	panic("Unsupported operation!")
}

// PresignImage makes a presigned request that can be used to get a raw image.
func (s *LocalStorer) PresignImage(id string) (*PresignedRequest, error) { // nolint:revive
	panic("Unsupported operation!")
}
