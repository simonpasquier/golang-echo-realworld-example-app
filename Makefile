export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export DEBUG=true
export APP=golang-echo-realworld-example-app
export LDFLAGS="-w -s"

APP_REVISION := $(shell git rev-parse HEAD)
APP_VERSION := $(shell git describe --tags --abbrev=0)
APP_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

all: build test

build:
	go build -tags netgo -ldflags "-X 'main.appVersion=$(APP_VERSION)' -X 'main.appRevision=$(APP_REVISION)' -X 'main.appBranch=$(APP_BRANCH)'"

run:
	go run -race .

############################################################
# Test
############################################################

test:
	go test -v -race ./...

container:
	docker build -t echo-realworld .

run-container:
	docker run --rm -it echo-realworld

.PHONY: build run test container
