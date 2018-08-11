NAME=c2c-demo
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

setup:
	go get github.com/golang/dep

install:
	dep ensure

update:
	dep ensure -update

test:
	go test ./...

run/%: cmd/%/main.go
	go run $<

build/%: cmd/%/main.go
	go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)-$*_$(VERSION) $<

info:
	@echo version: ${VERSION}
	@echo revision: ${REVISION}
