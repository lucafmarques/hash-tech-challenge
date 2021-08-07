## This help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

## Builds and runs checkout and discount services with docker-compose
run-services:
	@docker-compose build
	@docker-compose up --remove-orphans

## Builds and runs checkout service in docker
run:
	@docker build -t checkout .
	@docker run -d -p 8080:8080 checkout

## Runs test suite
test:
	@go test -v

## Runs test suite and outputs coverage
coverage:
	@go test -cover

## Builds docs, starts server and Swagger
docs:
	@swag init