.PHONY = default build clean

default: build

build:
	go build -o bin/gomate ./cmd/gomate

clean:
	rm bin/gomate
