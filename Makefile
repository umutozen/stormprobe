.PHONY: build test lint clean

build:
	go build -o stormprobe ./cmd

test:
	go vet ./...
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f stormprobe
	rm -rf outputs/
