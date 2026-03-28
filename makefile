IMAGE_NAME := sab_monitarr

.PHONY: build up down test coverage

## build: build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

## up: start the service with docker compose (requires .env)
up:
	docker compose -f docker-compose.yml up -d --build

## down: stop and remove containers
down:
	docker compose down

## test: run all Go tests
test:
	go test ./... -v

## coverage: run tests and display total coverage percentage
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out | grep "^total:"
	@rm -f coverage.out
