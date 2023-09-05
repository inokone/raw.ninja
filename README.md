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

## Build
```
cd cmd
go build .
```

## Run
```
cd cmd
go run .
```

## CI
The project has Github actions set up on every push.
Steps included
- OpenAPI generation
- Build 
- Unit tests
- Static code analysis (planned)

## API doc
When the application is running, the OpenAPI documentation is available with [Swagger](http://localhost:8080/swagger/doc.json).
