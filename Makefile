.PHONY: all build test

all: build

build:
	go build -o .bin/misen ./cmd/main.go

test:
	go test -v ./...
