# RAW.Ninja

Source code of a decomissioned cheap photo storage for professional purposes. A personal project to learn Go language.

The target was an SAAS web application that supports Professional needs of image storage  - including uploading, editing and sharing raw and processed image files, assisting contract-based image retention.

The application was written in Go and React + JS, hosted in AWS and was available at https://raw.ninja

## Features

### User Management

- User authentication with email/password or social login (Google, Facebook)
- Email verification and password recovery
- User roles (Free tier, Administrator)
- User quotas and storage limits
- Profile management and statistics

### Image Management

- Upload RAW and processed images (single or batch)
- Multiple storage backends (local filesystem, AWS S3)
- Automatic thumbnail generation
- Image organization with albums
- Image rating and favorites
- Fullscreen image viewer with zoom
- Quick search across images and albums
- Editing raw image files
- Customizable retention rule sets based on upload time
- Analytics on uploded images and trends

### Security & Compliance

- HTTPS support (HTTP for local development)
- JWT-based authentication
- Cookie consent for EU compliance
- Terms of service and privacy policy
- Cloudflare protection and domain management
- Googlr recaptcha for login and registration flows

### Technical Infrastructure

- Docker-based deployment
- AWS S3 for image and blob storage - with retrieval
- PostgreSQL database for the rest with SSL and encryption
- SendGrid email integration
- AWS deployment (EC2, ALB)
- Automated database migrations
- API doc using Swagger

### Administrative Features

- User management dashboard
- System statistics and monitoring
- Storage and quota management
- Upload tracking and analytics

## Set up for development

### Backend

Created using Golang in the [backend](/backend) folder. The RAW processing is based on Libraw. The application uses 2 databases:

- **File storage:** either local storage or a CDN used for storing image blobs and thumbnail blobs
- **Postgres:** relational database - storing all data except blob

Configuration of the backend is environment variables or env file. Here is a [sample](/environments/local.env) env file. When in production [AWS Secrets Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html) should be used instead.
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
docker run -p 80:80 -v ~/git/raw.ninja/environments/local/certificates:/etc/rawninja/certificates rawninja-frontend

docker run -p 8080:8080 -v ~/git/raw.ninja/environments/local:/etc/rawninja --mount type=tmpfs,destination=/tmp/photos,tmpfs-size=4096 rawninja-backend
```

Also if we just want to spin the application up with local database:

``` sh
docker compose up
```
