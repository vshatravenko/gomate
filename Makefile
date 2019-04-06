.PHONY = default build clean

default: build

build:
	go build -o bin/gomate ./cmd/gomate

test:
	go test ./...

clean:
	rm bin/gomate
