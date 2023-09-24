package app

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
)

func initDb(c common.RDBConfig) error {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", c.Host, c.Username, c.Password, c.Database, c.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime)
	DB = db
	return nil
}

func initStore(c common.ImageStoreConfig) error {
	if c.Type == "file" {
		var r image.Repository
		r, error := image.NewLocalStore(c.Path)
		if error != nil {
			return error
		}
		IS = &r
		return nil
	}
	return gorm.ErrNotImplemented
}
