package image

import (
	"bytes"
	"context"
	"io"
	"path/filepath"

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
	rawBucket   string
	thumbBucket string
	client      *s3.Client
}

// NewS3Storer creates a new LocalStorer with the specified storage path.
func NewS3Storer(rawBucket, thumbBucket, awsAPIKey, awsAPISecret string) (*S3Storer, error) {
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
		rawBucket:   rawBucket,
		thumbBucket: thumbBucket,
		client:      cl,
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
