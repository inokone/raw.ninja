package photo

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/image/importer"
	"github.com/inokone/photostorage/photo/descriptor"
	"github.com/rs/zerolog/log"
)

// UploadService is a service entity handling photo uploads
type UploadService struct {
	photos Storer
	images image.Storer
	config common.ImageStoreConfig
}

// NewUploadService creates an `UploadService` instance based on storers and configuration
func NewUploadService(photos Storer, images image.Storer, config common.ImageStoreConfig) *UploadService {
	return &UploadService{
		photos: photos,
		images: images,
		config: config,
	}
}

// UploadResult is a struct to store result of a single photo's upload
type UploadResult struct {
	ID  uuid.UUID
	Err error
}

// Upload is a method og `UploadService`, capable of uploading a single file. The method is concurrency safe, target
// is parallelization when multiple files are uploaded. Results of the upload is added to the channel uploadResult.
func (s UploadService) Upload(usr *user.User, file *multipart.FileHeader, ch chan UploadResult, wg *sync.WaitGroup) {
	var (
		err error
		mp  multipart.File
		raw []byte
		id  uuid.UUID
	)

	defer wg.Done()
	mp, err = file.Open()
	if err != nil {
		ch <- UploadResult{uuid.UUID{}, err}
	}
	defer closeRequestFile(mp)
	raw, err = io.ReadAll(mp)
	if err != nil {
		ch <- UploadResult{uuid.UUID{}, err}
	}
	id, err = s.uploadBinary(usr, raw, file.Filename)
	ch <- UploadResult{id, err}
}

func (s UploadService) uploadBinary(usr *user.User, raw []byte, filename string) (uuid.UUID, error) {
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
	err = s.images.Store(target.ID.String(), target.Raw, target.Thumbnail)
	if err != nil {
		log.Err(err).Msg("Failed to store photo!")
		return uuid.UUID{}, errors.New("uploaded file could not be stored")
	}
	log.Debug().Str("file", filename).Dur("elapsed", time.Since(start)).Msg("photo stored")
	return id, err
}

func (s UploadService) exceededGlobalQuota(fileSize int64) (bool, error) {
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

func (s UploadService) exceededUserQuota(usr *user.User, fileSize int64) (bool, error) {
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
			FileName: filename,
			Format:   descriptor.ParseFormat(extension),
			Uploaded: time.Now(),
			Metadata: *metadata,
		},
		User:      user,
		Raw:       raw,
		Thumbnail: thumbnail,
		UsedSpace: len(raw) + len(thumbnail),
	}
	return res, nil
}

// LoadService is a service for retrieving raw and thumbnail files and links
type LoadService struct {
	photos Storer
	images image.Storer
	cfg    common.ImageStoreConfig
}

// NewLoadService creates a `LoadService` instance based on the storers and configuration.
func NewLoadService(photos Storer, images image.Storer, cfg common.ImageStoreConfig) *LoadService {
	return &LoadService{
		photos: photos,
		images: images,
		cfg:    cfg,
	}
}

// AsResponse transforms a Photo array to Response array
func (s LoadService) AsResponse(result []Photo, baseURL string) ([]Response, error) {
	var (
		err  error
		imgs []Response
	)
	imgs = make([]Response, len(result))
	for i, photo := range result {
		imgs[i] = photo.AsResp()
		if err = s.decorateWithRequest(&imgs[i], baseURL+imgs[i].ID); err != nil {
			log.Err(err).Str("photo_id", imgs[i].ID).Msg("Failed to generate presigned raw.")
			return nil, err
		}
	}
	return imgs, nil
}

// ThumbnailURL generates presigned URL for a thumbnail
func (s LoadService) ThumbnailURL(photoID uuid.UUID, baseURL string) (*image.PresignedRequest, error) {
	if s.cfg.UsePresigned {
		return s.images.PresignThumbnail(photoID.String())
	}
	return presign(baseURL + photoID.String() + "/thumbnail"), nil
}

func (s LoadService) decorateWithRequest(photo *Response, baseURL string) error {
	var (
		id  = photo.ID
		err error
	)
	if s.cfg.UsePresigned {
		photo.Raw, err = s.images.PresignImage(id)
		if err != nil {
			return err
		}
		photo.Thumbnail, err = s.images.PresignThumbnail(id)
		if err != nil {
			return err
		}
	} else {
		photo.Raw = presign(baseURL + "/raw")
		photo.Thumbnail = presign(baseURL + "/thumbnail")
	}
	return nil
}

func presign(URL string) *image.PresignedRequest {
	return &image.PresignedRequest{
		URL:    URL,
		Method: "GET",
		Header: http.Header{}, // TODO: this is not really presigned, just a request. should rename it
		Mode:   "cors",
	}
}
