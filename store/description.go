package store

import (
	"fmt"

	"github.com/inokone/photostorage/photo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/common"
)

type DescriptorStore interface {
	Store(id string, desc photo.Descriptor) error

	Descriptions(ids ...string) ([]photo.Descriptor, error)

	Delete(id string) error
}

type GormDescriptorStore struct {
	config *common.RDBConfig
	db     *gorm.DB
}

func (s GormDescriptorStore) New() error {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=enable",
		s.config.Host, s.config.Username, s.config.Password, s.config.Database, s.config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	s.db = db

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(s.config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(s.config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(s.config.ConnMaxLifetime)
	return nil
}
