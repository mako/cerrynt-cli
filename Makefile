APP_NAME=cerrynt

.PHONY: run build test vet fmt tidy ci

run:
	go run ./cmd/$(APP_NAME)/main.go

build:
	go build ./cmd/$(APP_NAME)/main.go

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

ci: fmt vet test build
