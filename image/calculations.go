package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"math"
	"time"

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
	log.Printf("Generated thumbnail in %v", time.Since(start))
	return result, nil
}

func ExportJpeg(image image.Image) ([]byte, error) {
	start := time.Now()
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)
	log.Printf("Exported thumbnail in %v", time.Since(start))
	return buf.Bytes(), err
}

func ImportJpeg(b []byte) (image.Image, error) {
	return jpeg.Decode(bytes.NewReader(b))
}
