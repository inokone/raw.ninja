package importer

import (
	"image"

	img "github.com/inokone/photostorage/image"
)

// Importer is an interface for importing RAW camera images.
type Importer interface {
	Image(raw []byte) (*image.Image, error)

	Describe(raw []byte) (*img.Metadata, error)

	Thumbnail(raw []byte) (*img.ThumbnailImg, error)
}
