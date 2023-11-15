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
	    --config [path]
		    Path of the configuration folder where the app.env config file
			is present. Default value is "."
*/
package main

import (
	"flag"

	"github.com/inokone/photostorage/app"
)

// @title                     RAW.Ninja API
// @version                   0.1
// @description               RAW.Ninja is an application to store RAW image files.
// @BasePath                  /api/v1
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	var (
		isMigration = flag.Bool("migrate", false, "Start migration of the database. Default: [false]")
		application = flag.Bool("application", true, "Start the web application on the provided port. Default: [true].")
		config      = flag.String("config", ".", "Path of the configuration folder where the app.env file is. Default: [.]")
	)
	flag.Parse()

	if *isMigration {
		app.Migrate(*config)
	}
	if *application {
		app.App(*config)
	}
}
