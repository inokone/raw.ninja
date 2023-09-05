package photo

import (
	"image"
	"time"
)

type Format int64

const (
	GP2 Format = iota
	ARW
)

type Descriptor struct {
	FileName  string
	Uploaded  time.Time
	Format    Format
	Width     int
	Height    int
	Thumbnail image.Image
	Tags      []string
}

type Photo struct {
	Image image.Image
	Desc  Descriptor
}
