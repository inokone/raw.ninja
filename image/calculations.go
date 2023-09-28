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
	result := canvas(original.Bounds().Size().X, original.Bounds().Size().Y)
	draw.NearestNeighbor.Scale(result, result.Rect, original, original.Bounds(), draw.Over, nil)
	log.Debug().Dur("Elapsed time", time.Since(start)).Msg("Generated thumbnail.")
	return result, nil
}

func canvas(width int, height int) *image.RGBA {
	ratio := math.Min(thumbWidth/float64(width), thumbHeight/float64(height))
	newWidth := int(ratio * float64(width))
	newHeight := int(ratio * float64(height))
	result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	return result
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
