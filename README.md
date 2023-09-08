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
```
go install github.com/cosmtrek/air@latest // hot-reload for Gin server
go install github.com/swaggo/swag/cmd/swag@latest // OpenAPI spec generation for REST endpoints
```

## Build
```
cd cmd
swag init  // re-generate the OpenAPI spec files
go build . // build Go based application
```
If `swag` is not working you might have to add `~/go/bin` to  your PATH.

## Run
```
docker-compose up -d    // initialize Postgres database

cd cmd
go run --migrate        // migrate the database

go run .                // start the web-application, or
air                     // start the web application with hot-reload

docker-compose down     // stop runninf Postgres database 
```

## CI
The project has Github actions set up for every push.
Steps included
- OpenAPI re-generation
- Build 
- Run unit tests
- Static code analysis (planned)

## API doc
When the application is running, the OpenAPI documentation is available with [Swagger](http://localhost:8080/swagger/doc.json).

