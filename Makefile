.PHONY: build run test fmt vet lint

build:
	go build -o bin/yearcollage ./cmd/yearcollage

run:
	go run ./cmd/yearcollage --input ./bilder --output collage.jpg

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...

test:
	go test ./...

lint: vet
	@echo "add golangci-lint here when ready"
