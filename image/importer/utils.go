package importer

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func tempFile(target string, content []byte) (string, error) {
	f, err := os.CreateTemp("", target+"_*")
	if err != nil {
		return "", err
	}
	defer closeTempFile(f)

	_, err = f.Write(content)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func tempPath(target string) string {
	tempDir := os.TempDir()
	return filepath.Join(tempDir, target+"_"+uuid.New().String())
}

func removeTempFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Warn().AnErr("Cause", err).Str("Path", path).Msg("Could not remove temp file.")
	}
}

func closeTempFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Warn().AnErr("Cause", err).Str("Name", f.Name()).Msg("Could not close temp file.")
	}
}
