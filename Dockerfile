FROM golang:1.13 AS build

WORKDIR /src
COPY . /src

RUN apt-get -y install libsqlite3-0 make git
RUN GO111MODULE=on go mod download
RUN GO111MODULE=on make build

FROM debian:latest

WORKDIR /app
COPY --from=build /src/golang-echo-realworld-example-app /app/

EXPOSE 8585

ENTRYPOINT "/app/golang-echo-realworld-example-app"
