NAME=c2c-demo
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

setup:
	go get github.com/golang/dep
	go install github.com/golang/dep/cmd/dep

install:
	dep ensure

update:
	dep ensure -update

test:
	go test ./...

run:
	go run cmd/api/main.go

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)-$*_$(VERSION) cmd/api/main.go

info:
	@echo version: ${VERSION}
	@echo revision: ${REVISION}
