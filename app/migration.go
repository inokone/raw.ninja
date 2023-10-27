package app

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/auth/role"
	"github.com/inokone/photostorage/auth/user"
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

	res := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	if err := db.AutoMigrate(&photo.Photo{}, &role.Role{}, &user.User{}, &descriptor.Descriptor{}, &image.Metadata{}, &auth.AuthenticationState{}); err != nil {
		log.Error().Err(err).Msg("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	res = db.Exec("INSERT INTO roles (role_type, quota, display_name) VALUES (0, -1, 'Admin')")
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	res = db.Exec("INSERT INTO roles (role_type, quota, display_name) VALUES (1, 500000000, 'Free Tier')") // 500 Mb limit for free tier
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	log.Info().Msg("Database migration finished")
}
