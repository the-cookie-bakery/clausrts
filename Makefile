.PHONY: build run test vet lint assets clean

BINARY := bin/clausrts

build:
	go build -o $(BINARY) ./cmd/clausrts

run:
	go run ./cmd/clausrts

test:
	go test ./...

vet:
	go vet ./...

lint:
	@golangci-lint run 2>/dev/null || echo "golangci-lint not installed; run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

assets:
	go run ./cmd/genassets

clean:
	rm -rf bin
