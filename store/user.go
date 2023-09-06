package store

import (
	"database/sql"
	"fmt"

	"github.com/inokone/photostorage/common"

	_ "github.com/lib/pq"
)

type UserStore interface {
	Store(id string, user common.User) error

	User(id string) (common.User, error)

	Delete(id string) error
}

type RDBUserStore struct {
	config *common.RDBConfig
	db     *sql.DB
}

func (s RDBUserStore) New() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.config.Host, s.config.Port, s.config.Username, s.config.Password, s.config.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(s.config.MaxIdleConns)
	db.SetMaxOpenConns(s.config.MaxOpenConns)
	db.SetConnMaxLifetime(s.config.ConnMaxLifetime)
	s.db = db
	return nil
}

func (s RDBUserStore) Store(id string, user common.User) error {

}

func (s RDBUserStore) User(id string) (common.User, error) {

}

func (s RDBUserStore) Delete(id string) error {

}
