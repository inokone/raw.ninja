# photostorage

Cheap photo storage for professional purposes. A personal project to learn Go language.
The target is a web application capable of handling reltively large image files (>50Mb), generating thumbnails.

## Functional requirements (planned)

- Authentication/Authorization including OpenId Connect
- Uploading and storing RAW image files of various types
- Browse accessible files
- Download/Delete single/selected file(s)
- Share files between users
- Setting up lifecycle for uploaded images
- Marking favorite images (exempt from lifecycle rules)

## Set up for development

``` sh
brew install libraw                                            # RAW processing library on OSX, or
sudo apt-get install libraw-dev                                # on Ubuntu

go install github.com/cosmtrek/air@latest                      # Hot-reload for Gin server
go install github.com/swaggo/swag/cmd/swag@latest              # OpenAPI spec generator
go install github.com/go-critic/go-critic/cmd/gocritic@latest  # Static code anlanysis for Go
```

## Build

``` sh
swag i -d "./,./app,./common,./photo,./descriptor,./web"  # Gernerate OpenAPI spec files
go build main.go                                          # Build app

go test -v ./...                                          # Run unit tests
gocritic check ./...                                      # Run static code analysis
```

On OSX if `swag` and `gocritic` are not working you might have to add `~/go/bin` to your PATH.

## Run

``` sh
docker-compose up -d      # Initialize Postgres database

go run main.go --migrate  # Migrate the database and launch app, or
go run main.go            # Start the web-application, or
air                       # Start the web application with hot-reload for development

docker-compose down       # Stop running Postgres database 
```

For production deployment please use `GIN_MODE=release` env variable.

## CI

The project has Github actions set up for every push.
Steps included

- OpenAPI re-generation
- Build
- Run unit tests
- Static code analysis

## API doc

When the application is running, the OpenAPI documentation is available with [Swagger](http://localhost:8080/swagger/doc.json).
