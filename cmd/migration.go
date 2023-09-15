package main

import (
	"fmt"
	"os/user"

	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/photo"
)

func Migrate() {
	InitDb(config.Database)
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.AutoMigrate(&user.User{}, &descriptor.Descriptor{}, &photo.Photo{})
	fmt.Println("Database migration finished")
}
