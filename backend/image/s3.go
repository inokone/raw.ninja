package image

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

const (
	prefix    = "photos"
	rawFile   = "raw"
	thumbFile = "thumbnail.jpg"
	rawCT     = "application/octet-stream"
	thumbCT   = "image/jpeg"
)

// S3Storer is an implementation of the Storer interface as pointer that stores images on Amazon S3 buckets.
type S3Storer struct {
	rawBucket    string
	thumbBucket  string
	client       *s3.Client
	presign      *s3.PresignClient
	presignedTTL int64
}

// NewS3Storer creates a new LocalStorer with the specified storage path.
func NewS3Storer(rawBucket, thumbBucket, awsAPIKey, awsAPISecret string, presignedTTL int64) (*S3Storer, error) {
	var (
		c   aws.Config
		err error
		cl  *s3.Client
	)

	c, err = config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsAPIKey,
			awsAPISecret, "")),
		config.WithRegion("eu-central-1"),
	)
	if err != nil {
		return nil, err
	}

	cl = s3.NewFromConfig(c)

	return &S3Storer{
		rawBucket:    rawBucket,
		thumbBucket:  thumbBucket,
		client:       cl,
		presign:      s3.NewPresignClient(cl),
		presignedTTL: presignedTTL,
	}, nil
}

// Store stores an image on the local disk.
func (s *S3Storer) Store(id string, raw []byte, thumbnail []byte) error {
	path := filepath.Join(prefix, id)
	var err error
	if err = s.writeS3(s.rawBucket, filepath.Join(path, rawName), rawCT, raw); err != nil {
		log.Error().Err(err).Str("path", path).Str("id", id).Msg("Failed to write raw")
		return err
	}
	if err = s.writeS3(s.thumbBucket, filepath.Join(path, thumbnailName), thumbCT, thumbnail); err != nil {
		log.Error().Err(err).Str("path", path).Str("id", id).Msg("Failed to write thumbnail")
		return err
	}
	return nil
}

func (s *S3Storer) writeS3(bucket string, path string, contentType string, content []byte) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
	})
	return err
}

// Delete deletes an image from Amazon S3.
func (s *S3Storer) Delete(id string) error {
	var (
		path, rawFile, thumbFile string
		err                      error
	)

	path = filepath.Join(prefix, imageFolder, id)
	rawFile = filepath.Join(path, rawName)
	thumbFile = filepath.Join(path, thumbnailName)

	_, err = s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.rawBucket),
		Key:    aws.String(rawFile),
	})
	if err != nil {
		return err
	}

	_, err = s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.thumbBucket),
		Key:    aws.String(thumbFile),
	})
	return err
}

// LoadImage loads a image specified by the id from from Amazon S3.
func (s *S3Storer) LoadImage(id string) ([]byte, error) {
	log.Debug().Str("id", id).Msg("Collecting image")
	path := filepath.Join(prefix, id, rawName)
	return s.loadS3(s.rawBucket, path)
}

// LoadThumbnail loads the thumbnail of the image specified by the id from Amazon S3.
func (s *S3Storer) LoadThumbnail(id string) ([]byte, error) {
	log.Debug().Str("id", id).Msg("Collecting thumbnail")
	path := filepath.Join(prefix, id, thumbnailName)
	return s.loadS3(s.rawBucket, path)
}

func (s *S3Storer) loadS3(bucket string, key string) ([]byte, error) {
	var (
		result *s3.GetObjectOutput
		err    error
	)
	log.Debug().Str("bucket", bucket).Str("key", key).Msg("Collecting file")

	result, err = s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Err(err).Msg("Failed to collect file")
		return nil, err
	}
	defer result.Body.Close()

	res, err := io.ReadAll(result.Body)
	if err != nil {
		log.Err(err).Msg("Failed to read collected file")
	}
	return res, err
}

// SupportsPresign indicates whether the store supports presign
func (s *S3Storer) SupportsPresign() bool {
	return true
}

// GetUrl makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (s *S3Storer) getURL(bucket string, key string, ttlSec int64) (*PresignedRequest, error) {
	request, err := s.presign.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(ttlSec * int64(time.Second))
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get a presigned request to get %v:%v. cause: %v", bucket, key, err)
	}
	return &PresignedRequest{
		Method: request.Method,
		URL:    request.URL,
		Header: request.SignedHeader,
		Mode:   "no-cors",
	}, nil
}

// PresignThumbnail makes a presigned request that can be used to get a thumbnail.
func (s *S3Storer) PresignThumbnail(id string) (*PresignedRequest, error) {
	path := filepath.Join(prefix, id, thumbnailName)
	return s.getURL(s.thumbBucket, path, 300)
}

// PresignImage makes a presigned request that can be used to get a raw image.
func (s *S3Storer) PresignImage(id string) (*PresignedRequest, error) {
	path := filepath.Join(prefix, id, rawName)
	return s.getURL(s.rawBucket, path, s.presignedTTL)
}
