FROM golang:1.21.4-bullseye as build
RUN apt-get update && apt-get install -y libraw-dev
RUN apt-get install -y --force-yes --no-install-recommends apt-transport-https curl ca-certificates \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY . .
RUN go build -v -o rawninja main.go

FROM debian:11.8-slim
RUN apt-get update && apt-get install -y libraw-dev
COPY --from=build /app/rawninja /usr/local/bin/rawninja
ENV GIN_MODE=release
CMD [ "/usr/local/bin/rawninja", "--migrate" ] 
EXPOSE 8080