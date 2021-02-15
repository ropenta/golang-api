.PHONY: all clean fmt build run test

default: test

all:
		make clean
		make fmt
		make build
		make test

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
