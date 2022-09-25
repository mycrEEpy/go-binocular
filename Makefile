all: test

get:
	go mod download

check:
	go vet ./...

build: get
	go build -v ./...

test: check
	go test -v -coverprofile coverage.out -race ./...

cover: test
	go tool cover -html coverage.out

bench:
	go test -v -bench=. -run=^$

tidy:
	go mod tidy

fmt:
	go fmt ./...

clean:
	go clean
