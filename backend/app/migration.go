package app

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/inokone/photostorage/auth/account"
	"github.com/inokone/photostorage/auth/role"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo"
)

// Migrate executes the necessary database initialization and migration
func Migrate(path string) {
	var err error

	if err = initConf(path); err != nil {
		log.Err(err).Msg("Failed to load application configuration.")
		os.Exit(1)
	}

	if err = initDb(config.Database); err != nil {
		log.Err(err).Msg("Failed to set up connection to database. Application spinning down.")
		os.Exit(1)
	}

	res := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if res.Error != nil {
		log.Err(res.Error).Msg("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	if err := db.AutoMigrate(&photo.Photo{}, &role.Role{}, &user.User{}, &descriptor.Descriptor{}, &image.Metadata{}, &account.Account{}); err != nil {
		log.Err(err).Msg("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	var rCount int
	db.Raw("SELECT count(*) FROM roles").Scan(&rCount)
	if rCount > 0 {
		log.Info().Msg("Roles already present, database migration finished")
		return
	}

	res = db.Exec("INSERT INTO roles (role_type, quota, display_name) VALUES (0, -1, 'Admin')")
	if res.Error != nil {
		log.Err(res.Error).Msg("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	res = db.Exec("INSERT INTO roles (role_type, quota, display_name) VALUES (1, 524288000, 'Free Tier')") // 500 Mb limit for free tier
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	log.Info().Msg("Database migration finished")
}
