.PHONY: all

all: build

build:
	go build -o .bin/misen ./cmd/main.go

