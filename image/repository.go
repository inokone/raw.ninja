package image

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	Create(id string, raw []byte, thumbnail []byte) error

	Thumbnail(id string) ([]byte, error)

	Image(id string) ([]byte, error)

	Delete(id string) error

	UsedSpace(ids []string) (int64, error)
}

const (
	photoFolder   = "photos"
	rawName       = "raw"
	thumbnailName = "thumbnail.jpg"
)

type LocalRepository struct {
	Path string
}

func NewLocalStore(path string) (*LocalRepository, error) {
	fs := new(LocalRepository)
	fs.Path = filepath.Join(path, photoFolder)
	if err := os.MkdirAll(fs.Path, os.ModePerm); err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *LocalRepository) Create(id string, raw []byte, thumbnail []byte) error {
	path := filepath.Join(s.Path, photoFolder, id)
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
	return os.WriteFile(path, content, 0755)
}

func (s *LocalRepository) Delete(id string) error {
	path := filepath.Join(s.Path, photoFolder, id)
	return os.RemoveAll(path)
}

func (s *LocalRepository) Image(id string) ([]byte, error) {
	path := filepath.Join(s.Path, photoFolder, id, rawName)
	return os.ReadFile(path)
}

func (s *LocalRepository) Thumbnail(id string) ([]byte, error) {
	path := filepath.Join(s.Path, photoFolder, id, thumbnailName)
	thumbnail, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return thumbnail, nil
}

func (s *LocalRepository) UsedSpace(ids []string) (int64, error) {
	var sum int64 = 0
	for _, id := range ids {
		path := filepath.Join(s.Path, photoFolder, id)
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

// TODO: implement AWS S3 storage and try GCP storage
