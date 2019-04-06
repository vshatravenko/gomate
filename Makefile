.PHONY = default build clean

default: build

build:
	go build -o bin/gomate ./cmd/gomate

test:
	go test ./...

bench:
	go test ./... -bench .

clean:
	rm bin/gomate
