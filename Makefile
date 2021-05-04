.PHONY: run build clean tool lint

all: build

run:
	go run .

build:
	go build -v .

test:
	go test -v ./...

tool:
	go vet ./...; true
	gofmt -w .

lint:
	golangci-lint run

clean:
	rm -rf main
	go clean -i .
