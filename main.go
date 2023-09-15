package main

import (
	"flag"

	"github.com/inokone/photostorage/app"
)

func main() {
	var (
		isMigration = flag.Bool("migrate", false, "Start migration of the database. Default: [false]")
		application = flag.Bool("application", true, "Start the web application on the provided port. Default: [true].")
		port        = flag.Int("port", 8080, "Port of the webapplication. Default: [8080]")
	)
	flag.Parse()

	if *isMigration {
		app.Migrate()
	}
	if *application {
		app.App(*port)
	}
}
