all: fmt vet

fmt:
	go fmt ./...

vet:
	go vet ./...

api: run

run:
	go run cmd/api/*.go

migrations:
	go run cmd/migrator/*.go
