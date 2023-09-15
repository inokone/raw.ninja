package app

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/photo"
)

func Migrate() {
	err := initDb(config.Database)
	if err != nil {
		log.Fatal("Could not set up connection to database. Application spinning down.")
		os.Exit(1)
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.AutoMigrate(&user.User{}, &descriptor.Descriptor{}, &photo.Photo{})
	fmt.Println("Database migration finished")
}
