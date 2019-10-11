#FROM registry.access.redhat.com/ubi8/ubi AS build
FROM golang:1.13 AS build

WORKDIR /src
COPY . /src

#RUN dnf -y install golang sqlite make git
RUN apt-get -y install libsqlite3-0 make git
RUN GO111MODULE=on go mod download
RUN GO111MODULE=on make build

#FROM registry.access.redhat.com/ubi8/ubi-minimal
#FROM debian:latest
FROM golang:1.13

WORKDIR /app
COPY --from=build /src/golang-echo-realworld-example-app /app/

#RUN microdnf install sqlite
#RUN apt-get update && apt-get install -y libsqlite3-0 && rm -rf /var/lib/apt/lists/*

EXPOSE 8585

ENTRYPOINT "/app/golang-echo-realworld-example-app"
