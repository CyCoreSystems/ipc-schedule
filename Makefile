all: assets test build

test:
	go test

build:
	go build

assets:
	webpack
	go generate

