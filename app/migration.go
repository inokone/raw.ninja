package app

import (
	"fmt"
	"log"
	"os"

	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/photo"
)

func Migrate() {
	var err error
	if err = initDb(config.Database); err != nil {
		log.Fatal("Could not set up connection to database. Application spinning down.")
		os.Exit(1)
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err = DB.AutoMigrate(&photo.Photo{}, &auth.User{}, &descriptor.Descriptor{}); err != nil {
		log.Fatal("Database migration failed. Application spinning down.")
		os.Exit(1)
	}

	fmt.Println("Database migration finished")
}
