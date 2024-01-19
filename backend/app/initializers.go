package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/inokone/photostorage/auth/account"
	"github.com/inokone/photostorage/auth/role"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/collection"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/onetime"
	"github.com/inokone/photostorage/photo"
	"github.com/inokone/photostorage/ruleset"
	"github.com/inokone/photostorage/ruleset/rule"
	"github.com/inokone/photostorage/web"

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

func initDb(c common.RDBConfig, l common.LogConfig) error {
	cs := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v ", c.Host, c.Username, c.Password, c.Database, c.Port)
	if c.SSLMode == "disable" {
		cs += "sslmode=disable"
	} else {
		cs += fmt.Sprintf("sslrootcert=%v sslmode=%v", c.SSLCert, c.SSLMode)
	}
	gormDb, err := gorm.Open(postgres.Open(cs), &gorm.Config{
		Logger: logger.Default.LogMode(levelFor(l.LogLevel)),
	})
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

func levelFor(level string) logger.LogLevel {
	l := strings.ToLower(strings.TrimSpace(level))
	switch l {
	default:
		return logger.Warn
	case "debug":
		return logger.Info
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	}
}

func initStorers(c common.ImageStoreConfig) {
	storers.Photos = photo.NewGORMStorer(db)
	storers.Images = image.NewStorer(c)
	storers.Users = user.NewGORMStorer(db)
	storers.Roles = role.NewGORMStorer(db)
	storers.Accounts = account.NewGORMStorer(db)
	storers.Collections = collection.NewGORMStorer(db)
	storers.Rules = rule.NewGORMStorer(db)
	storers.RuleSets = ruleset.NewGORMStorer(db)
	storers.OneTime = onetime.NewGORMStorer(db)
}

func initServices(c common.ImageStoreConfig, storers web.Storers) {
	services.Load = *photo.NewLoadService(storers.Photos, storers.Images, c)
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
