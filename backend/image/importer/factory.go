package importer

import "slices"

var (
	libraw     = []string{"cr2", "crw", "dng", "arw", "raw", "tiff", "nef", "raf"}
	compressed = []string{"jpg", "jpeg", "png", "gif"}
)

// NewImporter is a factory method of `Importer` based on file formats
func NewImporter(format string) Importer {
	if slices.Contains(libraw, format) {
		return NewLibrawImporter()
	}
	if slices.Contains(compressed, format) {
		return NewDefaultImporter()
	}
	panic("No importer found for format " + format)
}
