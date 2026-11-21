.PHONY: all build run test clean docker-build docker-run

APP_NAME = jiosaavn-api

all: build

build:
	@echo "Building binary..."
	go build -ldflags="-s -w" -o bin/$(APP_NAME) ./cmd/api

run:
	@echo "Starting server..."
	go run ./cmd/api

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin/

docker-build:
	@echo "Building docker image..."
	docker build -t $(APP_NAME) .

docker-run:
	@echo "Running docker container..."
	docker run -p 8080:8080 $(APP_NAME)
