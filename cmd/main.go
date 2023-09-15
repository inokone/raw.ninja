package main

import (
	"flag"
	"log"
	"os"

	"gorm.io/gorm"

	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/image"
)

var config *common.AppConfig
var DB *gorm.DB
var IS *image.Store

func init() {
	conf, err := common.LoadConfig()
	if err != nil {
		log.Fatal("Could not load application configuration", err)
	}
	config = conf
}

func main() {
	var (
		migrate     = flag.Bool("migrate", false, "Start migration of the database. Default: [false]")
		application = flag.Bool("application", true, "Start the web application on the provided port. Default: [true].")
		port        = flag.Int("port", 8080, "Port of the webapplication. Default: [8080]")
	)
	flag.Parse()

	err := InitDb(config.Database)
	if err != nil {
		log.Fatal("Could not set up connection to database. Application spinning down.")
		os.Exit(1)
	}

	err = InitStore(config.Store)
	if err != nil {
		log.Fatal("Could not set up image store. Application spinning down.")
		os.Exit(1)
	}

	if *migrate {
		Migrate()
	}
	if *application {
		App(*port)
	}
}
