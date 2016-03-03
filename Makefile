all: assets test build

test:
	go test

build:
	go build

assets:
	rice embed-go

