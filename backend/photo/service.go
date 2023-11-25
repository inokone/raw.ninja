package photo

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/image/importer"
	"github.com/rs/zerolog/log"
)

type service struct {
	photos Storer
	config common.ImageStoreConfig
}

func newService(photos Storer, config common.ImageStoreConfig) *service {
	return &service{
		photos: photos,
		config: config,
	}
}

type uploadResult struct {
	id  uuid.UUID
	err error
}

func (s service) upload(usr *user.User, file *multipart.FileHeader, ch chan uploadResult, wg *sync.WaitGroup) {
	var (
		err error
		mp  multipart.File
		raw []byte
		id  uuid.UUID
	)

	defer wg.Done()
	mp, err = file.Open()
	if err != nil {
		ch <- uploadResult{uuid.UUID{}, err}
	}
	defer closeRequestFile(mp)
	raw, err = io.ReadAll(mp)
	if err != nil {
		ch <- uploadResult{uuid.UUID{}, err}
	}
	id, err = s.uploadBinary(usr, raw, file.Filename)
	ch <- uploadResult{id, err}
}

func (s service) uploadBinary(usr *user.User, raw []byte, filename string) (uuid.UUID, error) {
	start := time.Now()
	var (
		target        *Photo
		id            uuid.UUID
		quotaExceeded bool
		err           error
	)
	target, err = createPhoto(
		*usr,
		filepath.Base(filename),
		filepath.Ext(filename)[1:],
		raw,
	)
	if err != nil {
		log.Err(err).Msg("Failed to create photo entity!")
		return uuid.UUID{}, fmt.Errorf("Uploaded file format is not supported! Cause: %v", err)
	}
	quotaExceeded, err = s.exceededUserQuota(usr, target.Desc.Metadata.DataSize)
	if quotaExceeded || err != nil {
		return uuid.UUID{}, errors.New("you can not upload files, you have reached your quota")
	}
	quotaExceeded, err = s.exceededGlobalQuota(target.Desc.Metadata.DataSize)
	if quotaExceeded || err != nil {
		log.Error().Msg("Global quota exceeded!")
		return uuid.UUID{}, errors.New("you can not upload files, please contact an administrator")
	}
	id, err = s.photos.Store(target)
	if err != nil {
		log.Err(err).Msg("Failed to store photo!")
		return uuid.UUID{}, errors.New("uploaded file could not be stored")
	}
	log.Debug().Str("file", filename).Dur("elapsed", time.Since(start)).Msg("photo stored")
	return id, err
}

func (s service) exceededGlobalQuota(fileSize int64) (bool, error) {
	var (
		quota int64
		stats Stats
		err   error
	)

	quota = s.config.Quota
	if quota <= 0 {
		return false, nil
	}

	stats, err = s.photos.Stats()
	if err != nil {
		return false, err
	}
	return stats.UsedSpace+fileSize > quota, nil
}

func (s service) exceededUserQuota(usr *user.User, fileSize int64) (bool, error) {
	var (
		stats UserStats
		err   error
	)

	if usr.Role.Quota <= 0 {
		return false, nil
	}

	stats, err = s.photos.UserStats(usr.ID.String())
	if err != nil {
		return false, err
	}
	return stats.UsedSpace+fileSize > usr.Role.Quota, nil
}

func closeRequestFile(mp multipart.File) {
	mp.Close()
}

func createPhoto(user user.User, filename, extension string, raw []byte) (*Photo, error) {
	i := importer.NewImporter(string(descriptor.ParseFormat(extension)))
	thumbnail, err := i.Thumbnail(raw)
	if err != nil {
		return nil, err
	}
	metadata, err := i.Describe(raw)
	if err != nil {
		return nil, err
	}
	res := &Photo{
		Desc: descriptor.Descriptor{
			FileName:  filename,
			Format:    descriptor.ParseFormat(extension),
			Uploaded:  time.Now(),
			Thumbnail: thumbnail,
			Metadata:  *metadata,
		},
		User:      user,
		Raw:       raw,
		UsedSpace: len(raw) + len(thumbnail),
	}
	return res, nil
}
