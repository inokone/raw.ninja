package app

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/auth/account"
	"github.com/inokone/photostorage/auth/role"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var db *gorm.DB

func initConf(s string) error {
	conf, err := common.LoadConfig(s)
	if err != nil {
		return err
	}
	config = conf
	initLog()
	log.Info().Msg("Photostorage app starting up...")
	return nil
}

func initDb(c common.RDBConfig) error {
	cs := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v ", c.Host, c.Username, c.Password, c.Database, c.Port)
	if c.SSLMode == "disable" {
		cs += "sslmode=disable"
	} else {
		cs += fmt.Sprintf("sslrootcert=%v sslmode=%v", c.SSLCert, c.SSLMode)
	}
	gormDb, err := gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := gormDb.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime)
	db = gormDb
	return nil
}

func initStore(c common.ImageStoreConfig) {
	var is image.Storer = image.NewStorer(c)
	storers.Photos = photo.NewGORMStorer(db, is)
	storers.Users = user.NewGORMStorer(db)
	storers.Roles = role.NewGORMStorer(db)
	storers.Accounts = account.NewGORMStorer(db)
}

func initLog() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339
	level, err := zerolog.ParseLevel(config.Log.LogLevel)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse log level, default is debug.")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}
	if config.Log.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
