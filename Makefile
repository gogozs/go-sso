all: fmt lint test

fmt:
	go fmt

lint:
	golangci-lint run

test:
	go test -v ./...