/*
Photostorage is an application to store RAW image files.
Main focus is on cutting costs while maintaining security, high
availability and durability.

Usage:

	main [flags]

The flags are:

	    --migrate [=true/false]
	        When provided the application initiates a database migration
			before starting the application. Default value is false.
	    --application [=true/false]
	        Starts the web application for the photostorage. Default value
			is true.
	    --port [=0-65535]
	        The TCP port for the web application. Default value is 8080.
*/
package main

import (
	"flag"

	"github.com/inokone/photostorage/app"
)

func main() {
	// TODO: set up configuration options for logging here
	// TODO: add logging where necessary
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
