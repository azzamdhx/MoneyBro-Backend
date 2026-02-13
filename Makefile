.PHONY: run build generate migrate-up migrate-down test clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

generate:
	go get github.com/99designs/gqlgen
	go run github.com/99designs/gqlgen generate

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

deps:
	go mod download
	go mod tidy
