# RAW.Ninja

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

Created using Golang in the [backend](/backend) folder. The RAW processing is based on Libraw. The application uses 2 databases:

- **File storage:** either local storage or a CDN used for storing image blobs and thumbnail blobs
- **Postgres:** relational database - storing all data except blob

Configuration of the backend is environment variables or env file. Here is a [sample](/environments/local.env) env file.
The backend also uses the following set of external files and folders:

- **tmp:** The system temp folder is used by Libraw to store temporary files while processing RAW images.
- **local storage folder:** If local storage is used a folder need to be set up for it and configured for the application using `IMG_STORE_PATH` env variable.
- **web ssl:** If https is the target for the application (must for a prod deployment) it can be set using `TLS_CERT_PATH` and `TLS_KEY_PATH`
- **database ssl:** Database access can be encrypted too, can be set using `DB_SSL_MODE` and `DB_SSL_CERT` env variables.

#### Development Prerequisites

Before starting backend development the following need to be set up:

``` sh
brew install go                                    # Install Go
brew install libraw                                # Install RAW processing library on OSX, or
sudo apt-get install libraw-dev                    # on Ubuntu

go install github.com/cosmtrek/air@latest          # Hot-reload for Gin server
go install github.com/swaggo/swag/cmd/swag@latest  # OpenAPI spec generator

brew tap golangci/tap                              # Setting source for brew, then
brew install golangci/tap/golangci-lint            # Static code anlanysis for Go
```

Note: on M1/M2 OSX you need to manually install Libraw based on the [official doc](https://www.libraw.org/docs/Install-LibRaw-eng.html).

#### Build and Development

Building the application is not explicitly required for development. The following commands can be used:

``` sh
go mod download    # Download Go dependencies
swag i    # Generate OpenAPI spec files
go build main.go   # Build app

go test -v ./...   # Run unit tests
golangci-lint run  # Run static code analysis
```

On OSX if `swag` is not working you might have to add `~/go/bin` to your PATH.

#### Run

``` sh
go run main.go --migrate --config ../environments/development  # Migrate the database and launch app
air                                                            # Start the web application with hot-reload for development
```

#### API doc

When the application is running, the OpenAPI documentation is available with [Swagger](http://localhost:8080/swagger/doc.json).

#### Production environment

In production we need SSL/TLS set up - if we do not have a reverse proxy set up. For that we need a certificate and a private key
The `app.env` file needs to contain the keys pointing to the certificate. The application will automatically pick it up.
Also the database certification should be set up. e.g. for AWS RDS on eu-central-1 region:

``` sh
DB_SSL_CERT=eu-central-1-bundle.pem 
DB_SSL_MODE=verify-full
```

``` sh
GIN_MODE=release go run main.go
```

### Frontend

Created using React.js with `npx create-react-app` in the [frontend](/frontend) folder.

#### Before Development

``` sh
brew install node                          # Install Node.js
npm install react@latest react-dom@latest  # Fix react version 
```

Install "React Developer Tools" browser extension.

#### Development

Standard build mechanism with Node for frontend

``` sh
npm install           # Download frontend dependencies 
npm start             # Run the app in dev mode
HTTPS=true npm start  # Run the app in dev mode with self signed certification over https
```

#### Production

In production we need SSL/TLS set up, for that we need a certificate and a private key. These can be used with the following command to run the application:

``` sh
npm run build  # Build frontend for production

serve -s build -p 443 --ssl-cert "/etc/rawninja/certificates/mycert.crt" --ssl-key "/etc/rawninja/certificates/mykey.key"
```

## Docker

First version of containerization result is 2 docker images. One for the frontend and one for the backend.

### Build

The following commands can be used for building the docker images for developmen:

``` sh
cd frontend
docker build -t rawninja-frontend -f Dockerfile . # for production build use the --build-arg PRODUCTION=1 flag

cd backend
docker build -t rawninja-backend -f Dockerfile . 
```

### Run

There are multiple options for running the application locally.

The following commands can be used for running the individual docker images:

``` sh
docker run -p 80:80 -v /Users/inokone/git/raw.ninja/environments/local/certificates:/etc/rawninja/certificates rawninja-frontend

docker run -p 8080:8080 -v /Users/inokone/git/raw.ninja/environments/local:/etc/rawninja --mount type=tmpfs,destination=/tmp/photos,tmpfs-size=4096 rawninja-backend
```

Also if we just want to spin the application up with local database:

``` sh
docker compose up
```

## Deploying to EC2

Build:

``` sh
cd frontend
docker build --build-arg PRODUCTION=1 -t rawninja-frontend -f Dockerfile .
cd backend
docker build -t rawninja-backend -f Dockerfile . 
```

From dev machine:

``` sh
cd ~/Downloads
docker save -o backend.tar rawninja-backend
docker save -o frontend.tar rawninja-frontend

scp -i rawninja-ec2-kp.pem frontend.tar ec2-user@3.123.42.65:~/
scp -i rawninja-ec2-kp.pem backend.tar ec2-user@3.123.42.65:~/
scp -r -i rawninja-ec2-kp.pem ../git/raw.ninja/environments/production ec2-user@3.123.42.65:~/

```

SSH into EC2:

Test database access:

``` sh
psql -h rawninja-rds.c9xvfg3kuua1.eu-central-1.rds.amazonaws.com -p 5432 -U postgres -d postgres
```

Start up the containers for the service:

``` sh
sudo docker load -i backend.tar
sudo docker load -i frontend.tar

sudo docker run -d --restart always -p 80:80 rawninja-frontend &
sudo docker run -d --restart always -p 8080:8080 -v ~/production:/etc/rawninja --mount type=tmpfs,destination=/tmp/photos,tmpfs-size=4096 --mount type=bind,source=/etc/ssl/certs,target=/etc/ssl/certs rawninja-backend &
```

or just the frontend:

``` sh
sudo nohup serve -p 80 -s build
```

## CI

The project has Github actions set up for every push.
Steps included

- [Backend](.github/workflows/build.yaml)
  - OpenAPI re-generation
  - Build
  - Run unit tests
- [Backend Static code analysis](.github/workflows/golangci-lint.yml)
