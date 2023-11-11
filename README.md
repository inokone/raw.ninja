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

### Backend

Created using Golang in the base folder of the repository of the repository.

#### Prerequisites

``` sh
brew install go                                                # Install Go
brew install libraw                                            # RAW processing library on OSX, or
sudo apt-get install libraw-dev                                # on Ubuntu

go install github.com/cosmtrek/air@latest                      # Hot-reload for Gin server
go install github.com/swaggo/swag/cmd/swag@latest              # OpenAPI spec generator

brew tap golangci/tap                                          # Setting source for brew, then
brew install golangci/tap/golangci-lint                        # Static code anlanysis for Go
```

Note: on M1/M2 OSX you need to manually install Libraw based on the [official doc](https://www.libraw.org/docs/Install-LibRaw-eng.html).

#### Build

``` sh
go mod download    # Download Go dependencies
swag i             # Generate OpenAPI spec files
go build main.go   # Build app

go test -v ./...   # Run unit tests
golangci-lint run  # Run static code analysis
```

On OSX if `swag` and `gocritic` are not working you might have to add `~/go/bin` to your PATH.

#### Run

``` sh
docker-compose up -d      # Initialize Postgres database

go run main.go --migrate  # Migrate the database and launch app, or
go run main.go            # Start the web-application, or
air                       # Start the web application with hot-reload for development

docker-compose down       # Stop running Postgres database 
```

#### Production environment

In production we need SSL/TLS set up, for that we need a certificate and a private key
The `app.env` file needs to contain the keys pointing to the certificate. The application will automatically pick it up.
Also the database certification should be set up. e.g. for AWS RDS on eu-central-1 region:

``` sh
DB_SSL_CERT=eu-central-1-bundle.pem 
DB_SSL_MODE=verify-full
```

``` sh
GIN_MODE=release go run main.go
```

#### API doc

When the application is running, the OpenAPI documentation is available with [Swagger](http://localhost:8080/swagger/doc.json).

### Frontend

Created using React.js with `npx create-react-app` in the [frontend](/web/frontend/photostorage) folder.

#### Prerequisites

``` sh
brew install node                          # Install Node.js
npm install react@latest react-dom@latest  # Fix react version 
```

Install "React Developer Tools" browser extension.

#### Build

Standard build mechanism with Node for frontend

``` sh
npm install           # Download frontend dependencies 
npm start             # Run the app in dev mode
HTTPS=true npm start  # Run the app in dev mode with self signed certification over https
```

#### Production environment

In production we need SSL/TLS set up, for that we need a certificate and a private key. These can be used with the following command to run the application:

``` sh
npm run build  # Build frontend for production

serve -s build --p 443 --ssl-cert "/etc/ssl/certs/mycert.crt" --ssl-key "/etc/ssl/private/mykey.key"

```

## CI

The project has Github actions set up for every push.
Steps included

- Backend
  - OpenAPI re-generation
  - Build
  - Run unit tests
  - Static code analysis
- Frontend
