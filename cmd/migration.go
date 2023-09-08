package main

import (
	"fmt"
	"os/user"

	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/photo"
)

func Migrate() {
	common.InitDb(config.Database)
	common.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	common.DB.AutoMigrate(&user.User{}, &descriptor.Descriptor{}, &photo.Photo{})
	fmt.Println("Database migration finished")
}
