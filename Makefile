.PHONY: all clean fmt build run test

default: test

clean:
		go clean -cache -i ./...

fmt:
	go fmt ./...

build:
		go build ./...

test:
		go test -v ./...

run:
		go run ./...

all: clean fmt build test
