package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"math"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
)

const (
	thumbWidth  float64 = 200
	thumbHeight float64 = 200
)

func Thumbnail(original image.Image) (image.Image, error) {
	start := time.Now()
	ratio := math.Min(thumbWidth/float64(original.Bounds().Size().X), thumbHeight/float64(original.Bounds().Size().Y))
	width := int(ratio * float64(original.Bounds().Size().X))
	height := int(ratio * float64(original.Bounds().Size().Y))
	result := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(result, result.Rect, original, original.Bounds(), draw.Over, nil)
	log.Debug().Dur("Elapsed time", time.Since(start)).Msg("Generated thumbnail.")
	return result, nil
}

func ExportJpeg(image image.Image) ([]byte, error) {
	start := time.Now()
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)
	log.Debug().Dur("Elapsed time", time.Since(start)).Msg("Exported thumbnail.")
	return buf.Bytes(), err
}

func ImportJpeg(b []byte) (image.Image, error) {
	return jpeg.Decode(bytes.NewReader(b))
}
