package app

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo"
)

// Migrate executes the necessary database initialization and migration
func Migrate() {
	if err := initDb(config.Database); err != nil {
		log.Error().Err(err).Msg("Failed to set up connection to database. Application spinning down.")
		os.Exit(1)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err := db.AutoMigrate(&photo.Photo{}, &auth.User{}, &descriptor.Descriptor{}, &image.Metadata{}); err != nil {
		log.Error().Err(err).Msg("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	log.Info().Msg("Database migration finished")
}
