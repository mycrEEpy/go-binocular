all: test

get:
	go mod download

check-init:
	go get golang.org/x/lint/golint
	go get honnef.co/go/tools/cmd/staticcheck
	go get github.com/kisielk/errcheck

check: check-init
	go vet ./...
	golint ./...
	staticcheck ./...
	errcheck ./...

build: get
	go build -v ./...

test: check
	go test -v -coverprofile coverage.out -race ./...

cover: test
	go tool cover -html coverage.out

tidy:
	go mod tidy

fmt:
	go fmt ./...

clean:
	go clean
