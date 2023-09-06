package migration

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/inokone/photostorage/common"
)

func latest(host string, port int, database string) error {
	m, err := migrate.New(
		"file:///migrations",
		fmt.Sprintf("postgres://%v:%v/%v?sslmode=enable", host, port, database))
	if err != nil {
		log.Fatalf("Migartions failed. Cause: %v", err)
		return err
	}
	m.Up()
	return nil
}

func Migrate() {
	var config common.RDBConfig
	config.New()
	latest(config.Host, config.Port, config.Database)
}
