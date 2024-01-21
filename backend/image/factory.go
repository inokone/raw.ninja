package image

import (
	"strings"

	"github.com/inokone/photostorage/common"
	"github.com/rs/zerolog/log"
)

const (
	localType = "local"
	s3Type    = "s3"
)

// NewStorer is a factory method of `Storer` based on configuration
func NewStorer(config *common.ImageStoreConfig) Storer {
	standard := strings.ToLower(strings.TrimSpace(config.Type))
	if standard == localType {
		result, err := NewLocalStorer(config.Path)
		if err != nil {
			log.Err(err).Msg("Failed to set up local storer!")
			panic("Failed to set up local storer!")
		}
		return result
	}
	if standard == s3Type {
		result, err := NewS3Storer(
			config.RawBucket,
			config.ThumbBucket,
			config.AwsKey,
			config.AwsSecret,
			config.PresignedTTL,
		)
		if err != nil {
			log.Err(err).Msg("Failed to set up S3 storer!")
			panic("Failed to set up S3 storer!")
		}
		return result
	}

	panic("No store found for type " + config.Type)
}
