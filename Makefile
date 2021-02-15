.PHONY: build run test

default: test

build:
		go build ./...

run:
		go run ./...

test:
		go test ./...
