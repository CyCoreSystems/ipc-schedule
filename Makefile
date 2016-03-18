all: assets test build

test:
	go test

generate:
	go generate

build: generate
	go build

assets:
	webpack
	go generate

