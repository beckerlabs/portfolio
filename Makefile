
APP_NAME := beckerlabs/portfolio
VERSION ?= latest
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
DOCKERFILE := Dockerfile

.PHONY: all build test clean run

all: test build

build:
	@echo "Building the Docker image..."
	@docker build -t $(DOCKER_IMAGE) -f $(DOCKERFILE) .

test:
	@echo "Running Go tests..."
	@go test ./...

clean:
	@echo "Cleaning up Docker images..."
	@docker rmi -f $(DOCKER_IMAGE)

run:
	@echo "Running the application..."
	@docker run --rm -p 4000:4000 $(DOCKER_IMAGE)